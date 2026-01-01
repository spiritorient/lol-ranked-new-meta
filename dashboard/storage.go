package dashboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"lol-ranked-new-meta/types"
)

// DashboardMatch stores the COMPLETE match data from Riot API
type DashboardMatch struct {
	// Metadata
	SavedAt     time.Time `json:"saved_at"`
	DashboardID string    `json:"dashboard_id"`
	Region      string    `json:"region"`

	// Full Riot Match Data - stores EVERYTHING from the API
	RiotMatch types.RiotMatch `json:"riot_match"`
}

// DashboardData contains all saved matches for a dashboard
type DashboardData struct {
	DashboardID string           `json:"dashboard_id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Matches     []DashboardMatch `json:"matches"`
}

// Storage handles persistent storage of dashboard data
type Storage struct {
	basePath string
	mu       sync.RWMutex
}

// NewStorage creates a new dashboard storage instance
func NewStorage(basePath string) (*Storage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create dashboard storage directory: %w", err)
	}

	return &Storage{
		basePath: basePath,
	}, nil
}

// getDashboardPath returns the file path for a dashboard
func (s *Storage) getDashboardPath(dashboardID string) string {
	return filepath.Join(s.basePath, dashboardID+".json")
}

// LoadDashboard loads a dashboard by ID
func (s *Storage) LoadDashboard(dashboardID string) (*DashboardData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := s.getDashboardPath(dashboardID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new dashboard
			return &DashboardData{
				DashboardID: dashboardID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Matches:     make([]DashboardMatch, 0),
			}, nil
		}
		return nil, fmt.Errorf("failed to read dashboard file: %w", err)
	}

	var dashboard DashboardData
	if err := json.Unmarshal(data, &dashboard); err != nil {
		return nil, fmt.Errorf("failed to parse dashboard data: %w", err)
	}

	return &dashboard, nil
}

// SaveDashboard saves a dashboard to disk
func (s *Storage) SaveDashboard(dashboard *DashboardData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dashboard.UpdatedAt = time.Now()

	filePath := s.getDashboardPath(dashboard.DashboardID)
	tmpFile := filePath + ".tmp"

	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dashboard); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to encode dashboard data: %w", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpFile, filePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// AddMatch adds a match to a dashboard
func (s *Storage) AddMatch(dashboardID string, match DashboardMatch) error {
	dashboard, err := s.LoadDashboard(dashboardID)
	if err != nil {
		return err
	}

	// Check if match already exists
	for i, m := range dashboard.Matches {
		if m.RiotMatch.Metadata.MatchID == match.RiotMatch.Metadata.MatchID {
			// Update existing match
			dashboard.Matches[i] = match
			return s.SaveDashboard(dashboard)
		}
	}

	// Add new match
	dashboard.Matches = append(dashboard.Matches, match)
	return s.SaveDashboard(dashboard)
}

// CreateDashboardMatch creates a DashboardMatch from a RiotMatch
func CreateDashboardMatch(riotMatch *types.RiotMatch, dashboardID, region string) DashboardMatch {
	return DashboardMatch{
		SavedAt:     time.Now(),
		DashboardID: dashboardID,
		Region:      region,
		RiotMatch:   *riotMatch,
	}
}

// ListDashboards returns all dashboard IDs
func (s *Storage) ListDashboards() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			id := entry.Name()[:len(entry.Name())-5] // Remove .json
			ids = append(ids, id)
		}
	}

	return ids, nil
}
