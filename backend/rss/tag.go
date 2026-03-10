package rss

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JonSchaeffer/go-reader/db"
)

// GET /api/tags
func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := db.GetAllTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// POST /api/tags  body: {"name":"..."}
func CreateTag(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	tag, err := db.CreateTag(body.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

// DELETE /api/tags?id=
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTag(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/tags/item?type=&id=
func GetItemTags(w http.ResponseWriter, r *http.Request) {
	itemType := r.URL.Query().Get("type")
	itemID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || itemType == "" {
		http.Error(w, "type and id required", http.StatusBadRequest)
		return
	}
	tags, err := db.GetTagsForItem(itemType, itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// POST /api/tags/item  body: {"itemType":"...","itemId":1,"tagId":2}
func AddItemTag(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ItemType string `json:"itemType"`
		ItemID   int    `json:"itemId"`
		TagID    int    `json:"tagId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ItemType == "" {
		http.Error(w, "itemType, itemId, tagId required", http.StatusBadRequest)
		return
	}
	if err := db.AddTagToItem(body.ItemType, body.ItemID, body.TagID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/tags/item?type=&id=&tagId=
func RemoveItemTag(w http.ResponseWriter, r *http.Request) {
	itemType := r.URL.Query().Get("type")
	itemID, err1 := strconv.Atoi(r.URL.Query().Get("id"))
	tagID, err2 := strconv.Atoi(r.URL.Query().Get("tagId"))
	if err1 != nil || err2 != nil || itemType == "" {
		http.Error(w, "type, id, tagId required", http.StatusBadRequest)
		return
	}
	if err := db.RemoveTagFromItem(itemType, itemID, tagID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
