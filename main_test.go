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
	numberResponses := 2
	c := make(chan *Response, numberResponses)
	fillChannelWithResponses(numberResponses, c)
	var wg sync.WaitGroup
	wg.Add(1)
	go processResponses(c, &wg, f)
	close(c)
	// If this is wrong it will hang
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