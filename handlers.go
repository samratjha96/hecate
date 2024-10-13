package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/hecate"
)

func ingestSubredditHandler(db *database.DB) http.HandlerFunc {
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
func ingestAllSubredditsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		request := hecate.IngestAllFrontendRequest{}
		err := decoder.Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = hecate.IngestAllSubreddit(db, request.SortBy)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJson(w, http.StatusOK, make(map[string]string))
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
