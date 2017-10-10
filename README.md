# go-statsd-client

Client for StatsD (Golang)

[![Go Report Card](https://goreportcard.com/badge/github.com/sokil/go-statsd-client)](https://goreportcard.com/report/github.com/sokil/go-statsd-client)
[![GoDoc](https://godoc.org/github.com/sokil/go-statsd-client?status.svg)](https://godoc.org/github.com/sokil/go-statsd-client)
[![Build Status](https://travis-ci.org/sokil/go-statsd-client.svg?branch=master)](https://travis-ci.org/sokil/go-statsd-client)

## Useage

Client may be in buffered and unbuffered mode.

In buffered mode adding metric only adds it to buffer. Then `client.Flush()` builds all metrics to
packed and sends them to StatsD server by one request.

In unbuffered mode each metric sends to StatsD immediately.

Creating unbuffered client:

```go
client := NewClient("127.0.0.1", 9876)  # create client
client.Open()                           # open connection to StatsD
client.Count('a.b.c', 42, 0.7)          # set count metric and send it to StatsD
```

Creating buffered client:

```go
client := NewClient("127.0.0.1", 9876)  # create client
client.Open()                           # open connection to StatsD
client.Count('a.b.c', 42, 0.7)          # set count metric and add it to buffer
client.Count('a.b.d', 43)               # set timing metric and add it to buffer
client.Flush()							# send all metrics as one packet to StatsD
```
