package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Header struct {
	Name  string
	Value string
}

type Arguments struct {
	parallelismFactor int
	timeBetweenBatch  int
	sample            bool
	sampleSize        int
	inputFile         *bufio.Reader
	outputFile        *bufio.Writer
	headers           []Header
	method            string
}

type ArgumentsReader struct {
}

type ArgumentsParser interface {
	Parse() *Arguments
}

func (argumentsParser* ArgumentsReader) Parse() *Arguments {
	inputFileName, outputFileName, arguments := parseArguments()

	var err error
	inputFile := os.Stdin
	if *inputFileName != "" {
		inputFile, err = os.Open(*inputFileName)
		checkError(err)
	}
	arguments.inputFile = bufio.NewReader(inputFile)
	outputFile := os.Stdout
	if *outputFileName != "" {
		outputFile, err = os.Create(*outputFileName)
		checkError(err)
	}
	arguments.outputFile = bufio.NewWriter(outputFile)

	return arguments
}

func parseArguments() (*string, *string, *Arguments) {
	inputFileName := flag.String("input", "", "Input file path (default stdin)")
	parallelismFactor := flag.Int("parallel", runtime.NumCPU(), "Parallelism factor")
	outputFileName := flag.String("output", "", "Output file path (default stdout)")
	timeBetweenBatch := flag.Int("wait", 0, "Time between batch of parallel requests in millis")
	sample := flag.Bool("sample-mode", false, "Takes random samples from input to send requests (default false)")
	sampleSize := flag.Int("sample-size", 0, "Sample size. If zero, sample is disabled")
	headers := flag.String("headers", "", "Comma separated headers without (example: \"Auth-Token:123,Accept:text/html,Content-Type:application/json\")")
	method := flag.String("method", "GET", "Http method to be used in every request.")
	flag.Parse()
	m := strings.ToUpper(*method)
	if m != "GET" && m != "POST" && m != "PUT" && m != "PATCH" && m != "DELETE" {
		flag.Usage()
		os.Exit(1)
	}
	arguments := Arguments{
		parallelismFactor: *parallelismFactor,
		timeBetweenBatch:  *timeBetweenBatch,
		sample:            *sample && (*sampleSize != 0),
		sampleSize:        *sampleSize,
		headers:           parseHeaders(*headers),
		method:            m,
	}
	return inputFileName, outputFileName, &arguments
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseHeaders(headers string) []Header {
	r := make([]Header, 0)
	for _, h := range strings.SplitN(headers, ",", -1) {
		if h != "" {
			s := strings.Split(h, ":")
			r = append(r, Header{Name: s[0], Value: s[1]})
		}
	}
	fmt.Println(fmt.Sprintf("Headers for requests: %v", r))
	return r
}
