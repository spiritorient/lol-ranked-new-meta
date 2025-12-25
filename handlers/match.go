package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"lol-ranked-new-meta/openai"
	"lol-ranked-new-meta/riot"
	"lol-ranked-new-meta/types"
)

type MatchHandler struct {
	riotClient   *riot.Client
	openaiClient *openai.Client
}

// NewMatchHandler creates a new match handler
func NewMatchHandler(riotClient *riot.Client, openaiClient *openai.Client) *MatchHandler {
	return &MatchHandler{
		riotClient:   riotClient,
		openaiClient: openaiClient,
	}
}

// HandleAnalyzeMatch handles requests to analyze a match
func (h *MatchHandler) HandleAnalyzeMatch(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for frontend access
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MatchID == "" {
		h.sendError(w, "match_id is required", http.StatusBadRequest)
		return
	}

	// Fetch match data from Riot API
	log.Printf("Fetching match data for match ID: %s", req.MatchID)
	match, err := h.riotClient.GetMatch(req.MatchID)
	if err != nil {
		log.Printf("Error fetching match: %v", err)
		h.sendError(w, "Failed to fetch match data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Format match data for analysis (with optional champion/summoner filter)
	matchSummary := riot.FormatMatchForAnalysis(match, req.ChampionName, req.SummonerName)

	// Analyze match using OpenAI
	log.Printf("Analyzing match using OpenAI")
	if req.ChampionName != "" {
		log.Printf("Deep dive requested for champion: %s", req.ChampionName)
	} else if req.SummonerName != "" {
		log.Printf("Deep dive requested for summoner: %s", req.SummonerName)
	}
	if len(req.FocusAreas) > 0 {
		log.Printf("Focus areas requested: %v", req.FocusAreas)
	}
	analysis, err := h.openaiClient.AnalyzeMatch(r.Context(), matchSummary, req.ChampionName, req.SummonerName, req.FocusAreas)
	if err != nil {
		log.Printf("Error analyzing match: %v", err)
		h.sendError(w, "Failed to analyze match: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set match ID in response
	analysis.MatchID = req.MatchID

	// Send response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// HandleAnalyzeMatchGET handles GET requests (for testing convenience)
func (h *MatchHandler) HandleAnalyzeMatchGET(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for frontend access
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	matchID := r.URL.Query().Get("match_id")
	if matchID == "" {
		h.sendError(w, "match_id query parameter is required", http.StatusBadRequest)
		return
	}

	// Get optional champion/summoner filter from query params
	championFilter := r.URL.Query().Get("champion_name")
	summonerFilter := r.URL.Query().Get("summoner_name")
	
	// Get optional focus areas from query params (comma-separated)
	var focusAreas []string
	if focusAreasStr := r.URL.Query().Get("focus_areas"); focusAreasStr != "" {
		focusAreas = strings.Split(focusAreasStr, ",")
		for i := range focusAreas {
			focusAreas[i] = strings.TrimSpace(focusAreas[i])
		}
	}

	// Fetch match data from Riot API
	log.Printf("Fetching match data for match ID: %s", matchID)
	match, err := h.riotClient.GetMatch(matchID)
	if err != nil {
		log.Printf("Error fetching match: %v", err)
		h.sendError(w, "Failed to fetch match data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Format match data for analysis (with optional champion/summoner filter)
	matchSummary := riot.FormatMatchForAnalysis(match, championFilter, summonerFilter)

	// Analyze match using OpenAI
	log.Printf("Analyzing match using OpenAI")
	if championFilter != "" {
		log.Printf("Deep dive requested for champion: %s", championFilter)
	} else if summonerFilter != "" {
		log.Printf("Deep dive requested for summoner: %s", summonerFilter)
	}
	if len(focusAreas) > 0 {
		log.Printf("Focus areas requested: %v", focusAreas)
	}
	analysis, err := h.openaiClient.AnalyzeMatch(r.Context(), matchSummary, championFilter, summonerFilter, focusAreas)
	if err != nil {
		log.Printf("Error analyzing match: %v", err)
		h.sendError(w, "Failed to analyze match: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set match ID in response
	analysis.MatchID = matchID

	// Send response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *MatchHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := types.MatchResponse{
		Error: message,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
