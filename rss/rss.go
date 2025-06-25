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

func UpdateArticleReadStatus(w http.ResponseWriter, r *http.Request) {
	// Take in the ID of an article
	// Take in the status of the article bool. Read = true, unread = false
	// Query would look like this api/article?id=1&read=true

	// Call database function that would update the database entry

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

// TODO: Add UpdateArticleReadStatus(w http.ResponseWriter, r *http.Request) handler to mark articles as read/unread
// TODO: Add GetAllArticles(w http.ResponseWriter, r *http.Request) handler for paginated article listing across all feeds
// TODO: Add SearchArticles(w http.ResponseWriter, r *http.Request) handler for article search functionality
// TODO: Add UpdateRSS(w http.ResponseWriter, r *http.Request) handler to update RSS feed settings
// TODO: Add GetRSSStats(w http.ResponseWriter, r *http.Request) handler to return feed statistics
