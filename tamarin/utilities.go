package tamarin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// URLVarAtPosition returns the path segment at pos in a request URL path
func URLVarAtPosition(url string, pos int) string {
	split := strings.Split(url, "/")
	if len(split) <= pos {
		return ""
	}
	return split[pos]
}

func GetRequestBodyAndHeader(req *http.Request) ([]byte, http.Header, error) {
	if req == nil {
		return nil, nil, errors.New("nil request")
	}
	if req.Body == nil {
		return nil, nil, errors.New("nil request body")
	}
	if req.Header == nil {
		return nil, nil, errors.New("nil request header")
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read request body : %v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return bodyBytes, req.Header, nil
}

func UnmarshallJSONRequestBodyTo[T any](req *http.Request, target T) (*T, http.Header, error) {
	bodyBytes, header, err := GetRequestBodyAndHeader(req)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid request : %v", err)
	}
	err = json.Unmarshal(bodyBytes, &target)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal request body to target type : %v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return &target, header, nil
}
