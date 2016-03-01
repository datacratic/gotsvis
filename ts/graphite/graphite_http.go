// Copyright (c) 2014 Datacratic. All rights reserved.

package graphite

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/datacratic/gotsvis/ts"
)

const (
	TimeoutError = "timeout-error"
)

type Graphite struct {
	URL string

	Client *http.Client
}

func (graph *Graphite) Init() {
	if graph.URL == "" {
		panic("graphite URL can't be empty")
	}
	if graph.Client == nil {
		graph.Client = http.DefaultClient
	}
}

func (graph *Graphite) Do(req GraphiteRequest) *Response {

	query := req.GetQuery()
	query.Add("format", "raw")
	respHTTP, err := graph.Client.Get(graph.URL + "/render?" + query.Encode())

	resp := &Response{
		Request: req,
		Error:   err,
		data:    make(ts.TimeSeriesSlice, 0),
	}
	if respHTTP == nil || err != nil {
		if err, ok := err.(*url.Error); ok {
			if err, ok := err.Err.(net.Error); ok {
				if err.Timeout() {
					resp.Error = fmt.Errorf("%s: %s", TimeoutError, err)
					return resp
				}
			}
		}
		fmt.Println("respHTTP nil or err:", respHTTP, err)
	}

	resp.Code = respHTTP.StatusCode
	resp.Body, err = ioutil.ReadAll(respHTTP.Body)
	respHTTP.Body.Close()
	if resp.Error != nil {
		resp.Error = fmt.Errorf("%s, %s", resp.Error, err)
	} else {
		resp.Error = err
	}

	return resp
}

type GraphiteRequest interface {
	GetQuery() url.Values
}

type Request struct {
	Key   string
	From  time.Time
	Until time.Time
}

func (req *Request) GetQuery() url.Values {
	query := make(url.Values)

	query.Add("target", req.Key)

	if !req.From.IsZero() {
		query.Add("from", strconv.FormatInt(req.From.Unix(), 10))
	}
	if !req.Until.IsZero() {
		query.Add("until", strconv.FormatInt(req.Until.Unix(), 10))
	}
	return query
}

type Requests []Request

func (reqs Requests) GetQuery() url.Values {
	query := make(url.Values)
	var from, until time.Time

	for _, r := range reqs {
		query.Add("target", r.Key)
		if from.IsZero() {
			from = r.From
		} else if from.After(r.From) {
			from = r.From
		}
		if until.IsZero() {
			until = r.Until
		} else if until.Before(r.Until) {
			until = r.Until
		}
	}
	if !from.IsZero() {
		query.Add("from", strconv.FormatInt(from.Unix(), 10))
	}
	if !until.IsZero() {
		query.Add("until", strconv.FormatInt(until.Unix(), 10))
	}
	return query
}

// Response struct from a graphite request.
type Response struct {
	Request GraphiteRequest
	Body    []byte
	Code    int
	Error   error

	data ts.TimeSeriesSlice
}

// Single is used to get a single time series, and verify that there was only one
// time series in the response.
func (resp *Response) Single() (*ts.TimeSeries, error) {
	if err := resp.checkError(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(resp.Body)
	for line, err := buf.ReadBytes('\n'); err == nil; line, err = buf.ReadBytes('\n') {
		if len(resp.data) > 0 {
			resp.Error = errors.New("response length is larger than a single time series")
			return nil, resp.Error
		}

		if ts, err := resp.readLine(line); err != nil {
			resp.Error = err
			return nil, resp.Error
		} else {
			resp.data = append(resp.data, *ts)
		}
	}

	if len(resp.data) == 0 {
		resp.Error = errors.New("no data in response")
		return nil, resp.Error
	}
	return &resp.data[0], nil
}

// First is used to retrieve the first time series from a response.
func (resp *Response) First() (*ts.TimeSeries, error) {
	if err := resp.checkError(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(resp.Body)
	for line, err := buf.ReadBytes('\n'); err == nil; line, err = buf.ReadBytes('\n') {
		if ts, err := resp.readLine(line); err != nil {
			resp.Error = err
			return nil, resp.Error
		} else {
			resp.data = append(resp.data, *ts)
			return ts, nil
		}
	}

	if len(resp.data) == 0 {
		resp.Error = errors.New("no data in response")
		return nil, resp.Error
	}
	return &resp.data[0], nil
}

// All is used to retrive all the time series in a response.
func (resp *Response) All() (ts.TimeSeriesSlice, error) {
	if err := resp.checkError(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(resp.Body)
	for line, err := buf.ReadBytes('\n'); err == nil; line, err = buf.ReadBytes('\n') {
		if ts, err := resp.readLine(line); err != nil {
			resp.Error = err
			return nil, resp.Error
		} else {
			resp.data = append(resp.data, *ts)
		}
	}

	if len(resp.data) == 0 {
		resp.Error = errors.New("no data in response")
		return nil, resp.Error
	}
	return resp.data, nil
}

func (resp *Response) GetMatchingAll(key string) (ts.TimeSeriesSlice, error) {
	tss, err := resp.All()
	if err != nil {
		return nil, err
	}
	return tss.GetMatching(key), nil
}

func (resp *Response) GetMatchingSingle(key string) (*ts.TimeSeries, error) {
	matching, err := resp.GetMatchingAll(key)
	if err != nil {
		return nil, err
	}
	if len(matching) == 0 {
		return nil, fmt.Errorf("key matching '%s' not found in response", key)
	}
	return &matching[0], nil
}

func (resp *Response) checkError() error {
	if resp.Code != 200 {
		return fmt.Errorf("response returned status '%d' != 200", resp.Code)
	}
	return resp.Error
}

func (resp *Response) readLine(line []byte) (*ts.TimeSeries, error) {
	var key string
	var start, end time.Time
	var step time.Duration
	var data []float64

	buf := bytes.NewBuffer(line)
	if header, err := buf.ReadBytes('|'); err != nil {
		return nil, err
	} else {
		header = header[:len(header)-1]
		headers := bytes.Split(header, []byte(","))
		if len(headers) < 4 {
			resp.Error = errors.New("graphite data header size <= 4")
			return nil, resp.Error
		}

		key = string(bytes.Join(headers[:len(headers)-3], []byte(",")))
		startInt, err := strconv.ParseInt(string(headers[len(headers)-3]), 10, 64)
		if err != nil {
			return nil, errors.New("conv to int fail")
		}

		endInt, err := strconv.ParseInt(string(headers[len(headers)-2]), 10, 64)
		if err != nil {
			return nil, errors.New("conv to int fail")
		}

		step, err = time.ParseDuration(string(headers[len(headers)-1]) + "s")
		if err != nil {
			return nil, errors.New("step couldn't be determined")
		}

		start = time.Unix(startInt, 0)
		end = time.Unix(endInt, 0)

		points := end.Sub(start) / step
		data = make([]float64, points)
	}

	var err, errP error
	var point []byte
	for i := 0; err == nil; i++ {
		point, err = buf.ReadBytes(',')
		point = point[:len(point)-1]

		if bytes.Equal(point, []byte("None")) {
			data[i] = math.NaN()
			continue
		}
		data[i], errP = strconv.ParseFloat(string(point), 64)
		if errP != nil {
			resp.Error = errP
			return nil, errP
		}
	}
	if err != io.EOF {
		resp.Error = err
		return nil, resp.Error
	}
	return ts.NewTimeSeriesOfData(key, start, step, data)
}
