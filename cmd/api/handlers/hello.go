package handlers

import "net/http"

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte("Hello gopher"))
	return err
}
