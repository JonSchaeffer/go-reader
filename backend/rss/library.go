package rss

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/JonSchaeffer/go-reader/db"
)

func SaveToLibrary(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		URL       string `json:"url"`
		ArticleID int    `json:"articleId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var item *db.SavedItem
	var err error

	if reqData.ArticleID > 0 {
		item, err = saveArticleToLibrary(reqData.ArticleID)
	} else if reqData.URL != "" {
		item, err = saveURLToLibrary(reqData.URL)
	} else {
		http.Error(w, "url or articleId is required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func saveArticleToLibrary(articleID int) (*db.SavedItem, error) {
	exists, err := db.SavedItemExistsBySourceID(articleID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("article already saved to library")
	}

	articles, err := db.GetSingleArticle(articleID)
	if err != nil || len(articles) == 0 {
		return nil, fmt.Errorf("article not found")
	}
	a := articles[0]

	return db.CreateSavedItem("rss_article", &a.ID, a.Link, a.Title, a.Author, a.Description, a.PublishDate)
}

func saveURLToLibrary(url string) (*db.SavedItem, error) {
	exists, err := db.SavedItemExistsByURL(url)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("URL already saved to library")
	}

	fiveURL := fmt.Sprintf("%s/makefulltextfeed.php?url=%s&max=1&links=preserve", config.FiveFiltersURL, url)

	resp, err := http.Get(fiveURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var rssData RSS
	if err := xml.Unmarshal(body, &rssData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	title := rssData.Channel.Title
	var content, pubDate string

	if len(rssData.Channel.Items) > 0 {
		item := rssData.Channel.Items[0]
		if item.Title != "" {
			title = item.Title
		}
		processor := NewContentProcessor()
		content = processor.ProcessContent(item.Description)
		pubDate = item.PubDate
	}

	return db.CreateSavedItem("url", nil, url, title, "", content, pubDate)
}

func GetLibraryItems(w http.ResponseWriter, r *http.Request) {
	offset, limit, archived := 0, 50, false

	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		offset = o
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 200 {
		limit = l
	}
	if r.URL.Query().Get("archived") == "true" {
		archived = true
	}

	items, err := db.GetAllSavedItems(offset, limit, archived)
	if err != nil {
		http.Error(w, "Failed to get library items", http.StatusInternalServerError)
		return
	}

	total, err := db.GetTotalSavedItemCount(archived)
	if err != nil {
		http.Error(w, "Failed to get count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":   items,
		"total":   total,
		"offset":  offset,
		"limit":   limit,
		"hasMore": offset+limit < total,
	})
}

func GetSingleLibraryItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	item, err := db.GetSavedItemByID(id)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func UpdateLibraryItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	var reqData struct {
		Read     bool `json:"read"`
		Archived bool `json:"archived"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.UpdateSavedItemStatus(id, reqData.Read, reqData.Archived); err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteLibraryItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	if err := db.DeleteSavedItem(id); err != nil {
		http.Error(w, "Failed to delete item", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SearchLibrary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "q parameter is required", http.StatusBadRequest)
		return
	}

	limit := 20
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}

	items, err := db.SearchSavedItems(query, limit)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
