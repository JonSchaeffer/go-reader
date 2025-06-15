package main

//
//curl -X POST http://localhost:8080/api/rss \                                  k3s-homelab.porgy-monitor.ts.net
//> -H "Content-Type: application/json" \
//> -d '{"url": "https://example.com/rss"}'
//
// curl http://localhost:8080/api/rss
// curl http://localhost:8080/api/rss?id=1

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

func main() {
	// Initialize database
	err := db.Init("postgres://postgres:postgres@postgres:5432")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.CreateRSSTable()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateArticleTable()
	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP routes
	http.HandleFunc("/api/rss", routeRss)

	// Start RSS fetcher in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go startRSSFetcher(ctx)

	// Handle shutdown signals in background
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
		os.Exit(0)
	}()

	// Start HTTP server (this blocks)
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Server failed to start: %v", err)
	}
}

func routeRss(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRss(w, r)
	case http.MethodPost:
		postRss(w, r)
	case http.MethodDelete:
		deleteRSSbyID(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func getRss(w http.ResponseWriter, r *http.Request) {
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

func postRss(w http.ResponseWriter, r *http.Request) {
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

	fiveURL := getRSSFiveURL(requestData.URL)

	fiveResponse, err := http.Get(fiveURL)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(fiveResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	if fiveResponse.Body != nil {
		defer fiveResponse.Body.Close()
	}

	var rssURL RSS
	err = xml.Unmarshal(body, &rssURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create DB Entry
	rss, err := db.CreateRSS(requestData.URL, fiveURL, rssURL.Channel.Title, rssURL.Channel.Description, 1, 1)
	if err != nil {
		log.Printf("Error creating RSS: %v", err)
	}

	// Get RSS Feed, and save it to the DB
	saveRSSFeed(rss.FiveURL, rss.ID)

	// Return Success Response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "RSS URL added successfully",
		"id":      rss.ID,
		"url":     rss.URL,
	}
	json.NewEncoder(w).Encode(response)
}

func deleteRSSbyID(w http.ResponseWriter, r *http.Request) {
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

func getRSSFiveURL(RSSUrl string) string {
	return fmt.Sprintf("http://fullfeedrss:80/makefulltextfeed.php?url=%s&max=4&links=preserve", RSSUrl)
}

func saveRSSFeed(FeedURL string, FeedID int) {
	response, err := http.Get(FeedURL)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range rss.Channel.Items {
		fmt.Printf("Data: %+v\n", item.Title)
		_, err := db.CreateArticle(FeedID, item.Title, item.Link,
			item.GUID, item.Description, item.PubDate,
			item.Format, item.Identifier, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v saved successfully.\n", item.Title)
	}
}

func startRSSFetcher(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("RSS fetcher started")

	// Run once immediately
	fetchNewArticles()

	for {
		select {
		case <-ctx.Done():
			log.Println("RSS fetcher stopping...")
			return
		case <-ticker.C:
			log.Println("Starting scheduled RSS fetch...")
			fetchNewArticles()
		}
	}
}

func fetchNewArticles() {
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

			saveRSSFeed(item.FiveURL, item.ID)
		}()
	}

	log.Println("Finished fetching articles")
}
