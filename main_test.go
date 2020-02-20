package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	// Start a local HTTP server
	requestCount := 0
	endpoint := "/some/path"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() != endpoint {
			t.Fail()
		}
		requestCount++
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()
	c := make(chan *Response, 1)
	defer close(c)
	var wg sync.WaitGroup
	wg.Add(1)
	go doRequest("GET", fmt.Sprintf("%s%s", server.URL, endpoint), make([]Header, 0), c, &wg)
	wg.Wait()
	response, ok := <-c
	if !ok || response.statusCode() != 200 {
		t.Fail()
	}
}

func TestXXXThenYYY(t *testing.T) {
	outputFileName := "test_responses.txt"
	f, _ := os.Create(outputFileName)
	defer func() {
		e := os.Remove(outputFileName)
		if e != nil {

		}
	}()
	numberResponses := 2
	c := make(chan *Response, numberResponses)
	fillChannelWithResponses(numberResponses, c)
	var wg sync.WaitGroup
	wg.Add(1)
	go processResponses(c, &wg, f)
	close(c)
	wg.Wait()
}

func fillChannelWithResponses(q int, c chan<- *Response) {
	for i := 0; i < q; i++ {
		response := httptest.ResponseRecorder{Code: http.StatusOK}
		response.WriteString(fmt.Sprintf("%d", i))
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		c <- &Response{ response: response.Result(), elapsed: time.Since(start) }
	}
}

func TestXXXThenYY(t *testing.T) {
	v := make([]time.Duration, 0)
	for i := 0; i < 5; i++ {
		r := rand.Intn(10)
		println(r)
		start := time.Now()
		time.Sleep(time.Duration(r) * time.Millisecond)
		v = append(v, time.Since(start))
	}
	j := int64(0)
	for _, t := range v {
		j += int64(t / time.Millisecond)
	}
	println(j / int64(len(v)))
}