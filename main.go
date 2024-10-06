package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/hecate"
)

func main() {
	// Initialize database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tables
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Hecate"))
	})

	r.Route("/subreddits", func(r chi.Router) {

		r.Get("/", subredditGetHandler(db))
		r.Get("/{subredditName}", subredditPostsGetHandler(db))
		r.Post("/ingest", subscribeHandler(db))
	})

	server := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	fmt.Printf("Starting server on server %v", server)
	http.ListenAndServe(server, r)
}

func subscribeHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		subreddit := hecate.SubscribeFrontendRequest{}
		err := decoder.Decode(&subreddit)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		subscriptions, err := hecate.IngestSubreddit(db, subreddit.Subscription)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJson(w, http.StatusCreated, subscriptions)
	}
}

func subredditGetHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subreddits, err := hecate.GetAllSubreddits(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusCreated, subreddits)
	}
}

func subredditPostsGetHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subredditName := chi.URLParam(r, "subredditName")
		posts, err := hecate.GetAllPostsForSubreddit(db, subredditName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, posts)
	}
}
