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
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

const (
	jsonFileName = "rss_urls.json"
	jsonFeedName = "rss_feed.json"
)

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
	// Load data from file
	rssData, err := loadJSONFromFile[RSSData](jsonFileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading data: %v", err), http.StatusInternalServerError)
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

		for _, entry := range rssData.Entries {
			if entry.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(entry)
				return
			}
		}

		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Return all entries
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "RSS URLs retrieved successfully",
		"count":   len(rssData.Entries),
		"entries": rssData.Entries,
	})
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

	// Load Existing Data
	rssData, err := loadJSONFromFile[RSSData](jsonFileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading data: %v", err), http.StatusInternalServerError)
		return
	}

	// Create new entry
	newEntry := RSSEntry{
		ID:  getNextID(rssData.Entries),
		URL: requestData.URL,
	}

	// Add to Entries
	rssData.Entries = append(rssData.Entries, newEntry)

	// Save to file
	if err := saveJSONToFile(rssData, jsonFileName); err != nil {
		http.Error(w, fmt.Sprintf("Error saving data: %v", err), http.StatusInternalServerError)
		return
	}

	// Return Success Response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "RSS URL added successfully",
		"id":      newEntry.ID,
		"url":     newEntry.URL,
	}
	json.NewEncoder(w).Encode(response)
}

func loadJSONFromFile[T any](fileName string) (T, error) {
	var data T

	// Check if file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return data, fmt.Errorf("file does not exist: %s", fileName)
	}

	// Read the file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return data, fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal into the generic type
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return data, fmt.Errorf("error parsing JSON: %v", err)
	}

	return data, nil
}

func saveJSONToFile[T any](data T, fileName string) error {
	// Convert to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

// TODO: Rewrite this fucntion to get the last entry, get the ID, the add 1.
func getNextID(entries []RSSEntry) int {
	maxID := 0
	for _, entry := range entries {
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}
	return maxID + 1
}

func getRSSIDbyURL(entries []RSSEntry, url string) *RSSEntry {
	for i := range entries {
		if entries[i].URL == url {
			return &entries[i]
		}
	}
	return nil
}

func getRSSURLbyID(entries []RSSEntry, id int) *RSSEntry {
	for i := range entries {
		if entries[i].ID == id {
			return &entries[i]
		}
	}
	return nil
}

func getRSSURLbyContains(entries []RSSEntry, url string) []RSSEntry {
	var results []RSSEntry
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.URL), strings.ToLower(url)) {
			results = append(results, entry)
		}
	}
	return results
}

// Load the RSS data, pop the relevant entry, write to file

func deleteRSSbyID(w http.ResponseWriter, r *http.Request) {
	// Load data from file
	rssData, err := loadJSONFromFile[RSSData](jsonFileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loadting data: %v", err), http.StatusInternalServerError)
		return
	}

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

		found := false
		for i, entry := range rssData.Entries {
			if entry.ID == id {
				rssData.Entries = append(rssData.Entries[:i], rssData.Entries[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Entry not found", http.StatusNotFound)
			return
		}

		err = saveJSONToFile(rssData, jsonFileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error saving data: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry deleted successfully"))

	}
}

func getRSSFiveURL(RSSUrl string) string {
	return fmt.Sprintf("http://fullfeedrss:80/makefulltextfeed.php?url=%s&max=3&links=preserve", RSSUrl)
}

func saveRSSFeed(FeedURL string) {
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
	saveJSONToFile(rss, jsonFeedName)
	fmt.Printf("%+v\n", rss.Channel.Items[0].Title)
}
