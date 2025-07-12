package rss

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JonSchaeffer/go-reader/db"
)

type RSSEntry struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type RSSData struct {
	Entries []RSSEntry `json:"entries"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Format      string `xml:"format"`
	Identifier  string `xml:"identifier"`
}

func GetRss(w http.ResponseWriter, r *http.Request) {
	rss, err := db.GetAllRSS()
	if err != nil {
		http.Error(w, "No data returned", http.StatusBadRequest)
		return
	}

	// Check for optional ID parameter
	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		// Return specific entry by ID
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		rss, err := db.GetRSSByID(id)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rss)
		return
	}

	// Return all entries
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rss)
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	// Get article from database
	article, err := db.GetAllArticles()
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Return article as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetArticlesByRSSID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("rssid")
	limitParam := r.URL.Query().Get("limit")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if limitParam == "" {
		limitParam = "100"
	}

	// Convert ID parameter to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	// Get article from database
	article, err := db.GetArticleByRSSID(id, limit)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Return article as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetSingleArticle(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Convert ID parameter to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	// Get article from database
	article, err := db.GetSingleArticle(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Return article as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func SearchArticles(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("query")
	limitParam := r.URL.Query().Get("limit")

	if queryParam == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	if limitParam == "" {
		limitParam = "20"
	}

	// Convert limit param to int
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	article, err := db.SearchArticles(queryParam, limit)
	if err != nil {
		http.Error(w, "No search results found", http.StatusNotFound)
		return
	}
	// Return article as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateArticleReadStatus(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	readParam := r.URL.Query().Get("read")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if readParam == "" {
		http.Error(w, "Read parameter is required", http.StatusBadRequest)
	}

	// Convert ID parameter to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	// Convert read parameter to boolean
	read, err := strconv.ParseBool(readParam)
	if err != nil {
		http.Error(w, "Invalid read parameter", http.StatusBadRequest)
		return
	}

	err = db.UpdateArticleReadStatus(id, read)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating read status for article %d", id), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Article %d read status set to %t", id, read)))
}

func PostRss(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var requestData struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate that URL is not empty
	if requestData.URL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	fiveURL := GetRSSFiveURL(requestData.URL)

	fiveResponse, err := http.Get(fiveURL)
	if err != nil {
		log.Printf("Error fetching RSS feed: %v", err)
		http.Error(w, "Failed to fetch RSS feed", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(fiveResponse.Body)
	if err != nil {
		log.Printf("Error reading RSS response: %v", err)
		http.Error(w, "Failed to read RSS feed", http.StatusInternalServerError)
		return
	}

	if fiveResponse.Body != nil {
		defer fiveResponse.Body.Close()
	}

	var rssURL RSS
	err = xml.Unmarshal(body, &rssURL)
	if err != nil {
		log.Printf("Error parsing RSS XML: %v", err)
		http.Error(w, "Failed to parse RSS feed", http.StatusInternalServerError)
		return
	}

	// Create DB Entry
	rss, err := db.CreateRSS(requestData.URL, fiveURL, rssURL.Channel.Title, rssURL.Channel.Description, 1, 1)
	if err != nil {
		log.Printf("Error creating RSS: %v", err)
	}

	// Get RSS Feed, and save it to the DB
	SaveRSSArticles(rss.FiveURL, rss.ID)

	// Return Success Response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "RSS URL added successfully",
		"id":      rss.ID,
		"url":     rss.URL,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteRSSbyID(w http.ResponseWriter, r *http.Request) {
	// Check for ID paramater
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if idParam != "" {
		// Return specific entry by ID
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		err = db.DeleteRSSByID(id)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry deleted successfully"))

	}
}

func GetRSSFiveURL(RSSUrl string) string {
	return fmt.Sprintf("http://fullfeedrss:80/makefulltextfeed.php?url=%s&max=4&links=preserve", RSSUrl)
}

func SaveRSSArticles(FeedURL string, FeedID int) {
	response, err := http.Get(FeedURL)
	if err != nil {
		log.Printf("Error fetching feed URL %s: %v", FeedURL, err)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading feed response from %s: %v", FeedURL, err)
		return
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Printf("Error parsing RSS XML from %s: %v", FeedURL, err)
		return
	}

	processor := NewContentProcessor()

	for _, item := range rss.Channel.Items {
		// TODO: This approach isn't super efficient. It will process the description
		// for every Rss Article even if it already exists. I think putting in a Check
		// to compare GUID's first should help with this.

		// Process description
		processedDescription := processor.ProcessContent(item.Description)

		_, err := db.CreateArticle(FeedID, item.Title, item.Link,
			item.GUID, processedDescription, item.PubDate,
			item.Format, item.Identifier, false)
		if err != nil {
			log.Printf("Error saving article '%s': %v", item.Title, err)
			continue
		}
		fmt.Printf("%+v saved successfully.\n", item.Title)
	}
}

func StartRSSFetcher(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("RSS fetcher started")

	// Run once immediately
	FetchNewArticles()

	for {
		select {
		case <-ctx.Done():
			log.Println("RSS fetcher stopping...")
			return
		case <-ticker.C:
			log.Println("Starting scheduled RSS fetch...")
			FetchNewArticles()
		}
	}
}

func FetchNewArticles() {
	log.Println("Starting to fetch new articles...")

	rss, err := db.GetAllRSS()
	if err != nil {
		log.Printf("ERROR: Failed to get RSS feeds: %v", err)
		return // Changed from log.Fatal to return
	}

	log.Printf("Processing %d RSS feeds", len(rss))

	for i, item := range rss {
		log.Printf("Processing feed %d/%d: %s", i+1, len(rss), item.FiveURL)

		// Wrap in a function to catch panics
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic processing RSS feed %s: %v", item.FiveURL, r)
				}
			}()

			SaveRSSArticles(item.FiveURL, item.ID)
		}()
	}

	log.Println("Finished fetching articles")
}

func UpdateRSS(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	urlParam := r.URL.Query().Get("url")
	feedSizeParam := r.URL.Query().Get("feedsize")
	syncParam := r.URL.Query().Get("sync")
	categoryIDParam := r.URL.Query().Get("categoryid")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if urlParam == "" && feedSizeParam == "" && syncParam == "" && categoryIDParam == "" {
		http.Error(w, "At least one parameter (url, feedsize, sync, categoryid) is required", http.StatusBadRequest)
		return
	}

	// Convert ID parameter to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	updatedFields := []string{}
	updatedValues := map[string]interface{}{}

	if urlParam != "" {
		err := db.UpdateRSS(id, "url", urlParam)
		if err != nil {
			http.Error(w, "Error updating RSS URL", http.StatusBadRequest)
			return
		}
		updatedFields = append(updatedFields, "url")
		updatedValues["url"] = urlParam
	}

	if feedSizeParam != "" {
		feedSize, err := strconv.Atoi(feedSizeParam)
		if err != nil {
			http.Error(w, "Invalid feed size parameter", http.StatusBadRequest)
			return
		}
		err = db.UpdateRSS(id, "feedsize", feedSize)
		if err != nil {
			http.Error(w, "Error updating RSS feed size", http.StatusBadRequest)
			return
		}
		updatedFields = append(updatedFields, "feedsize")
		updatedValues["feedsize"] = feedSize
	}

	if syncParam != "" {
		sync, err := strconv.Atoi(syncParam)
		if err != nil {
			http.Error(w, "Invalid sync parameter", http.StatusBadRequest)
			return
		}
		err = db.UpdateRSS(id, "sync", sync)
		if err != nil {
			http.Error(w, "Error updating RSS feed sync", http.StatusBadRequest)
			return
		}
		updatedFields = append(updatedFields, "sync")
		updatedValues["sync"] = sync
	}

	if categoryIDParam != "" {
		if categoryIDParam == "null" {
			// Set to NULL for uncategorized
			err = db.UpdateRSSCategoryID(id, nil)
			if err != nil {
				http.Error(w, "Error updating RSS feed category", http.StatusBadRequest)
				return
			}
			updatedFields = append(updatedFields, "categoryid")
			updatedValues["categoryid"] = nil
		} else {
			categoryID, err := strconv.Atoi(categoryIDParam)
			if err != nil {
				http.Error(w, "Invalid category ID parameter", http.StatusBadRequest)
				return
			}
			err = db.UpdateRSSCategoryID(id, &categoryID)
			if err != nil {
				http.Error(w, "Error updating RSS feed category", http.StatusBadRequest)
				return
			}
			updatedFields = append(updatedFields, "categoryid")
			updatedValues["categoryid"] = categoryID
		}
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"message":        fmt.Sprintf("RSS feed %d updated successfully", id),
		"id":             id,
		"updated_fields": updatedFields,
		"updated_values": updatedValues,
	}
	json.NewEncoder(w).Encode(response)
}

func GetRSSStats(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Convert ID parameter to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	// Get stats from database
	stats, err := db.GetRSSStats(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving stats for RSS feed %d: %v", id, err), http.StatusNotFound)
		return
	}

	// Return stats as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Convert ID parameter to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	// Delete article from database
	err = db.DeleteArticle(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting article %d: %v", id, err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Article %d deleted successfully", id)))
}

// Category management handlers

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := db.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func PostCategory(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if reqData.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	// Set default color if not provided
	if reqData.Color == "" {
		reqData.Color = "#3b82f6"
	}

	category, err := db.CreateCategory(reqData.Name, reqData.Color)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create category: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var reqData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if reqData.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	if reqData.Color == "" {
		reqData.Color = "#3b82f6"
	}

	err = db.UpdateCategory(id, reqData.Name, reqData.Color)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update category: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"message": fmt.Sprintf("Category %d updated successfully", id),
		"id":      id,
		"name":    reqData.Name,
		"color":   reqData.Color,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	err = db.DeleteCategoryByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete category: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Category %d deleted successfully", id)))
}
