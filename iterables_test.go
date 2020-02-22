package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestScannerIterator_HasNext(t *testing.T) {
	input := "test"
	scannerIterator := ScannerIterator{scanner: bufio.NewScanner(strings.NewReader(input))}
	if !scannerIterator.HasNext() {
		t.Fail()
	}
}

func TestScannerIterator_HasNotNext(t *testing.T) {
	input := ""
	scannerIterator := ScannerIterator{scanner: bufio.NewScanner(strings.NewReader(input))}
	if scannerIterator.HasNext() {
		t.Fail()
	}
}

func TestScannerIterator_Next(t *testing.T) {
	input := "test\ninput\nfor\nthis"
	scannerIterator := ScannerIterator{scanner: bufio.NewScanner(strings.NewReader(input))}
	for scannerIterator.HasNext() {
		if scannerIterator.Next() == "" {
			t.Fail()
		}
	}

}

func TestRandomSliceIterator_Next(t *testing.T) {
	v := []string{"this", "is", "a", "test", "array"}
	randomSliceIterator := RandomSliceIterator{slice: v}
	i := 0
	for randomSliceIterator.HasNext() && i < len(v){
		contains := false
		next := randomSliceIterator.Next()
		for _, s := range v {
			if s == next {
				contains = true
			}
		}
		if !contains {
			t.Fail()
		}
		i++
	}
}

func TestRandomSliceIterator_HasNext(t *testing.T) {
	randomSliceIterator := RandomSliceIterator{slice: []string{"test"}}
	if !randomSliceIterator.HasNext() {
		t.Fail()
	}
}

func TestRandomSliceIterator_HasNotNext(t *testing.T) {
	randomSliceIterator := RandomSliceIterator{slice: []string{}}
	if randomSliceIterator.HasNext() {
		t.Fail()
	}
}