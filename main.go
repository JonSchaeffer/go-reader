package main

//
//curl -X POST http://localhost:8080/api/rss \                                  k3s-homelab.porgy-monitor.ts.net
//> -H "Content-Type: application/json" \
//> -d '{"url": "https://example.com/rss"}'
//
// curl http://localhost:8080/api/rss
// curl http://localhost:8080/api/rss?id=1

// TODO: Before saving the description to the database, process the html into something
// more readable :)

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/JonSchaeffer/go-reader/db"
	"github.com/JonSchaeffer/go-reader/rss"
)

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
	http.HandleFunc("/api/articles", routeAllArticles)           // All articles
	http.HandleFunc("/api/articles/single", routeSingleArticle)  // Single article by ?id=
	http.HandleFunc("/api/articles/by-rss", routeArticlesByRSS)  // Articles by RSS ID
	http.HandleFunc("/api/articles/update", routeUpdateArticle)  // Update article read status
	http.HandleFunc("/api/articles/search", routeSearchArticles) // Search articles

	// Start RSS fetcher in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// go rss.startRSSFetcher(ctx)
	go rss.StartRSSFetcher(ctx)

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
		rss.GetRss(w, r)
	case http.MethodPost:
		rss.PostRss(w, r)
	case http.MethodDelete:
		rss.DeleteRSSbyID(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func routeAllArticles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rss.GetAllArticles(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func routeSingleArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rss.GetSingleArticle(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func routeArticlesByRSS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rss.GetArticlesByRSSID(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func routeUpdateArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		rss.UpdateArticleReadStatus(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}

func routeSearchArticles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rss.SearchArticles(w, r)
	default:
		http.Error(w, "Method is not allowed or supported", http.StatusMethodNotAllowed)
	}
}
