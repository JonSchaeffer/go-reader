package main

//
//curl -X POST http://localhost:8080/api/rss \                                  k3s-homelab.porgy-monitor.ts.net
//> -H "Content-Type: application/json" \
//> -d '{"url": "https://example.com/rss"}'
//
// curl http://localhost:8080/api/rss
// curl http://localhost:8080/api/rss?id=1

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	PubDate     string `xml:"pubdate"`
	Format      string `xml:"dc:format"`
	Identifier  string `xml:"dc:identifier"`
}

func main() {
	http.HandleFunc("/api/rss", routeRss)

	// saveRSSFeed(getRSSFiveURL("https://news.ycombinator.com/rss"))
	err := db.Init("postgres://postgres:postgres@postgres:5432")
	if err != nil {
		log.Fatal(err)
	}
	err = db.CreateRSSTable()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateArticleTable()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
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

	// Create DB Entry
	rss, err := db.CreateRSS(requestData.URL, fiveURL, "Title", "description", 1, 1)
	if err != nil {
		log.Printf("Error creating RSS: %v", err)
	}

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
	return fmt.Sprintf("http://fullfeedrss:80/makefulltextfeed.php?url=%s&max=3&links=preserve", RSSUrl)
}
