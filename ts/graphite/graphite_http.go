// Copyright (c) 2014 Datacratic. All rights reserved.

package graphite

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gotsvis2/ts"
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

func (graph *Graphite) Do(req Request) *Response {

	query := req.GetQuery()
	query.Add("format", "raw")
	respHTTP, err := graph.Client.Get(graph.URL + "/render?" + query.Encode())
	//fmt.Println(string(resp.Body))

	resp := &Response{
		Request: &req,
		Error:   err,
		Code:    respHTTP.StatusCode,
		data:    make(ts.TimeSeriesSlice, 0),
	}
	resp.Body, err = ioutil.ReadAll(respHTTP.Body)
	respHTTP.Body.Close()
	if resp.Error != nil {
		resp.Error = fmt.Errorf("%s, %s", resp.Error, err)
	} else {
		resp.Error = err
	}

	return resp
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

// Response struct from a graphite request.
type Response struct {
	Request *Request
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
		fmt.Println("line", string(line), err)
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
	buf = bytes.NewBuffer(resp.Body)
	line, err := buf.ReadBytes('\n')
	fmt.Println("line2", string(line), err)

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
			fmt.Println("headers", headers)
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

		fmt.Println("point", string(point))
		if bytes.Equal(point, []byte("None")) {
			data[i] = math.NaN()
			continue
		}
		data[i], errP = strconv.ParseFloat(string(point), 64)
		if errP != nil {
			fmt.Println("here", string(point), errP)
			resp.Error = errP
			return nil, errP
		}
	}
	if err != io.EOF {
		resp.Error = err
		return nil, resp.Error
	}
	fmt.Println(string(resp.Body))
	fmt.Println(key, start, step, data)
	fmt.Println(len(data))
	return ts.NewTimeSeriesOfData(key, start, step, data)
}
