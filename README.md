# HTTP

## Install

go1.12.5+ required

```bash
go get github.com/fedebonisconti/tbd
```

## Usage

```bash
$ tbd --help                                                                                                                                                                                                                                                            130 ↵
Usage of tbd:
  -headers string
    	Comma separated headers without (example: "Auth-Token:123,Accept:text/html,Content-Type:application/json")
  -input string
    	InputFile (default "stdin")
  -output string
    	Output inputFile name (default "stdout")
  -parallel int
    	Parallelism factor (default 4)
  -sample
    	Takes random samples from input to send requests (default false)
  -sample-size int
    	Sample size (if zero, sample is disabled)
  -wait int
    	Time between parallel requests in millis
    	
$ tbd -headers Auth-Token:123,Accept:text/html,Content-Type:application/json -sample -sample-size 10 -parallel 10 -input test.txt -output out.txt
```
