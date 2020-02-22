package main

import (
	"testing"
)

func TestParseHeaders(t *testing.T) {
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
