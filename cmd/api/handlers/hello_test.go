package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fbuster/cmd/api/handlers"
)

type testWriter struct {
	http.ResponseWriter
}

func (w *testWriter) Write(b []byte) (int, error) {
	return 0, errors.New("write error")
}

var req *http.Request

func reset() {
	req = httptest.NewRequest("GET", "http://example.com", nil)
}

func init() {
	reset()
}

func Test_Hello(t *testing.T) {
	defer t.Cleanup(reset)

	// t.Run("ShouldTest", func(t *testing.T) {

	// })

	req := httptest.NewRequest("GET", "http://example.com", nil)
	rec := httptest.NewRecorder()

	hello := handlers.Hello{}
	err := hello.ServeHTTP(rec, req)

	if err != nil {
		t.Errorf("got %v; want %v", err, nil)
	}
	if got, want := rec.Body.String(), "Hello gopher"; got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}

func Test_HelloErrorWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &testWriter{ResponseWriter: rec}
	req := httptest.NewRequest("GET", "http://example.com", nil)

	// ww, ok := rw.ResponseWriter.(*httptest.ResponseRecorder)

	hello := handlers.Hello{}
	err := hello.ServeHTTP(rw, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("error expected")
	}
}
