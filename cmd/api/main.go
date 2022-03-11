package main

import (
	"log"
	"net/http"

	"github.com/fbuster/cmd/api/handlers"
	"github.com/fbuster/internal/middleware"
	"github.com/fbuster/internal/web"
)

func main() {
	logger := log.Default()
	app := web.NewApp(logger,
		middleware.Recovery(),
	)

	{
		app.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) error {
			_, err := w.Write([]byte("Hello Gophers!"))
			panic("Panic dans le monde")
			return err
		})
	}

	{
		user := handlers.User{Logger: logger}
		app.HandleFunc("GET", "/users", user.Find)
		app.HandleFunc("POST", "/users", user.Create)
	}

	server := http.Server{
		Handler: app,
		Addr:    ":3000",
	}

	server.ListenAndServe()
}
