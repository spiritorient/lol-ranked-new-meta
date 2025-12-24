package types

// MatchRequest represents the incoming request for match analysis
type MatchRequest struct {
	MatchID       string `json:"match_id"`
	Region        string `json:"region,omitempty"`        // Optional: overrides default region
	ChampionName  string `json:"champion_name,omitempty"` // Optional: for deep dive analysis on specific champion
	SummonerName  string `json:"summoner_name,omitempty"` // Optional: for deep dive analysis on specific summoner
}

// MatchResponse represents the response from the match advisor
type MatchResponse struct {
	MatchID          string            `json:"match_id"`
	Analysis         string            `json:"analysis"`
	Suggestions      []string          `json:"suggestions"`
	CoachingTips     []string          `json:"coaching_tips"`
	ChampionDeepDive string            `json:"champion_deep_dive,omitempty"` // Optional: deep dive analysis for specific champion
	StructuredInsights *StructuredInsights `json:"structured_insights,omitempty"` // New: structured data-driven insights
	Error            string            `json:"error,omitempty"`
}

// StructuredInsights provides specific, data-driven insights about the match
type StructuredInsights struct {
	WhatWentWell     []SpecificEvent `json:"what_went_well"`
	WhatWentWrong    []SpecificEvent `json:"what_went_wrong"`
	CriticalMoments  []CriticalMoment `json:"critical_moments"`
	ItemAnalysis     *ItemAnalysis    `json:"item_analysis,omitempty"`
	MatchupAnalysis  *MatchupAnalysis `json:"matchup_analysis,omitempty"`
	KeyStatistics    KeyStatistics    `json:"key_statistics"`
}

// SpecificEvent represents a concrete event that happened in the match
type SpecificEvent struct {
	Title       string   `json:"title"`        // e.g., "Early First Blood"
	Description string   `json:"description"`  // What actually happened
	Impact      string   `json:"impact"`       // Why it mattered
	Data        []string `json:"data"`         // Supporting data/numbers
	Category    string   `json:"category"`     // "objective", "combat", "vision", "farming", etc.
}

// CriticalMoment represents a key moment in the game
type CriticalMoment struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Outcome     string   `json:"outcome"`
	Impact      string   `json:"impact"`
	Data        []string `json:"data"`
}

// ItemAnalysis provides detailed item build analysis
type ItemAnalysis struct {
	BuildPath      []ItemTiming `json:"build_path"`
	TimingAnalysis string       `json:"timing_analysis"`
	OpponentMatchup string      `json:"opponent_matchup"` // How items countered opponent composition
	Recommendations []string    `json:"recommendations"`
}

// ItemTiming represents when an item was purchased
type ItemTiming struct {
	ItemID    int    `json:"item_id"`
	ItemName  string `json:"item_name,omitempty"`
	TimeBought string `json:"time_bought"` // e.g., "12:34" or "early/mid/late"
	Context   string `json:"context"`      // Why this item at this time
}

// MatchupAnalysis provides champion matchup and team composition analysis
type MatchupAnalysis struct {
	LaneMatchup      string   `json:"lane_matchup"`       // How the champion fared vs opponent
	TeamComposition  string   `json:"team_composition"`   // Overall team comp analysis
	Synergies        []string `json:"synergies"`          // Good synergies with teammates
	Counters         []string `json:"counters"`           // Champion counters in this match
	WinConditions    []string `json:"win_conditions"`     // What needed to happen to win
}

// KeyStatistics highlights important numbers from the match
type KeyStatistics struct {
	Combat      []StatPair `json:"combat"`
	Objectives  []StatPair `json:"objectives"`
	Economy     []StatPair `json:"economy"`
	Vision      []StatPair `json:"vision"`
}

// StatPair represents a key statistic
type StatPair struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Context string `json:"context,omitempty"` // Additional context or comparison
}

