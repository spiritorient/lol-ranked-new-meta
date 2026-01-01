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

// DashboardMatch represents a saved match for the dashboard
type DashboardMatch struct {
	// Match identification
	MatchID   string    `json:"match_id"`
	Region    string    `json:"region"`
	SavedAt   time.Time `json:"saved_at"`
	
	// Game info
	GameMode     string `json:"game_mode"`
	GameDuration int64  `json:"game_duration"` // seconds
	GameVersion  string `json:"game_version"`
	
	// Team results
	BlueTeamWin bool `json:"blue_team_win"`
	
	// All participants with full stats
	Participants []ParticipantStats `json:"participants"`
	
	// Team objectives
	BlueTeamObjectives TeamObjectives `json:"blue_team_objectives"`
	RedTeamObjectives  TeamObjectives `json:"red_team_objectives"`
}

// ParticipantStats contains all stats for a participant
type ParticipantStats struct {
	// Identity
	SummonerName  string `json:"summoner_name"`
	RiotIDName    string `json:"riot_id_name"`
	RiotIDTagline string `json:"riot_id_tagline"`
	ChampionName  string `json:"champion_name"`
	ChampionID    int    `json:"champion_id"`
	TeamID        int    `json:"team_id"` // 100 = Blue, 200 = Red
	TeamPosition  string `json:"team_position"`
	Role          string `json:"role"`
	
	// Result
	Win bool `json:"win"`
	
	// Combat
	Kills                       int `json:"kills"`
	Deaths                      int `json:"deaths"`
	Assists                     int `json:"assists"`
	TotalDamageDealtToChampions int `json:"total_damage_dealt_to_champions"`
	PhysicalDamageDealt         int `json:"physical_damage_dealt"`
	MagicDamageDealt            int `json:"magic_damage_dealt"`
	TrueDamageDealt             int `json:"true_damage_dealt"`
	TotalDamageTaken            int `json:"total_damage_taken"`
	DamageSelfMitigated         int `json:"damage_self_mitigated"`
	TotalHeal                   int `json:"total_heal"`
	TotalShielded               int `json:"total_shielded"`
	
	// Multi-kills
	DoubleKills        int `json:"double_kills"`
	TripleKills        int `json:"triple_kills"`
	QuadraKills        int `json:"quadra_kills"`
	PentaKills         int `json:"penta_kills"`
	LargestKillingSpree int `json:"largest_killing_spree"`
	LargestMultiKill   int `json:"largest_multi_kill"`
	
	// Objectives
	TurretKills    int  `json:"turret_kills"`
	InhibitorKills int  `json:"inhibitor_kills"`
	DragonKills    int  `json:"dragon_kills"`
	BaronKills     int  `json:"baron_kills"`
	FirstBlood     bool `json:"first_blood"`
	FirstTower     bool `json:"first_tower"`
	
	// Economy
	GoldEarned int `json:"gold_earned"`
	GoldSpent  int `json:"gold_spent"`
	
	// Farming
	TotalMinionsKilled   int `json:"total_minions_killed"`
	NeutralMinionsKilled int `json:"neutral_minions_killed"`
	
	// Vision
	VisionScore         int `json:"vision_score"`
	WardsPlaced         int `json:"wards_placed"`
	WardsKilled         int `json:"wards_killed"`
	ControlWardsBought  int `json:"control_wards_bought"`
	
	// Items
	Item0 int `json:"item0"`
	Item1 int `json:"item1"`
	Item2 int `json:"item2"`
	Item3 int `json:"item3"`
	Item4 int `json:"item4"`
	Item5 int `json:"item5"`
	Item6 int `json:"item6"` // Trinket
	
	// Time
	TotalTimeSpentDead     int `json:"total_time_spent_dead"`
	LongestTimeSpentLiving int `json:"longest_time_spent_living"`
	ChampLevel             int `json:"champ_level"`
}

// TeamObjectives contains team objective stats
type TeamObjectives struct {
	TowerKills      int  `json:"tower_kills"`
	InhibitorKills  int  `json:"inhibitor_kills"`
	DragonKills     int  `json:"dragon_kills"`
	BaronKills      int  `json:"baron_kills"`
	RiftHeraldKills int  `json:"rift_herald_kills"`
	FirstTower      bool `json:"first_tower"`
	FirstBlood      bool `json:"first_blood"`
	FirstDragon     bool `json:"first_dragon"`
	FirstBaron      bool `json:"first_baron"`
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
		if m.MatchID == match.MatchID {
			// Update existing match
			dashboard.Matches[i] = match
			return s.SaveDashboard(dashboard)
		}
	}
	
	// Add new match
	dashboard.Matches = append(dashboard.Matches, match)
	return s.SaveDashboard(dashboard)
}

