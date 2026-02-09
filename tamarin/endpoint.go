package tamarin

import (
	"encoding/json"
	"log"
	"net/http"
)

// EndpointHandlerFunc is a wrapper for http.HandlerFunc's that returns flavored Errors
type EndpointHandlerFunc func(http.ResponseWriter, *http.Request) *EndpointError

type endpoint struct {
	path     string
	method   string
	sequence []EndpointHandlerFunc
}

// NewEndpoint returns an instantiated Endpoint
func NewEndpoint(path string) *endpoint {
	return &endpoint{sequence: []EndpointHandlerFunc{}, path: path}
}

// WithMethod adds the HTTP Method (GET/POST) to the endpoint
func (e *endpoint) WithMethod(httpMethod string) *endpoint {
	e.method = httpMethod
	return e
}

// WithHandlers adds EndpointHandlerFunc's to the sequence of HandlerFuncs to be
// executed by this Endpoint
func (e *endpoint) WithHandlers(eFunc ...EndpointHandlerFunc) *endpoint {
	e.sequence = append(e.sequence, eFunc...)
	return e
}

// Handle satisfies the http.HandlerFunc interface but executes multiple wrapped
// HandlerFuncs in sequence, stopping if there is an error
func (e *endpoint) Handle(rw http.ResponseWriter, req *http.Request) {
	for _, f := range e.sequence {
		err := f(rw, req)
		if err != nil {
			if json.Valid([]byte(err.returnMessage)) {
				rw.Header().Add("Content-Type", "application/json")
			} else {
				rw.Header().Add("Content-Type", "text/plain")
			}
			rw.WriteHeader(err.returnCode)
			rw.Write([]byte(err.returnMessage))
			log.Printf("Stopping sequence for '%s' due to error : %v", e.path, err.error)
			log.Printf("User will see Error Code : %d / Message : %s", err.returnCode, err.returnMessage)
			break
		}
	}
}
