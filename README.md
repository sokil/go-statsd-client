# go-statsd-client

Client for StatsD (Golang)

[![Go Report Card](https://goreportcard.com/badge/github.com/sokil/go-statsd-client)](https://goreportcard.com/report/github.com/sokil/go-statsd-client)
[![GoDoc](https://godoc.org/github.com/sokil/go-statsd-client?status.svg)](https://godoc.org/github.com/sokil/go-statsd-client)
[![Build Status](https://travis-ci.org/sokil/go-statsd-client.svg?branch=master)](https://travis-ci.org/sokil/go-statsd-client)

## Useage

```go
client := NewClient("127.0.0.1", 9876)  # create client
client.SetAutoflush(true)               # if true send metric on each set. By default false
client.Open()                           # open connection to StatsD
client.Count('a.b.c', 42, 0.7)          # set count metric. If autoflush enabled, then send to statsd
```