// RiotMatch represents the structure of match data from Riot API
// This is a simplified version - you may need to expand based on your needs
type RiotMatch struct {
	Metadata RiotMatchMetadata `json:"metadata"`
	Info     RiotMatchInfo     `json:"info"`
}

type RiotMatchMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type RiotMatchInfo struct {
	GameCreation       int64                `json:"gameCreation"`
	GameDuration       int64                `json:"gameDuration"`
	GameEndTimestamp   int64                `json:"gameEndTimestamp"`
	GameID             int64                `json:"gameId"`
	GameMode           string               `json:"gameMode"`
	GameName           string               `json:"gameName"`
	GameStartTimestamp int64                `json:"gameStartTimestamp"`
	GameType           string               `json:"gameType"`
	GameVersion        string               `json:"gameVersion"`
	MapID              int                  `json:"mapId"`
	Participants       []RiotParticipant    `json:"participants"`
	PlatformID         string               `json:"platformId"`
	QueueID            int                  `json:"queueId"`
	Teams              []RiotTeam           `json:"teams"`
	TournamentCode     string               `json:"tournamentCode"`
}

type RiotParticipant struct {
	Assists                        int                    `json:"assists"`
	BaronKills                     int                    `json:"baronKills"`
	BountyLevel                    int                    `json:"bountyLevel"`
	Challenges                     map[string]interface{} `json:"challenges"`
	ChampLevel                     int                    `json:"champLevel"`
	ChampionID                     int                    `json:"championId"`
	ChampionName                   string                 `json:"championName"`
	ChampionTransform              int                    `json:"championTransform"`
	ConsumablesPurchased           int                    `json:"consumablesPurchased"`
	DamageDealtToBuildings         int                    `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int                    `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int                    `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int                    `json:"damageSelfMitigated"`
	Deaths                          int                    `json:"deaths"`
	DetectorWardsPlaced            int                    `json:"detectorWardsPlaced"`
	DoubleKills                     int                    `json:"doubleKills"`
	DragonKills                     int                    `json:"dragonKills"`
	FirstBloodAssist               bool                   `json:"firstBloodAssist"`
	FirstBloodKill                 bool                   `json:"firstBloodKill"`
	FirstTowerAssist               bool                   `json:"firstTowerAssist"`
	FirstTowerKill                 bool                   `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool                   `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool                   `json:"gameEndedInSurrender"`
	GoldEarned                     int                    `json:"goldEarned"`
	GoldSpent                      int                    `json:"goldSpent"`
	IndividualPosition             string                 `json:"individualPosition"`
	InhibitorKills                 int                    `json:"inhibitorKills"`
	InhibitorTakedowns             int                    `json:"inhibitorTakedowns"`
	InhibitorsLost                 int                    `json:"inhibitorsLost"`
	Item0                          int                    `json:"item0"`
	Item1                          int                    `json:"item1"`
	Item2                          int                    `json:"item2"`
	Item3                          int                    `json:"item3"`
	Item4                          int                    `json:"item4"`
	Item5                          int                    `json:"item5"`
	Item6                          int                    `json:"item6"`
	ItemsPurchased                 int                    `json:"itemsPurchased"`
	KillingSprees                  int                    `json:"killingSprees"`
	Kills                          int                    `json:"kills"`
	Lane                           string                 `json:"lane"`
	LargestCriticalStrike          int                    `json:"largestCriticalStrike"`
	LargestKillingSpree            int                    `json:"largestKillingSpree"`
	LargestMultiKill               int                    `json:"largestMultiKill"`
	LongestTimeSpentLiving         int                    `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int                    `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int                    `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int                    `json:"magicDamageTaken"`
	NeutralMinionsKilled           int                    `json:"neutralMinionsKilled"`
	NexusKills                     int                    `json:"nexusKills"`
	NexusTakedowns                 int                    `json:"nexusTakedowns"`
	NexusLost                      int                    `json:"nexusLost"`
	ObjectivesStolen               int                    `json:"objectivesStolen"`
	ObjectivesStolenAssists        int                    `json:"objectivesStolenAssists"`
	ParticipantID                  int                    `json:"participantId"`
	PentaKills                     int                    `json:"pentaKills"`
	Perks                          RiotPerks              `json:"perks"`
	PhysicalDamageDealt            int                    `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int                    `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int                    `json:"physicalDamageTaken"`
	ProfileIcon                    int                    `json:"profileIcon"`
	Puuid                          string                 `json:"puuid"`
	QuadraKills                    int                    `json:"quadraKills"`
	RiotIDName                     string                 `json:"riotIdName"`
	RiotIDTagline                  string                 `json:"riotIdTagline"`
	Role                           string                 `json:"role"`
	SightWardsBoughtInGame         int                    `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int                    `json:"spell1Casts"`
	Spell2Casts                    int                    `json:"spell2Casts"`
	Spell3Casts                    int                    `json:"spell3Casts"`
	Spell4Casts                    int                    `json:"spell4Casts"`
	Summoner1Casts                 int                    `json:"summoner1Casts"`
	Summoner1ID                    int                    `json:"summoner1Id"`
	Summoner2Casts                 int                    `json:"summoner2Casts"`
	Summoner2ID                    int                    `json:"summoner2Id"`
	SummonerID                     string                 `json:"summonerId"`
	SummonerLevel                  int                    `json:"summonerLevel"`
	SummonerName                   string                 `json:"summonerName"`
	TeamEarlySurrendered           bool                   `json:"teamEarlySurrendered"`
	TeamID                         int                    `json:"teamId"`
	TeamPosition                   string                 `json:"teamPosition"`
	TimeCCingOthers                int                    `json:"timeCCingOthers"`
	TimePlayed                     int                    `json:"timePlayed"`
	TotalDamageDealt               int                    `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int                    `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int                    `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int                    `json:"totalDamageTaken"`
	TotalHeal                      int                    `json:"totalHeal"`
	TotalHealsOnTeammates          int                    `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int                    `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int                    `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int                    `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int                    `json:"totalUnitsHealed"`
	TripleKills                    int                    `json:"tripleKills"`
	TrueDamageDealt                int                    `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int                    `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int                    `json:"trueDamageTaken"`
	TurretKills                    int                    `json:"turretKills"`
	TurretTakedowns                int                    `json:"turretTakedowns"`
	TurretsLost                    int                    `json:"turretsLost"`
	UnrealKills                    int                    `json:"unrealKills"`
	VisionScore                    int                    `json:"visionScore"`
	VisionWardsBoughtInGame        int                    `json:"visionWardsBoughtInGame"`
	WardsKilled                    int                    `json:"wardsKilled"`
	WardsPlaced                    int                    `json:"wardsPlaced"`
	Win                            bool                   `json:"win"`
}

type RiotPerks struct {
	StatPerks RiotStatPerks      `json:"statPerks"`
	Styles    []RiotPerkStyle    `json:"styles"`
}

type RiotStatPerks struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type RiotPerkStyle struct {
	Description string          `json:"description"`
	Selections  []RiotPerkSelection `json:"selections"`
	Style       int             `json:"style"`
}

type RiotPerkSelection struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type RiotTeam struct {
	Bans       []RiotBan `json:"bans"`
	Objectives RiotObjectives `json:"objectives"`
	TeamID     int       `json:"teamId"`
	Win        bool      `json:"win"`
}

type RiotBan struct {
	ChampionID int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type RiotObjectives struct {
	Baron      RiotObjective `json:"baron"`
	Champion   RiotObjective `json:"champion"`
	Dragon     RiotObjective `json:"dragon"`
	Inhibitor  RiotObjective `json:"inhibitor"`
	RiftHerald RiotObjective `json:"riftHerald"`
	Tower      RiotObjective `json:"tower"`
}

type RiotObjective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

