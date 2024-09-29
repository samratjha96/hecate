package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type SubscribePayload struct {
	Subscriptions []string `json:"subreddits"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Hecate"))
	})

	r.Route("/users/{userID}", func(r chi.Router) {
		r.Get("/", getUserHandler)
		r.Post("/subscribe", subscribeHandler)
	})

	http.ListenAndServe(":3000", r)
	// os.Exit(hecate.Main())
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	w.Write([]byte(fmt.Sprintf("hi user %v", userID)))
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	decoder := json.NewDecoder(r.Body)
	subreddits := SubscribePayload{}
	err := decoder.Decode(&subreddits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("hi user %v with payload %v", userID, subreddits.Subscriptions)))
}
