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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type RSSEntry struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type RSSData struct {
	Entries []RSSEntry `json:"entries"`
}

const jsonFileName = "rss_urls.json"

func main() {
	http.HandleFunc("/api/rss", routeRss)

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
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func getRss(w http.ResponseWriter, r *http.Request) {
	// Load data from file
	rssData, err := loadRSSData()
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
	rssData, err := loadRSSData()
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
	if err := saveRSSData(rssData); err != nil {
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

func loadRSSData() (*RSSData, error) {
	// Check if file exists, if not return empty data
	if _, err := os.Stat(jsonFileName); os.IsNotExist(err) {
		return &RSSData{Entries: []RSSEntry{}}, nil
	}

	// Read the file
	data, err := os.ReadFile(jsonFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %v", err)
	}

	// Parse JSON
	var rssData RSSData
	if err := json.Unmarshal(data, &rssData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &rssData, nil
}

func saveRSSData(data *RSSData) error {
	// Convert to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(jsonFileName, jsonData, 0644); err != nil {
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
