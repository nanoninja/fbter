package handlers

import (
	"log"
	"net/http"
)

type User struct {
	Logger *log.Logger
}

func (u User) Find(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u User) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}
