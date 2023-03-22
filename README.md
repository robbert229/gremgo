# gremgo

> This repository is a fork of `github.com/qasaur/gremgo` which is currently unmaintained.

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/qasaur/gremgo) [![Build Status](https://travis-ci.org/qasaur/gremgo.svg?branch=master)](https://travis-ci.org/qasaur/gremgo) [![Go Report Card](https://goreportcard.com/badge/github.com/qasaur/gremgo)](https://goreportcard.com/report/github.com/qasaur/gremgo)

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

	dialer := gremgo.NewDialer("ws://127.0.0.1:8182") // Returns a WebSocket dialer to connect to Gremlin Server
	g, err := gremgo.Dial(dialer, errs) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
    	return
	}
	res, err := g.Execute( // Sends a query to Gremlin Server with bindings
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

	dialer := gremgo.NewSecureDialer("127.0.0.1:8182", "username", "password") // Returns a WebSocket dialer to connect to Gremlin Server
	g, err := gremgo.Dial(dialer, errs) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
    	return
	}
	res, err := g.Execute( // Sends a query to Gremlin Server with bindings
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

License
==========

Copyright (c) 2016 Marcus Engvall

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
