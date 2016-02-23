// Copyright (c) 2014 Datacratic. All rights reserved.

package graphite

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var start = time.Date(2016, time.Month(1), 15, 17, 0, 0, 0, time.UTC)
var end = start.Add(7 * time.Minute)

func MockGraphite() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		fmt.Println(r.URL.Query())
		target := r.URL.Query()["target"][0]

		switch target {

		case "some.random.key":
			fmt.Fprintf(w, "some.random.key,%d,%d,%d|1,None,1.5,2,3,None,4\n",
				start.Unix(), end.Unix(), 60)

		case "some.*.key":
			fmt.Fprintf(w, `some.random.key,%d,%d,%d|1,None,1.5,2,3,None,4
some.other.key,%d,%d,%d|None,1,3,2,5,2\n`,
				start.Unix(), end.Unix(), 60, start.Unix(), end.Unix(), 60)

		default:
			panic("should not get here")
		}
	})

	return httptest.NewServer(handler)
}

func TestGraphiteHTTP(t *testing.T) {

	mock := MockGraphite()

	graphite := Graphite{
		URL: mock.URL,
	}
	graphite.Init()

	resp := graphite.Do(&Request{Key: "some.random.key"})

	_, err := resp.Single()
	if err != nil {
		fmt.Println(string(resp.Body))
		t.Fatal(err)
	}

	resp = graphite.Do(&Request{Key: "some.*.key"})

	_, err = resp.First()
	if err != nil {
		t.Fatal(err)
	}

	_, err = resp.All()
	if err != nil {
		t.Fatal(err)
	}
}
