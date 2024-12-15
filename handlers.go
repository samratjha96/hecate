package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/hecate"
)

const (
	statusOK       = http.StatusOK
	statusCreated  = http.StatusCreated
	statusBadReq   = http.StatusBadRequest
	statusIntError = http.StatusInternalServerError
)

// ingestSubredditHandler handles the ingestion of a single subreddit
func ingestSubredditHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var subreddit hecate.SubscribeFrontendRequest
		if err := decodeJSONBody(w, r, &subreddit); err != nil {
			log.Printf("Failed to decode request body: %v", err)
			return
		}

		log.Printf("Ingesting subreddit: %s", subreddit.Subreddit.Name)
		subscriptions, err := hecate.IngestSubreddit(r.Context(), db, subreddit.Subreddit)
		if err != nil {
			log.Printf("Failed to ingest subreddit %s: %v", subreddit.Subreddit.Name, err)
			respondWithError(w, statusIntError, fmt.Sprintf("Failed to ingest subreddit: %v", err))
			return
		}

		log.Printf("Successfully ingested subreddit: %s", subreddit.Subreddit.Name)
		respondWithJson(w, statusCreated, subscriptions)
	}
}

// ingestAllSubredditsHandler handles the ingestion of all subreddits
func ingestAllSubredditsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request hecate.IngestAllFrontendRequest
		if err := decodeJSONBody(w, r, &request); err != nil {
			log.Printf("Failed to decode request body: %v", err)
			return
		}

		log.Printf("Ingesting all subreddits with sort: %s", request.SortBy)
		if err := hecate.IngestAllSubreddit(r.Context(), db, request.SortBy); err != nil {
			log.Printf("Failed to ingest all subreddits: %v", err)
			respondWithError(w, statusIntError, fmt.Sprintf("Failed to ingest all subreddits: %v", err))
			return
		}

		log.Println("Successfully ingested all subreddits")
		respondWithJson(w, statusOK, map[string]string{"status": "success"})
	}
}

// subredditGetHandler handles retrieving all subreddits
func subredditGetHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Retrieving all subreddits")
		subreddits, err := hecate.GetAllSubreddits(db)
		if err != nil {
			log.Printf("Failed to retrieve subreddits: %v", err)
			respondWithError(w, statusIntError, fmt.Sprintf("Failed to retrieve subreddits: %v", err))
			return
		}
		log.Printf("Retrieved %d subreddits", len(subreddits))
		respondWithJson(w, statusOK, subreddits)
	}
}

// subredditPostsGetHandler handles retrieving posts for a specific subreddit
func subredditPostsGetHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subredditName := chi.URLParam(r, "subredditName")
		log.Printf("Retrieving posts for subreddit: %s", subredditName)
		posts, err := hecate.GetAllPostsForSubreddit(db, subredditName)
		if err != nil {
			log.Printf("Failed to retrieve posts for subreddit %s: %v", subredditName, err)
			respondWithError(w, statusIntError, fmt.Sprintf("Failed to retrieve posts: %v", err))
			return
		}
		log.Printf("Retrieved %d posts for subreddit: %s", len(posts), subredditName)
		respondWithJson(w, statusOK, posts)
	}
}

// searchPostsHandler handles searching posts across all subreddits
func searchPostsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			respondWithError(w, statusBadReq, "Search query is required")
			return
		}

		log.Printf("Searching posts with query: %s", query)
		posts, err := db.SearchPosts(query)
		if err != nil {
			log.Printf("Failed to search posts: %v", err)
			respondWithError(w, statusIntError, fmt.Sprintf("Failed to search posts: %v", err))
			return
		}

		response := hecate.SearchPostsResponse{
			Posts: make([]hecate.SubredditPostFrontendResponse, len(posts)),
		}
		for i, post := range posts {
			response.Posts[i] = hecate.SubredditPostFrontendResponse{
				Title:         post.Title,
				Content:       post.Content,
				DiscussionURL: post.DiscussionURL,
				CommentCount:  post.CommentCount,
				Upvotes:       post.Upvotes,
				SubredditName: post.SubredditName,
			}
		}

		log.Printf("Found %d posts matching query: %s", len(posts), query)
		respondWithJson(w, statusOK, response)
	}
}

// decodeJSONBody decodes the JSON body of a request into a given struct
func decodeJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		respondWithError(w, statusBadReq, fmt.Sprintf("Invalid request body: %v", err))
		return fmt.Errorf("failed to decode JSON body: %w", err)
	}
	return nil
}
