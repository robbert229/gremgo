# gremgo

> This repository is a fork of `github.com/qasaur/gremgo` which is currently unmaintained.

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/qasaur/gremgo) [![Go](https://github.com/robbert229/gremgo/actions/workflows/go.yml/badge.svg)](https://github.com/robbert229/gremgo/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/qasaur/gremgo)](https://goreportcard.com/report/github.com/robbert229/gremgo)

gremgo is a fast, efficient, and easy-to-use client for the TinkerPop graph database stack. It is a Gremlin language
driver which uses WebSockets to interface with Gremlin Server and has a strong emphasis on concurrency and scalability.
Please keep in mind that gremgo is still under heavy development and although effort is being made to fully cover gremgo
with reliable tests, bugs may be present in several areas.

Installation
==========

```
go get github.com/robbert229/gremgo
```

Documentation
==========

* [pkg.go.dev](https://pkg.go.dev/github.com/robbert229/gremgo)

Example
==========

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/robbert229/gremgo"
)

func main() {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs) // Example of connection error handling logic

	dialer := gremgo.NewDialer("ws://127.0.0.1:8182")                // Returns a WebSocket dialer to connect to Gremlin Server
	g, err := gremgo.DialContext(context.Background(), dialer, errs) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := g.ExecuteContext( // Sends a query to Gremlin Server with bindings
		context.Background(),
		"g.V(x)",
		map[string]string{"x": "1234"},
		map[string]string{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
```

Authentication
==========
The plugin accepts authentication creating a secure dialer where credentials are setted.
If the server where are you trying to connect needs authentication, and you do not provide the
credentials the complement will panic.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/robbert229/gremgo"
)

func main() {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs) // Example of connection error handling logic

	dialer := gremgo.NewDialer(
		"127.0.0.1:8182",
		gremgo.WithAuthentication("username", "password"),
	) // Returns a WebSocket dialer to connect to Gremlin Server

	g, err := gremgo.DialContext(context.Background(), dialer, errs) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := g.ExecuteContext( // Sends a query to Gremlin Server with bindings
		context.Background(),
		"g.V(x)",
		map[string]string{"x": "1234"},
		map[string]string{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
```