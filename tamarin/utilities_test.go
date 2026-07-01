package tamarin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestGetRequestBodyAndHeader(t *testing.T) {
	_, _, err := GetRequestBodyAndHeader(nil)
	if err == nil {
		t.Fail()
	}
	_, _, err = GetRequestBodyAndHeader(&http.Request{})
	if err == nil {
		t.Fail()
	}
	_, _, err = GetRequestBodyAndHeader(&http.Request{Body: io.NopCloser(bytes.NewReader([]byte{}))})
	if err == nil {
		t.Fail()
	}
	req, err := http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte("A")))
	if err != nil {
		t.Fail()
	}
	req.Header.Add("B", "C")
	body, header, err := GetRequestBodyAndHeader(req)
	if err != nil {
		t.Fail()
	}
	if !bytes.Equal(body, []byte("A")) {
		t.Fail()
	}
	if header.Get("B") != "C" {
		t.Fail()
	}
}

type testingStructA struct {
	S string
	F float64
}
type testingStructB struct {
	B bool
	I int
}

func TestUnmarshallJSONRequestBodyTo(t *testing.T) {
	tsA := testingStructA{S: "A", F: 1.23}
	bodyBytesA, err := json.Marshal(tsA)
	if err != nil {
		t.Fail()
	}
	req, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(bodyBytesA))
	if err != nil {
		t.Log("err request")
		t.Fail()
	}
	req.Header.Add("B", "C")

	result, _, err := UnmarshallJSONRequestBodyTo(req, testingStructA{})
	if err != nil {
		t.Log("err happy")
		t.Fail()
	}
	if result == nil || result.S != "A" {
		t.Logf("err result : %v", err)
		t.Fail()
	}
	_, _, err = UnmarshallJSONRequestBodyTo(req, "")
	if err == nil {
		t.Logf("err sad : %v", err)
		t.Fail()
	}
}

func TestURLVarAtPosition(t *testing.T) {
	cases := []struct {
		name string
		url  string
		pos  int
		want string
	}{
		{"variable at position 2", "/application/ct7wg-abc", 2, "ct7wg-abc"},
		{"static segment at position 1", "/application/ct7wg-abc", 1, "application"},
		{"leading empty at position 0", "/application/ct7wg-abc", 0, ""},
		{"trailing slash, missing var", "/application/", 2, ""},
		{"no var segment", "/application", 2, ""},
		{"position out of range", "/application", 9, ""},
		{"deep path", "/accounts/acc-1/keys/key-9", 4, "key-9"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := URLVarAtPosition(tc.url, tc.pos); got != tc.want {
				t.Errorf("URLVarAtPosition(%q, %d) = %q, want %q", tc.url, tc.pos, got, tc.want)
			}
		})
	}
}
