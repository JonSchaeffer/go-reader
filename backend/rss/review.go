package rss

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JonSchaeffer/go-reader/db"
)

// GET /api/review/queue?limit=5
func GetReviewQueue(w http.ResponseWriter, r *http.Request) {
	limit := 5
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	queue, err := db.GetReviewQueue(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queue)
}

// POST /api/review/complete?id=
func CompleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	if err := db.CompleteReview(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/review/count
func GetReviewCount(w http.ResponseWriter, r *http.Request) {
	count, err := db.GetReviewQueueSize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
