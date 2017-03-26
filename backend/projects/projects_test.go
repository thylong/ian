package projects

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fatih/color"
)

func TestStatus(t *testing.T) {
	cases := []struct {
		Project        string
		StatusCode     int
		ExpectedResult string
	}{
		{"test", 200, fmt.Sprintf("test: %s", color.GreenString("OK"))},
		{"test", 500, fmt.Sprintf("test: %s", color.RedString("ERROR"))},
	}
	for _, tc := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.StatusCode)
		}))
		defer ts.Close()

		result := Status(tc.Project, ts.URL, "/status")
		if result != tc.ExpectedResult {
			t.Errorf("Status func returned wrong result: got %#v want %#v",
				result, tc.ExpectedResult)
		}
	}
}

func TestStats(t *testing.T) {
	cases := []struct {
		Project     string
		StatusCode  int
		Body        io.Reader
		ExpectedErr error
	}{
		{"test", 200, bytes.NewBuffer([]byte(`{"id":66674923,"name_url":"regexrace","open_issues":10}`)), nil},
		{"test", 500, bytes.NewBuffer([]byte(`{"id":66674923,"name":"regexrace","open_issues":10}`)), ErrStatsUnavailable},
		{"test", 200, bytes.NewBuffer([]byte(`{"id"66674923,"name":"regexrace","open_issues":10}`)), ErrJSONPayloadInvalidFormat},
	}
	for _, tc := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.StatusCode)
			fmt.Fprintln(w, tc.Body)
		}))
		defer ts.Close()

		_, err := Stats(tc.Project, ts.URL)
		if err != tc.ExpectedErr {
			t.Errorf("Stats func returned wrong error: got %#v want %#v",
				err, tc.ExpectedErr)
		}
	}
}