// ConvertRiotMatch converts a Riot API match to dashboard format
func ConvertRiotMatch(riotMatch *types.RiotMatch, region string) DashboardMatch {
	match := DashboardMatch{
		MatchID:      riotMatch.Metadata.MatchID,
		Region:       region,
		SavedAt:      time.Now(),
		GameMode:     riotMatch.Info.GameMode,
		GameDuration: riotMatch.Info.GameDuration,
		GameVersion:  riotMatch.Info.GameVersion,
		Participants: make([]ParticipantStats, 0, len(riotMatch.Info.Participants)),
	}
	
	// Process teams
	for _, team := range riotMatch.Info.Teams {
		objectives := TeamObjectives{
			TowerKills:      team.Objectives.Tower.Kills,
			InhibitorKills:  team.Objectives.Inhibitor.Kills,
			DragonKills:     team.Objectives.Dragon.Kills,
			BaronKills:      team.Objectives.Baron.Kills,
			RiftHeraldKills: team.Objectives.RiftHerald.Kills,
			FirstTower:      team.Objectives.Tower.First,
			FirstBlood:      team.Objectives.Champion.First,
			FirstDragon:     team.Objectives.Dragon.First,
			FirstBaron:      team.Objectives.Baron.First,
		}
		
		if team.TeamID == 100 {
			match.BlueTeamWin = team.Win
			match.BlueTeamObjectives = objectives
		} else {
			match.RedTeamObjectives = objectives
		}
	}
	
	// Process participants
	for _, p := range riotMatch.Info.Participants {
		participant := ParticipantStats{
			SummonerName:                p.SummonerName,
			RiotIDName:                  p.RiotIDName,
			RiotIDTagline:               p.RiotIDTagline,
			ChampionName:                p.ChampionName,
			ChampionID:                  p.ChampionID,
			TeamID:                      p.TeamID,
			TeamPosition:                p.TeamPosition,
			Role:                        p.Role,
			Win:                         p.Win,
			Kills:                       p.Kills,
			Deaths:                      p.Deaths,
			Assists:                     p.Assists,
			TotalDamageDealtToChampions: p.TotalDamageDealtToChampions,
			PhysicalDamageDealt:         p.PhysicalDamageDealtToChampions,
			MagicDamageDealt:            p.MagicDamageDealtToChampions,
			TrueDamageDealt:             p.TrueDamageDealtToChampions,
			TotalDamageTaken:            p.TotalDamageTaken,
			DamageSelfMitigated:         p.DamageSelfMitigated,
			TotalHeal:                   p.TotalHeal,
			TotalShielded:               p.TotalDamageShieldedOnTeammates,
			DoubleKills:                 p.DoubleKills,
			TripleKills:                 p.TripleKills,
			QuadraKills:                 p.QuadraKills,
			PentaKills:                  p.PentaKills,
			LargestKillingSpree:         p.LargestKillingSpree,
			LargestMultiKill:            p.LargestMultiKill,
			TurretKills:                 p.TurretKills,
			InhibitorKills:              p.InhibitorKills,
			DragonKills:                 p.DragonKills,
			BaronKills:                  p.BaronKills,
			FirstBlood:                  p.FirstBloodKill,
			FirstTower:                  p.FirstTowerKill,
			GoldEarned:                  p.GoldEarned,
			GoldSpent:                   p.GoldSpent,
			TotalMinionsKilled:          p.TotalMinionsKilled,
			NeutralMinionsKilled:        p.NeutralMinionsKilled,
			VisionScore:                 p.VisionScore,
			WardsPlaced:                 p.WardsPlaced,
			WardsKilled:                 p.WardsKilled,
			ControlWardsBought:          p.VisionWardsBoughtInGame,
			Item0:                       p.Item0,
			Item1:                       p.Item1,
			Item2:                       p.Item2,
			Item3:                       p.Item3,
			Item4:                       p.Item4,
			Item5:                       p.Item5,
			Item6:                       p.Item6,
			TotalTimeSpentDead:          p.TotalTimeSpentDead,
			LongestTimeSpentLiving:      p.LongestTimeSpentLiving,
			ChampLevel:                  p.ChampLevel,
		}
		match.Participants = append(match.Participants, participant)
	}
	
	return match
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

