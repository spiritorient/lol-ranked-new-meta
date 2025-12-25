package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RequestRecord represents a single request record
type RequestRecord struct {
	Timestamp   time.Time `json:"timestamp"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Referer     string    `json:"referer"`
	Country     string    `json:"country,omitempty"`     // Optional: if we add geo IP lookup
	Region      string    `json:"region,omitempty"`      // Optional: if we add geo IP lookup
	StatusCode  int       `json:"status_code"`
	ResponseTime int64    `json:"response_time_ms"`      // Response time in milliseconds
}

// AnalyticsData holds all analytics data
type AnalyticsData struct {
	TotalRequests   int64                    `json:"total_requests"`
	UniqueIPs       map[string]int           `json:"unique_ips"`
	RequestsByPath  map[string]int64         `json:"requests_by_path"`
	RequestsByMethod map[string]int64        `json:"requests_by_method"`
	RequestsByDay   map[string]int64         `json:"requests_by_day"`
	UserAgents      map[string]int           `json:"user_agents"`
	RecentRequests  []RequestRecord          `json:"recent_requests"`  // Last N requests in memory
	AllRequests     []RequestRecord          `json:"all_requests,omitempty"`  // All requests stored on disk
	FirstRequest    time.Time                `json:"first_request"`
	LastRequest     time.Time                `json:"last_request"`
	mu              sync.RWMutex             `json:"-"`
}

// Tracker handles analytics tracking
type Tracker struct {
	data       *AnalyticsData
	storage    *Storage
	maxRecent  int           // Maximum number of recent requests to keep in memory
	maxDays    int           // Maximum number of days to keep requests (0 = unlimited)
	maxRecords int          // Maximum number of total records to keep (0 = unlimited)
}

// NewTracker creates a new analytics tracker
// maxRecent: number of recent requests to keep in memory for quick access
// maxDays: maximum days to keep requests (0 = unlimited, -1 = disabled)
// maxRecords: maximum total records to keep (0 = unlimited, -1 = disabled)
func NewTracker(storagePath string, maxRecent int, maxDays int, maxRecords int) (*Tracker, error) {
	storage, err := NewStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	data, err := storage.Load()
	if err != nil {
		// If loading fails, start with empty data
		data = &AnalyticsData{
			UniqueIPs:       make(map[string]int),
			RequestsByPath:   make(map[string]int64),
			RequestsByMethod: make(map[string]int64),
			RequestsByDay:    make(map[string]int64),
			UserAgents:       make(map[string]int),
			RecentRequests:   make([]RequestRecord, 0),
			AllRequests:      make([]RequestRecord, 0),
		}
	}

	// Initialize AllRequests if nil (for backward compatibility)
	if data.AllRequests == nil {
		data.AllRequests = make([]RequestRecord, 0)
	}

	tracker := &Tracker{
		data:       data,
		storage:    storage,
		maxRecent:  maxRecent,
		maxDays:    maxDays,
		maxRecords: maxRecords,
	}

	// Clean up old records on startup if limits are set
	if maxDays > 0 || maxRecords > 0 {
		tracker.cleanupOldRecords()
	}

	return tracker, nil
}

// Track records a new request
func (t *Tracker) Track(r *http.Request, statusCode int, responseTime time.Duration) {
	t.data.mu.Lock()
	defer t.data.mu.Unlock()

	// Get client IP (handles proxies/load balancers)
	ip := getClientIP(r)
	
	// Get user agent
	userAgent := r.UserAgent()
	if userAgent == "" {
		userAgent = "Unknown"
	}

	// Get referer
	referer := r.Referer()

	now := time.Now()
	record := RequestRecord{
		Timestamp:    now,
		IP:           ip,
		UserAgent:    userAgent,
		Method:       r.Method,
		Path:         r.URL.Path,
		Referer:      referer,
		StatusCode:   statusCode,
		ResponseTime: responseTime.Milliseconds(),
	}

	// Update statistics
	t.data.TotalRequests++
	
	// Track unique IPs
	t.data.UniqueIPs[ip]++
	
	// Track by path
	t.data.RequestsByPath[r.URL.Path]++
	
	// Track by method
	t.data.RequestsByMethod[r.Method]++
	
	// Track by day
	dayKey := now.Format("2006-01-02")
	t.data.RequestsByDay[dayKey]++
	
	// Track user agents (simplified - just browser type)
	browserType := simplifyUserAgent(userAgent)
	t.data.UserAgents[browserType]++
	
	// Update timestamps
	if t.data.FirstRequest.IsZero() {
		t.data.FirstRequest = now
	}
	t.data.LastRequest = now
	
	// Add to recent requests (keep only last N in memory for quick access)
	t.data.RecentRequests = append(t.data.RecentRequests, record)
	if len(t.data.RecentRequests) > t.maxRecent {
		t.data.RecentRequests = t.data.RecentRequests[1:]
	}

	// Add to all requests (stored on disk)
	t.data.AllRequests = append(t.data.AllRequests, record)

	// Cleanup old records if limits are set (do this periodically, not on every request)
	// We'll do cleanup every 100 requests to avoid performance impact
	if len(t.data.AllRequests)%100 == 0 {
		t.cleanupOldRecords()
	}

	// Save to disk asynchronously (don't block the request)
	// Save every 10 requests to balance between data safety and performance
	go func() {
		if t.data.TotalRequests%10 == 0 {
			if err := t.storage.Save(t.data); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to save analytics: %v\n", err)
			}
		}
	}()
}

// GetStats returns current analytics statistics
func (t *Tracker) GetStats() *AnalyticsData {
	t.data.mu.RLock()
	defer t.data.mu.RUnlock()

	// Create a copy to avoid race conditions
	stats := &AnalyticsData{
		TotalRequests:    t.data.TotalRequests,
		UniqueIPs:        make(map[string]int),
		RequestsByPath:   make(map[string]int64),
		RequestsByMethod: make(map[string]int64),
		RequestsByDay:    make(map[string]int64),
		UserAgents:       make(map[string]int),
		RecentRequests:   make([]RequestRecord, len(t.data.RecentRequests)),
		FirstRequest:     t.data.FirstRequest,
		LastRequest:      t.data.LastRequest,
	}

	// Copy maps
	for k, v := range t.data.UniqueIPs {
		stats.UniqueIPs[k] = v
	}
	for k, v := range t.data.RequestsByPath {
		stats.RequestsByPath[k] = v
	}
	for k, v := range t.data.RequestsByMethod {
		stats.RequestsByMethod[k] = v
	}
	for k, v := range t.data.RequestsByDay {
		stats.RequestsByDay[k] = v
	}
	for k, v := range t.data.UserAgents {
		stats.UserAgents[k] = v
	}
	copy(stats.RecentRequests, t.data.RecentRequests)
	
	// Copy all requests
	if t.data.AllRequests != nil {
		stats.AllRequests = make([]RequestRecord, len(t.data.AllRequests))
		copy(stats.AllRequests, t.data.AllRequests)
	}

	return stats
}

// Middleware creates an HTTP middleware that tracks requests
func (t *Tracker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Call the next handler
		next.ServeHTTP(rw, r)
		
		// Track the request
		responseTime := time.Since(start)
		t.Track(r, rw.statusCode, responseTime)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// cleanupOldRecords removes old records based on maxDays and maxRecords limits
func (t *Tracker) cleanupOldRecords() {
	if t.maxDays <= 0 && t.maxRecords <= 0 {
		return // No cleanup needed
	}

	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -t.maxDays)
	
	var filtered []RequestRecord
	
	for _, record := range t.data.AllRequests {
		// Apply day-based filter
		if t.maxDays > 0 && record.Timestamp.Before(cutoffTime) {
			continue // Skip records older than maxDays
		}
		filtered = append(filtered, record)
	}
	
	// Apply record count limit (keep most recent N)
	if t.maxRecords > 0 && len(filtered) > t.maxRecords {
		// Keep only the most recent records
		start := len(filtered) - t.maxRecords
		filtered = filtered[start:]
	}
	
	t.data.AllRequests = filtered
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := len(ip) - 1; idx >= 0 && ip[idx] == ']' {
		// IPv6 address
		return ip
	}
	if idx := len(ip); idx > 0 {
		for i := idx - 1; i >= 0; i-- {
			if ip[i] == ':' {
				return ip[:i]
			}
		}
	}
	return ip
}

// simplifyUserAgent extracts a simplified browser/device type from user agent
func simplifyUserAgent(ua string) string {
	if ua == "" || ua == "Unknown" {
		return "Unknown"
	}

	// Check for mobile
	if contains(ua, "Mobile", "Android", "iPhone", "iPad") {
		if contains(ua, "Chrome") {
			return "Mobile Chrome"
		}
		if contains(ua, "Safari") {
			return "Mobile Safari"
		}
		if contains(ua, "Firefox") {
			return "Mobile Firefox"
		}
		return "Mobile Other"
	}

	// Desktop browsers
	if contains(ua, "Chrome") && !contains(ua, "Edg") {
		return "Desktop Chrome"
	}
	if contains(ua, "Firefox") {
		return "Desktop Firefox"
	}
	if contains(ua, "Safari") && !contains(ua, "Chrome") {
		return "Desktop Safari"
	}
	if contains(ua, "Edg") {
		return "Desktop Edge"
	}
	if contains(ua, "Opera") {
		return "Desktop Opera"
	}

	// Bots/crawlers
	if contains(ua, "bot", "crawler", "spider", "scraper") {
		return "Bot/Crawler"
	}

	return "Other"
}

func contains(s string, substrs ...string) bool {
	sLower := strings.ToLower(s)
	for _, substr := range substrs {
		substrLower := strings.ToLower(substr)
		if strings.Contains(sLower, substrLower) {
			return true
		}
	}
	return false
}

// GetJSON returns analytics data as JSON
func (t *Tracker) GetJSON() ([]byte, error) {
	stats := t.GetStats()
	return json.MarshalIndent(stats, "", "  ")
}

