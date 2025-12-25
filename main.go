package main

import (
	"log"
	"net/http"

	"lol-ranked-new-meta/analytics"
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

	// Initialize analytics tracker
	// Keep last 100 requests in memory for quick access
	// Store ALL requests on disk (unlimited by default, configurable via env vars)
	analyticsTracker, err := analytics.NewTracker(
		cfg.AnalyticsDataPath,
		100,                      // maxRecent: keep 100 in memory
		cfg.AnalyticsMaxDays,     // maxDays: 0 = unlimited
		cfg.AnalyticsMaxRecords,  // maxRecords: 0 = unlimited
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize analytics tracker: %v", err)
		log.Printf("Analytics will not be tracked")
		analyticsTracker = nil
	} else {
		log.Printf("Analytics tracking enabled (data stored at: %s)", cfg.AnalyticsDataPath)
	}

	// Create handlers
	matchHandler := handlers.NewMatchHandler(riotClient, openaiClient)
	
	// Create analytics handler (if tracker is available)
	var analyticsHandler *handlers.AnalyticsHandler
	if analyticsTracker != nil {
		analyticsHandler = handlers.NewAnalyticsHandler(analyticsTracker)
	}

	// Create a new mux
	mux := http.NewServeMux()

	// Set up API routes (must be before frontend to take precedence)
	mux.HandleFunc("/analyze-match", matchHandler.HandleAnalyzeMatch)
	mux.HandleFunc("/analyze-match-get", matchHandler.HandleAnalyzeMatchGET) // Convenience GET endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Analytics endpoint (only if tracker is available)
	if analyticsHandler != nil {
		mux.HandleFunc("/analytics", analyticsHandler.HandleAnalytics)
		log.Printf("  GET  /analytics?key=<ANALYTICS_KEY> - View analytics (optional key protection)")
	}

	// Serve frontend static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/style.css", fs)
	mux.Handle("/script.js", fs)
	mux.Handle("/logo.png", fs)

	// Serve index.html for root and non-API routes (SPA routing)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Serve API routes normally (they're already registered above)
		if path == "/analyze-match" || path == "/analyze-match-get" || path == "/health" || path == "/analytics" {
			// This won't be reached since those routes are registered first, but good to check
			return
		}
		// Serve index.html for all other routes
		http.ServeFile(w, r, "./frontend/index.html")
	})

	// Wrap mux with analytics middleware if tracker is available
	var handler http.Handler = mux
	if analyticsTracker != nil {
		handler = analyticsTracker.Middleware(mux)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Endpoints available:")
	log.Printf("  POST /analyze-match - Analyze a match (requires JSON body with match_id)")
	log.Printf("  GET  /analyze-match-get?match_id=<match_id> - Analyze a match (convenience endpoint)")
	log.Printf("  GET  /health - Health check")

	// Bind to all interfaces (0.0.0.0) for cloud deployment compatibility
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
