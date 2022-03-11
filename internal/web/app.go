package web

import (
	"log"
	"net/http"

	"github.com/fbuster/internal/handler"
	"github.com/fbuster/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	logger *log.Logger
	router *chi.Mux
	mw     []middleware.Middleware
}

func NewApp(logger *log.Logger, mw ...middleware.Middleware) *App {
	return &App{
		router: chi.NewRouter(),
		logger: logger,
		mw:     mw,
	}
}

func (a App) Handle(method, pattern string, h handler.Handler, mw ...middleware.Middleware) {
	h = middleware.Chain(h, a.mw...)
	h = middleware.Chain(h, mw...)

	f := func(w http.ResponseWriter, r *http.Request) {
		err := h.ServeHTTP(w, r)

		if err != nil {
			a.logger.Println(err)
		}
	}
	a.router.MethodFunc(method, pattern, f)
}

func (a App) HandleFunc(method, pattern string, h handler.HandlerFunc, mw ...middleware.Middleware) {
	a.Handle(method, pattern, h, mw...)
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
