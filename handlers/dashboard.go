package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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

// HandleSaveMatch saves a match to a dashboard
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

	// Fetch match from Riot API
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

	// Convert to dashboard format
	dashMatch := dashboard.ConvertRiotMatch(riotMatch, region)

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
		Message:     "Match saved successfully",
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

