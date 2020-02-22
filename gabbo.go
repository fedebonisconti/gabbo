package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var writeMutex sync.Mutex

type Response struct {
	response *http.Response
	elapsed  time.Duration
	err      error
}

type Gabbo struct {
}

func (g Gabbo) run() {
	arguments := GetCommandLineArguments()
	defer func() {
		checkError(arguments.inputFile.Close())
		checkError(arguments.outputFile.Close())
	}()

	var requestsWg sync.WaitGroup
	var responsesWg sync.WaitGroup

	responsesChannel := make(chan *Response)

	responsesWg.Add(1)
	go processResponses(responsesChannel, &responsesWg, arguments.outputFile)

	sent := 0

	iterable := getIterable(arguments)
	for iterable.HasNext() {
		requestsWg.Add(1)
		go doRequest("GET", strings.TrimSpace(iterable.Next()), arguments.headers, responsesChannel, &requestsWg)
		sent++
		if sent%arguments.parallelismFactor == 0 {
			requestsWg.Wait()
			time.Sleep(time.Duration(arguments.timeBetweenBatch) * time.Millisecond)
		}
		if arguments.sample && sent == arguments.sampleSize {
			break
		}
	}
	if sent%arguments.parallelismFactor > 0 {
		requestsWg.Wait()
	}
	close(responsesChannel)
	responsesWg.Wait()
}

func doRequest(method string, url string, headers []Header, responsesChannel chan<- *Response, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request", err.Error())
		return
	}
	for _, header := range headers {
		req.Header.Add(header.Name, header.Value)
	}
	start := time.Now()
	resp, err := client.Do(req)
	response := Response{
		response: resp,
		elapsed:  time.Since(start),
	}
	if err == nil {
		responsesChannel <- &response
	} else {
		fmt.Println(err)
	}
}

func processResponses(responsesChannel <-chan *Response, responsesWg *sync.WaitGroup, outputFile *os.File) {
	defer responsesWg.Done()
	timesElapsed := make([]time.Duration, 0)
	statuses := make([]int, 0)
	for {
		response, ok := <-responsesChannel
		if !ok {
			break
		}
		timesElapsed = append(timesElapsed, response.elapsed)
		statuses = append(statuses, response.statusCode())
		if response.statusCode() < 299 {
			writeFile(fmt.Sprintf("%s\n", response.bodyString()), outputFile)
		}
	}
	var sumElapsed int64 = 0
	var max int64 = 0
	var min int64 = math.MaxInt64
	for _, t := range timesElapsed {
		_t := int64(t) / (int64(time.Millisecond))
		sumElapsed += _t
		if _t > max {
			max = _t
		}
		if _t < min {
			min = _t
		}
	}
	success := 0
	redirects := 0
	clientErrors := 0
	serverErrors := 0
	for _, s := range statuses {
		if s >= 200 && s <= 299 {
			success++
		} else if s >= 300 && s <= 399 {
			redirects++
		} else if s >= 400 && s <= 499 {
			clientErrors++
		} else if s >= 500 && s <= 599 {
			serverErrors++
		}
	}
	if len(timesElapsed) == 0 || max == 0 || min == 0 {
		fmt.Println("Couldn't print statistics")
		return
	}
	fmt.Println("")
	fmt.Println(fmt.Sprintf("\tMinimum response time %d ms", min))
	fmt.Println(fmt.Sprintf("\tAverage time between requests %d ms", sumElapsed/int64(len(timesElapsed))))
	fmt.Println(fmt.Sprintf("\tMaximum response time %d ms", max))
	fmt.Println(fmt.Sprintf("\t%-10v|%10v", "Status","Count"))
	fmt.Println(fmt.Sprintf("\t%-10v|%10v", "2xx", success))
	fmt.Println(fmt.Sprintf("\t%-10v|%10v", "3xx", redirects))
	fmt.Println(fmt.Sprintf("\t%-10v|%10v", "4xx", clientErrors))
	fmt.Println(fmt.Sprintf("\t%-10v|%10v", "5xx", serverErrors))
	fmt.Println("")
}

func writeFile(content string, outputFile *os.File) {
	writeMutex.Lock()
	outputFile.WriteString(content)
	writeMutex.Unlock()
}

func getIterable(arguments *Arguments) Iterable {
	scanner := bufio.NewScanner(arguments.inputFile)
	if arguments.sample {
		v := make([]string, 0)
		for scanner.Scan() {
			v = append(v, scanner.Text())
		}
		return &RandomSliceIterator{slice: v}
	} else {
		return &ScannerIterator{scanner: scanner}
	}
}

func (r *Response) statusCode() int {
	return r.response.StatusCode
}

func (r *Response) bodyString() string {
	defer r.response.Body.Close()
	body, e := ioutil.ReadAll(r.response.Body)
	if e != nil {
		return ""
	}
	return string(body)
}

