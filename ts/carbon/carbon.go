// Copyright (c) 2014 Datacratic. All rights reserved.

package carbon

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"sync"

	"github.com/datacratic/gotsvis/ts"
)

type Carbon struct {
	URL string

	once sync.Once
	conn net.Conn
	feed chan bytes.Buffer
}

func (carbon *Carbon) Init() {
	carbon.once.Do(carbon.init)
}

func (carbon *Carbon) init() {
	if carbon.URL == "" {
		panic("carbon URL can't be empty")
	}

	url, err := url.Parse(carbon.URL)
	if err != nil {
		panic(err)
	}

	carbon.conn, err = net.Dial("tcp", url.Host)
	if err != nil {
		panic(err)
	}
	log.Printf("carbon: connected to tcp://%s", url.Host)

	carbon.feed = make(chan bytes.Buffer)

	go carbon.Send()
}

func (carbon *Carbon) Send() {
	for {
		buffer := <-carbon.feed
		if _, err := io.Copy(carbon.conn, &buffer); err != nil {
			log.Printf("carbon.Send.error:", err)
		}
	}
}

func (carbon *Carbon) Write(ts *ts.TimeSeries) {
	carbon.Init()

	var buffer bytes.Buffer

	it := ts.IteratorTimeValue()
	for t, v, ok := it.Next(); ok; t, v, ok = it.Next() {
		fmt.Fprintf(&buffer, "%s %f %d\n", ts.Key(), v, t.Unix())
	}

	go func() {
		carbon.feed <- buffer
	}()
}

func (carbon *Carbon) WriteSlice(tss ts.TimeSeriesSlice) {
	for _, ts := range tss {
		carbon.Write(&ts)
	}
}
