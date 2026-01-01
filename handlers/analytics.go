package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"lol-ranked-new-meta/analytics"
)

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

type AnalyticsHandler struct {
	tracker *analytics.Tracker
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(tracker *analytics.Tracker) *AnalyticsHandler {
	return &AnalyticsHandler{
		tracker: tracker,
	}
}

// HandleAnalytics serves analytics data
// Optionally protected by a simple API key check
func (h *AnalyticsHandler) HandleAnalytics(w http.ResponseWriter, r *http.Request) {
	// Optional: Add simple authentication
	// You can set ANALYTICS_KEY environment variable to protect this endpoint
	apiKey := os.Getenv("ANALYTICS_KEY")
	if apiKey != "" {
		providedKey := r.URL.Query().Get("key")
		if providedKey != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// If browser requests HTML (and not explicitly asking for JSON), serve the dashboard
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := r.URL.Query().Get("format") == "json"
	wantsHTML := !wantsJSON && (acceptHeader == "" || 
		(len(acceptHeader) > 0 && acceptHeader[0:1] != "{" && 
		 (contains(acceptHeader, "text/html") || 
		  (contains(acceptHeader, "*/*") && !contains(acceptHeader, "application/json")))))
	
	// Serve HTML dashboard for browser requests
	if wantsHTML && r.Method == "GET" {
		http.ServeFile(w, r, "./frontend/analytics.html")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stats := h.tracker.GetStats()
	
	// Format response nicely
	response := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_requests":    stats.TotalRequests,
			"unique_ips":        len(stats.UniqueIPs),
			"total_stored":      len(stats.AllRequests),
			"first_request":     stats.FirstRequest,
			"last_request":      stats.LastRequest,
		},
		"by_path":      stats.RequestsByPath,
		"by_method":    stats.RequestsByMethod,
		"by_day":       stats.RequestsByDay,
		"user_agents":  stats.UserAgents,
		"top_ips":      getTopN(stats.UniqueIPs, 10),
		"recent_requests": stats.RecentRequests,
		"all_requests_count": len(stats.AllRequests),
	}
	
	// Include all requests if requested (via ?all=true query param)
	// Be careful: this could be a large response!
	if r.URL.Query().Get("all") == "true" {
		response["all_requests"] = stats.AllRequests
	} else {
		// Just show a sample of recent requests from all stored
		if len(stats.AllRequests) > 0 {
			// Show last 50 from all stored requests
			start := len(stats.AllRequests) - 50
			if start < 0 {
				start = 0
			}
			response["all_requests_sample"] = stats.AllRequests[start:]
			response["all_requests_note"] = "Showing last 50 requests. Use ?all=true to get all requests."
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// getTopN returns the top N items from a map sorted by value
func getTopN(m map[string]int, n int) map[string]int {
	if len(m) <= n {
		return m
	}

	// Simple implementation - get top N by value
	type kv struct {
		Key   string
		Value int
	}

	var items []kv
	for k, v := range m {
		items = append(items, kv{k, v})
	}

	// Sort by value (descending)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Value > items[i].Value {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// Take top N
	result := make(map[string]int)
	for i := 0; i < n && i < len(items); i++ {
		result[items[i].Key] = items[i].Value
	}

	return result
}

