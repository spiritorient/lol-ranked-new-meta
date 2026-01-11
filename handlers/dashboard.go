package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"lol-ranked-new-meta/dashboard"
	"lol-ranked-new-meta/riot"
)

// DashboardHandler handles dashboard-related requests
type DashboardHandler struct {
	storage    *dashboard.Storage
	riotClient *riot.Client
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(storage *dashboard.Storage, riotClient *riot.Client) *DashboardHandler {
	return &DashboardHandler{
		storage:    storage,
		riotClient: riotClient,
	}
}

// SaveMatchRequest represents a request to save a match to dashboard
type SaveMatchRequest struct {
	MatchID     string `json:"match_id"`
	DashboardID string `json:"dashboard_id"` // Optional - will create new if empty
}

// SaveMatchResponse represents the response after saving a match
type SaveMatchResponse struct {
	Success     bool   `json:"success"`
	DashboardID string `json:"dashboard_id"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

// HandleSaveMatch saves a match to a dashboard (stores ALL Riot API data)
func (h *DashboardHandler) HandleSaveMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(SaveMatchResponse{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return
	}

	var req SaveMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(SaveMatchResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.MatchID == "" {
		json.NewEncoder(w).Encode(SaveMatchResponse{
			Success: false,
			Error:   "match_id is required",
		})
		return
	}

	// Generate dashboard ID if not provided
	dashboardID := req.DashboardID
	if dashboardID == "" {
		dashboardID = generateDashboardID()
	}

	// Sanitize dashboard ID
	dashboardID = sanitizeDashboardID(dashboardID)

	// Fetch COMPLETE match data from Riot API
	riotMatch, err := h.riotClient.GetMatch(req.MatchID)
	if err != nil {
		json.NewEncoder(w).Encode(SaveMatchResponse{
			Success: false,
			Error:   "Failed to fetch match: " + err.Error(),
		})
		return
	}

	// Extract region from match ID (format: REGION_GAMEID)
	region := "unknown"
	if parts := strings.Split(req.MatchID, "_"); len(parts) > 0 {
		region = parts[0]
	}

	// Create dashboard match with FULL Riot data
	dashMatch := dashboard.CreateDashboardMatch(riotMatch, dashboardID, region)

	// Save to dashboard
	if err := h.storage.AddMatch(dashboardID, dashMatch); err != nil {
		json.NewEncoder(w).Encode(SaveMatchResponse{
			Success: false,
			Error:   "Failed to save match: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(SaveMatchResponse{
		Success:     true,
		DashboardID: dashboardID,
		Message:     "Match saved with ALL data from Riot API",
	})
}

// HandleGetDashboard retrieves a dashboard by ID
func (h *DashboardHandler) HandleGetDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if browser wants HTML
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := r.URL.Query().Get("format") == "json"
	wantsHTML := !wantsJSON && strings.Contains(acceptHeader, "text/html")

	// Get dashboard ID from URL path
	path := r.URL.Path
	dashboardID := strings.TrimPrefix(path, "/d/")
	dashboardID = strings.TrimPrefix(dashboardID, "/dashboard/")

	// If just /dashboard with no ID, serve the page which will prompt for ID
	if dashboardID == "" || dashboardID == "dashboard" || dashboardID == "d" {
		if wantsHTML {
			http.ServeFile(w, r, "./frontend/dashboard.html")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "dashboard_id is required. Use /d/YOUR_ID or /dashboard/YOUR_ID",
		})
		return
	}

	// Serve HTML dashboard page for browser requests
	if wantsHTML {
		http.ServeFile(w, r, "./frontend/dashboard.html")
		return
	}

	// Return JSON data
	w.Header().Set("Content-Type", "application/json")

	dashboardID = sanitizeDashboardID(dashboardID)

	dashboardData, err := h.storage.LoadDashboard(dashboardID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to load dashboard: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(dashboardData)
}

// HandleListDashboards returns all dashboard IDs with summary info
func (h *DashboardHandler) HandleListDashboards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if browser wants HTML
	acceptHeader := r.Header.Get("Accept")
	wantsHTML := strings.Contains(acceptHeader, "text/html")

	ids, err := h.storage.ListDashboards()
	if err != nil {
		if wantsHTML {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<h1>Error loading dashboards</h1><p>" + err.Error() + "</p>"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to list dashboards: " + err.Error(),
		})
		return
	}

	// Build dashboard list with match counts
	type DashboardSummary struct {
		ID          string `json:"id"`
		MatchCount  int    `json:"match_count"`
		LastUpdated string `json:"last_updated"`
	}

	summaries := make([]DashboardSummary, 0, len(ids))
	for _, id := range ids {
		data, err := h.storage.LoadDashboard(id)
		if err != nil {
			continue
		}
		summaries = append(summaries, DashboardSummary{
			ID:          id,
			MatchCount:  len(data.Matches),
			LastUpdated: data.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if wantsHTML {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head>
	<title>All Dashboards - New Meta</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		:root { --bg: #0a0a0f; --card: #12121a; --primary: #00d4aa; --text: #e0e0e0; }
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { font-family: 'Segoe UI', sans-serif; background: var(--bg); color: var(--text); padding: 2rem; }
		h1 { color: var(--primary); margin-bottom: 1.5rem; }
		.dashboard-list { display: grid; gap: 1rem; max-width: 800px; }
		.dashboard-card { background: var(--card); padding: 1rem 1.5rem; border-radius: 8px; display: flex; justify-content: space-between; align-items: center; border-left: 3px solid var(--primary); }
		.dashboard-card:hover { background: #1a1a24; }
		.dashboard-id { font-weight: 600; font-size: 1.1rem; }
		.dashboard-id a { color: var(--primary); text-decoration: none; }
		.dashboard-id a:hover { text-decoration: underline; }
		.dashboard-info { color: #888; font-size: 0.85rem; }
		.match-count { background: var(--primary); color: #000; padding: 0.25rem 0.75rem; border-radius: 20px; font-weight: 600; }
		.empty { color: #666; font-style: italic; }
	</style>
</head>
<body>
	<h1>📊 All Dashboards</h1>
	<div class="dashboard-list">`

		if len(summaries) == 0 {
			html += `<p class="empty">No dashboards created yet.</p>`
		} else {
			for _, s := range summaries {
				html += `<div class="dashboard-card">
					<div>
						<div class="dashboard-id"><a href="/d/` + s.ID + `">` + s.ID + `</a></div>
						<div class="dashboard-info">Last updated: ` + s.LastUpdated + `</div>
					</div>
					<span class="match-count">` + fmt.Sprintf("%d", s.MatchCount) + ` matches</span>
				</div>`
			}
		}

		html += `</div>
	<p style="margin-top: 2rem; color: #666;"><a href="/" style="color: var(--primary);">← Back to Analyzer</a></p>
</body>
</html>`
		w.Write([]byte(html))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dashboards": summaries,
		"total":      len(summaries),
	})
}

// generateDashboardID creates a random dashboard ID
func generateDashboardID() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// sanitizeDashboardID removes unsafe characters from dashboard ID
func sanitizeDashboardID(id string) string {
	// Only allow alphanumeric, dash, underscore
	var result strings.Builder
	for _, c := range id {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result.WriteRune(c)
		}
	}
	sanitized := result.String()
	if sanitized == "" {
		sanitized = generateDashboardID()
	}
	return sanitized
}
