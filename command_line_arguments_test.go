package main

import (
	"os"
	"runtime"
	"testing"
)

func TestParseHeadersThenHeadersShouldBeReturned(t *testing.T) {
	headersString := "Accept:text,Content-Type:application/json,Auth-Token:test"
	headers := parseHeaders(headersString)
	v := len(headers) == 3
	v = v && headers[0].Name == "Accept" && headers[0].Value == "text"
	v = v && headers[1].Name == "Content-Type" && headers[1].Value == "application/json"
	v = v && headers[2].Name == "Auth-Token" && headers[2].Value == "test"
	if !v {
		t.Fail()
	}
}

func TestWhenNoHeadersThenAnEmptyArrayShouldReturned(t *testing.T) {
	headersString := ""
	headers := parseHeaders(headersString)
	if len(headers) > 0 {
		t.Fail()
	}
}

func TestGetCLArgumentsThenDefaultArgumentsShouldBeReturned(t *testing.T) {
	arguments := GetCommandLineArguments()
	v := arguments.inputFile == os.Stdin
	v = v && arguments.outputFile == os.Stdout
	v = v && arguments.parallelismFactor == runtime.NumCPU()
	v = v && arguments.timeBetweenBatch == 0
	v = v && arguments.sample == false
	v = v && arguments.sampleSize == 0
	v = v && len(arguments.headers) == 0
	if !v {
		t.Fail()
	}
}
