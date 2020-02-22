# Gabbo

![Alt Text](https://media.giphy.com/media/3orieQxWzEndtoNSxi/giphy.gif)


## What is Gabbo?
Gabbo is an ApacheBench CLI tool implementation written in Go.

## Install

go1.12.5+ required

```bash
go get github.com/fedebonisconti/gabbo
```

## Usage

```bash
$ gabbo --help                                                                                                                                                                                                                                                            130 ↵
Usage of gabbo:
  -headers string
    	Comma separated headers without (example: "Auth-Token:123,Accept:text/html,Content-Type:application/json")
  -input string
    	Input file path (default stdin)
  -output string
    	Output file path (default stdout)
  -parallel int
    	Parallelism factor (default 4)
  -sample-mode
    	Takes random samples from input to send requests (default false)
  -sample-size int
    	Sample size. If zero, sample is disabled
  -wait int
    	Time between batch of parallel requests in millis
    	
$ gabbo -headers Auth-Token:123,Accept:text/html,Content-Type:application/json -sample-mode -sample-size 1500 -parallel 10 -input test.txt -output out.txt
Headers for requests: [{Auth-Token 123}, {Accept text/html}, {Content-Type application/json}]


	Minimum response time 704 ms
	Average time between requests 1114 ms
	Maximum response time 7476 ms
	Status    |     Count
	2xx       |      1500
	3xx       |         0
	4xx       |         0
	5xx       |         0


```


## Contribution

1. Fork ([https://github.com/fedebonisconti/gabbo/fork](https://github.com/fedebonisconti/gabbo/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test` command and confirm that it passes
1. Create a new Pull Request