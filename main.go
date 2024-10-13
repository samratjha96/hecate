package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samratjha96/hecate/internal/database"
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

	v1Router := chi.NewRouter()

	r.Mount("/v1", v1Router)

	v1Router.Route("/subreddits", func(router chi.Router) {

		router.Get("/", subredditGetHandler(db))
		router.Get("/{subredditName}", subredditPostsGetHandler(db))
		router.Post("/ingest", ingestSubredditHandler(db))
		router.Post("/ingest-all", ingestAllSubredditsHandler(db))
	})

	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		fmt.Println("Using default port 8000")
		serverPort = "8000"
	}
	serverAddr := fmt.Sprintf(":%s", serverPort)
	fmt.Printf("Starting server on port %v", serverAddr)
	http.ListenAndServe(serverAddr, r)
}
