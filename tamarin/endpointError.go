package tamarin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// EndpointError wraps a conventional error with additional content for return codes and specific response messages
type EndpointError struct {
	error
	returnCode    int
	returnMessage string
}

// FailWithErrorMessage is a convenience func for building an EndpointError with just a message
func FailWithErrorMessage(code int, message string, err error) *EndpointError {
	return &EndpointError{error: err, returnCode: code, returnMessage: message}
}

// FailWithJSONStatus is a convenience func for building an EndpointError with a JSON response body
func FailWithJSONStatus(code int, v any, err error) *EndpointError {
	jsonBytes, jErr := json.Marshal(v)
	if jErr != nil {
		err = fmt.Errorf("failed to marshal response JSON : %v. Original error was : %v", jErr, err)
		code = http.StatusInternalServerError
	}
	return &EndpointError{error: err, returnCode: code, returnMessage: string(jsonBytes)}
}

// SuceedWithJSONStatus returns a 200 and JSON-marshalled response body
func SuceedWithJSONStatus(responseBody any, rw http.ResponseWriter) *EndpointError {
	if responseBody == nil || rw == nil {
		return FailWithErrorMessage(http.StatusInternalServerError, "Internal Server Error", errors.New("responseBody or Response Writer was nil"))
	}
	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		return FailWithErrorMessage(http.StatusInternalServerError, "Internal Server Error", fmt.Errorf("unable to marshal response body : %v", err))
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
	return nil
}

// SuceedWithMessage returns a 200 and string response body
func SuceedWithMessage(responseMessage string, rw http.ResponseWriter) *EndpointError {
	if rw == nil {
		return FailWithErrorMessage(http.StatusInternalServerError, "Internal Server Error", errors.New("Response Writer was nil"))
	}
	rw.Header().Add("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(responseMessage))
	return nil
}
