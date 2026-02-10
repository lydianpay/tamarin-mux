<div align="center">

# Tamarin Mux

![tamarin.png](tamarin.png)

</div>

---

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/lydianpay/tamarin-mux)](https://goreportcard.com/report/lydianpay/tamarin-mux)
[![Maintainability](https://qlty.sh/gh/lydianpay/projects/tamarin-mux/maintainability.svg)](https://qlty.sh/gh/lydianpay/projects/tamarin-mux)
[![Code Coverage](https://qlty.sh/gh/lydianpay/projects/tamarin-mux/coverage.svg)](https://qlty.sh/gh/lydianpay/projects/tamarin-mux)
[![CodeQL](https://github.com/lydianpay/tamarin-mux/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/lydianpay/tamarin-mux/actions/workflows/github-code-scanning/codeql)

</div>

Written in Go ("Golang" for search engines) with zero external dependencies, this package implements a clean, 
non-bloated, HTTP request multiplexer.

---

## Installation & Usage
1. Once confirming you have [Go](https://go.dev/doc/install) installed, the command below will add
   `tamarin` as a dependency to your Go program.
```bash
go get -u github.com/lydianpay/tamarin-mux
```
2. Import the package into your code
```go
package main

import (
   "github.com/lydianpay/tamarin-mux"
)
```
3. Examples
* A Simple GET handler 
* Test it with `curl localhost/ping`  
```go
package main

import (
	"net/http"
	"github.com/lydianpay/tamarin-mux/tamarin"
)

func main() {
	handler := tamarin.NewHandler(true).WithHandleFuncs("/ping", http.MethodGet, pong)
	http.ListenAndServe(":80", handler)
}

func pong(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("pong"))
}
```
* An example with sequenced handlers 
* Test it with `curl localhost/goodguy --header 'X-I-Am-A-Good-Guy: yes' --header 'Content-Type: application/json' --data '{"SomeValue":1,"SomeID":"a1b2c3"}'`
```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lydianpay/tamarin-mux/tamarin"
)

const goodGuyKey = "X-I-Am-A-Good-Guy"

func main() {
	handler := tamarin.NewHandler(true).
		WithPostEndpoint(tamarin.NewEndpoint("/goodguy").WithHandlers(middleWare1, processRequest))
	http.ListenAndServe(":80", handler)
}

func middleWare1(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	if req.Header == nil {
		return tamarin.FailWithErrorMessage(
			http.StatusBadRequest,
			"missing header",
			errors.New("request had no header"))
	}
	goodGuy, OK := req.Header[goodGuyKey]
	if !OK {
		return tamarin.FailWithErrorMessage(
			http.StatusBadRequest,
			"malformed header",
			fmt.Errorf("Header missing goodGuyKey '%s'", goodGuyKey))
	}
	if len(goodGuy) != 1 || !strings.EqualFold(goodGuy[0], "yes") {
		return tamarin.FailWithErrorMessage(
			http.StatusBadRequest,
			"malformed header",
			fmt.Errorf("requestor did not assert goodness"))
	}
	return nil
}

func processRequest(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	if req.Body == nil {
		return tamarin.FailWithJSONStatus(
			http.StatusBadRequest,
			&ErrorBody{Message: "no body in request"},
			errors.New("passed middleware but no body"))
	}
	jb, err := io.ReadAll(req.Body)
	if err != nil {
		return tamarin.FailWithJSONStatus(
			http.StatusBadRequest,
			&ErrorBody{Message: "malformed request body"},
			fmt.Errorf("couldn't read request body : %v", err))
	}
	rb := RequestBody{}
	err = json.Unmarshal(jb, &rb)
	if err != nil {
		return tamarin.FailWithJSONStatus(
			http.StatusBadRequest,
			&ErrorBody{Message: "malformed request body"},
			fmt.Errorf("couldn't unmarshal request body : %v", err))
	}
	return tamarin.SuceedWithJSONStatus(
		&ResponseBody{SomeResponse: fmt.Sprintf("Good job. %d / %s", rb.SomeValue, rb.SomeID)}, rw)
}

type RequestBody struct {
	SomeValue int
	SomeID    string
}

type ResponseBody struct {
	SomeResponse string
}

type ErrorBody struct {
	Message string
}
```
