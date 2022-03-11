package middleware

import "github.com/fbuster/internal/handler"

type Middleware func(handler.Handler) handler.Handler

// Middleware(Middleware((handler)))
func Chain(h handler.Handler, mw ...Middleware) handler.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
