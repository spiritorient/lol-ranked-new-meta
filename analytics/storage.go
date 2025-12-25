package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Storage handles persistent storage of analytics data
type Storage struct {
	filePath string
}

// NewStorage creates a new storage instance
func NewStorage(filePath string) (*Storage, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create storage directory: %w", err)
		}
	}

	return &Storage{
		filePath: filePath,
	}, nil
}

// Load loads analytics data from disk
func (s *Storage) Load() (*AnalyticsData, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return empty data
			return &AnalyticsData{
				UniqueIPs:       make(map[string]int),
				RequestsByPath:  make(map[string]int64),
				RequestsByMethod: make(map[string]int64),
				RequestsByDay:   make(map[string]int64),
				UserAgents:      make(map[string]int),
				RecentRequests:  make([]RequestRecord, 0),
				AllRequests:     make([]RequestRecord, 0),
			}, nil
		}
		return nil, fmt.Errorf("failed to read analytics file: %w", err)
	}

	var analyticsData AnalyticsData
	if err := json.Unmarshal(data, &analyticsData); err != nil {
		return nil, fmt.Errorf("failed to parse analytics data: %w", err)
	}

	// Initialize maps if they're nil (for backward compatibility)
	if analyticsData.UniqueIPs == nil {
		analyticsData.UniqueIPs = make(map[string]int)
	}
	if analyticsData.RequestsByPath == nil {
		analyticsData.RequestsByPath = make(map[string]int64)
	}
	if analyticsData.RequestsByMethod == nil {
		analyticsData.RequestsByMethod = make(map[string]int64)
	}
	if analyticsData.RequestsByDay == nil {
		analyticsData.RequestsByDay = make(map[string]int64)
	}
	if analyticsData.UserAgents == nil {
		analyticsData.UserAgents = make(map[string]int)
	}
	if analyticsData.RecentRequests == nil {
		analyticsData.RecentRequests = make([]RequestRecord, 0)
	}
	if analyticsData.AllRequests == nil {
		analyticsData.AllRequests = make([]RequestRecord, 0)
	}

	return &analyticsData, nil
}

// Save saves analytics data to disk
func (s *Storage) Save(data *AnalyticsData) error {
	// Create a temporary file first, then rename (atomic write)
	tmpFile := s.filePath + ".tmp"
	
	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to encode analytics data: %w", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, s.filePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

