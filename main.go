package main

import (
	"log"
	"net/http"

	"lol-ranked-new-meta/config"
	"lol-ranked-new-meta/handlers"
	"lol-ranked-new-meta/openai"
	"lol-ranked-new-meta/riot"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize clients
	riotClient := riot.NewClient(cfg.RiotAPIKey, cfg.RiotAPIRegion)
	openaiClient := openai.NewClient(cfg.OpenAIAPIKey, cfg.OpenAIModel)

	// Create handlers
	matchHandler := handlers.NewMatchHandler(riotClient, openaiClient)

	// Set up API routes (must be before frontend to take precedence)
	http.HandleFunc("/analyze-match", matchHandler.HandleAnalyzeMatch)
	http.HandleFunc("/analyze-match-get", matchHandler.HandleAnalyzeMatchGET) // Convenience GET endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve frontend static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/style.css", fs)
	http.Handle("/script.js", fs)
	http.Handle("/logo.png", fs)

	// Serve index.html for root and non-API routes (SPA routing)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Serve API routes normally (they're already registered above)
		if path == "/analyze-match" || path == "/analyze-match-get" || path == "/health" {
			// This won't be reached since those routes are registered first, but good to check
			return
		}
		// Serve index.html for all other routes
		http.ServeFile(w, r, "./frontend/index.html")
	})

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Endpoints available:")
	log.Printf("  POST /analyze-match - Analyze a match (requires JSON body with match_id)")
	log.Printf("  GET  /analyze-match-get?match_id=<match_id> - Analyze a match (convenience endpoint)")
	log.Printf("  GET  /health - Health check")

	// Bind to all interfaces (0.0.0.0) for cloud deployment compatibility
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
