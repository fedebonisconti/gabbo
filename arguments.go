package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	STDIN  string = "stdin"
	STDOUT string = "stdout"
)

type Parameters struct {
	inputFile         *os.File
	outputFile        *os.File
	parallelismFactor int
	timeBetweenBatch  int
	sample            bool
}

type Header struct {
	Name  string
	Value string
}

type Arguments struct {
	parallelismFactor int
	timeBetweenBatch  int
	sample            bool
	sampleSize        int
	inputFile         *os.File
	outputFile        *os.File
	headers           []Header
}

func (p *Parameters) Close() {
	checkError(p.inputFile.Close())
	checkError(p.outputFile.Close())
}

func GetArguments() *Arguments {
	inputFileName, outputFileName, arguments := parseArguments()

	var err error
	arguments.inputFile = os.Stdin
	if *inputFileName != "" && *inputFileName != STDIN {
		arguments.inputFile, err = os.Open(*inputFileName)
		checkError(err)
	}
	arguments.outputFile = os.Stdout
	if *outputFileName != "" && *outputFileName != STDOUT {
		arguments.outputFile, err = os.Create(*outputFileName)
		checkError(err)
	}

	return arguments
}

func parseArguments() (*string, *string, *Arguments) {
	inputFileName := flag.String("input", STDIN, "InputFile")
	parallelismFactor := flag.Int("parallel", runtime.NumCPU(), "Parallelism factor")
	outputFileName := flag.String("output", STDOUT, "Output inputFile name")
	timeBetweenBatch := flag.Int("wait", 0, "Time between batches in millis")
	sample := flag.Bool("sample", false, "Takes random samples from input to send requests")
	sampleSize := flag.Int("sample-size", 0, "Sample size (if zero, sample is disabled)")
	headers := flag.String("headers", "", "Comma separated headers without (example: \"Auth-Token:123,Accept:text/html,Content-Type:application/json\")")
	flag.Parse()
	arguments := Arguments{
		parallelismFactor: *parallelismFactor,
		timeBetweenBatch:  *timeBetweenBatch,
		sample:            *sample && (*sampleSize != 0),
		sampleSize:        *sampleSize,
		headers:           parseHeaders(*headers),
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
