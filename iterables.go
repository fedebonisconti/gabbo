package main

import (
	"bufio"
	"math/rand"
)

type Iterable interface {
	HasNext() bool
	Next() string
}

type ScannerIterator struct {
	scanner *bufio.Scanner
}

func (s ScannerIterator) HasNext() bool {
	return s.scanner.Scan()
}

func (s ScannerIterator) Next() string {
	return s.scanner.Text()
}

type RandomSliceIterator struct {
	slice []string
}

func (r RandomSliceIterator) HasNext() bool {
	return len(r.slice) > 0
}

func (r RandomSliceIterator) Next() string {
	return r.slice[rand.Intn(len(r.slice))]
}