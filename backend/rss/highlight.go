package rss

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JonSchaeffer/go-reader/db"
)

// GET /api/highlights?type=library&id=5
func GetHighlights(w http.ResponseWriter, r *http.Request) {
	itemType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")

	if itemType == "" || idStr == "" {
		http.Error(w, "type and id required", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	highlights, err := db.GetHighlightsByItem(itemType, itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if highlights == nil {
		highlights = []db.Highlight{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(highlights)
}

// POST /api/highlights
func CreateHighlight(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ItemType     string `json:"itemType"`
		ItemID       int    `json:"itemId"`
		SelectedText string `json:"selectedText"`
		Color        string `json:"color"`
		Note         string `json:"note"`
		TextOffset   int    `json:"textOffset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.SelectedText == "" || body.ItemID == 0 {
		http.Error(w, "selectedText and itemId required", http.StatusBadRequest)
		return
	}
	if body.Color == "" {
		body.Color = "yellow"
	}

	h, err := db.CreateHighlight(body.ItemType, body.ItemID, body.SelectedText, body.Color, body.Note, body.TextOffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h)
}

// PUT /api/highlights?id=3
func UpdateHighlight(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Color string `json:"color"`
		Note  string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := db.UpdateHighlight(id, body.Color, body.Note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/highlights?id=3
func DeleteHighlight(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.DeleteHighlight(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
