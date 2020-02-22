package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestRequestIsDoneThenItShouldBeInTheChannel(t *testing.T) {
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

func TestWhenTwoArePushedInTheChannelThenBothResponsesShouldBeProcessed(t *testing.T) {
	outputFileName := "test_responses.txt"
	f, _ := os.Create(outputFileName)
	defer func() {
		checkError(os.Remove(outputFileName))
	}()
	numberResponses := 4
	c := make(chan *Response, numberResponses)
	fillChannelWithResponses(numberResponses, c)
	var wg sync.WaitGroup
	wg.Add(1)
	go processResponses(c, &wg, f)
	close(c)
	// If this is wrong it will hang
	wg.Wait()
}

func TestWhenGetResponseBodyThenBodyShouldBeReturned(t *testing.T) {
	body := "Test"
	response := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       &http.Request{},
		Header:        make(http.Header, 0),
	}
	r := Response{response: response}
	if r.bodyString() != body {
		t.Fail()
	}
}

func TestWhenGetStatusCodeThenTeapotShouldBeReturned(t *testing.T) {
	response := &http.Response{StatusCode: http.StatusTeapot}
	r := Response{response: response}
	if r.statusCode() != http.StatusTeapot {
		t.Fail()
	}
}

func TestWhenGetIterableAndSampleIsEnabledThenARandomSliceIteratorShouldBeReturned(t *testing.T) {
	arguments := Arguments{sample: true, sampleSize: 10, inputFile: os.Stdin}
	i := getIterable(&arguments)
	if !isInstanceOf(i, (*RandomSliceIterator)(nil)) {
		t.Fail()
	}
}

func TestWhenGetIterableAndSampleIsDisabledThenAScannerIteratorShouldBeReturned(t *testing.T) {
	arguments := Arguments{sample: false, sampleSize: 10, inputFile: os.Stdin}
	i := getIterable(&arguments)
	if !isInstanceOf(i, (*ScannerIterator)(nil)) {
		t.Fail()
	}
}

func isInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func fillChannelWithResponses(q int, c chan<- *Response) {
	status := 200
	for i := 0; i < q; i++ {
		response := httptest.ResponseRecorder{Code: status}
		response.WriteString(fmt.Sprintf("%d", i))
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		c <- &Response{ response: response.Result(), elapsed: time.Since(start) }
		status += 100
	}
}