package middleware

import (
	"errors"
	"net/http"

	"github.com/fbuster/internal/handler"
)

func Recovery() Middleware {

	f := func(next handler.Handler) handler.Handler {

		h := func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if r := recover(); r != nil {
					switch v := r.(type) {
					case string:
						err = errors.New(v)
					case error:
						err = v
					default:
						err = errors.New("unknown error")
					}
				}
			}()
			return next.ServeHTTP(w, r)
		}
		return handler.HandlerFunc(h)
	}

	return f
}
