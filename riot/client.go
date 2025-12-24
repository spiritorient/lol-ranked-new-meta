package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lol-ranked-new-meta/types"
)

type Client struct {
	apiKey string
	region string
	client *http.Client
}

// NewClient creates a new Riot API client
func NewClient(apiKey, region string) *Client {
	return &Client{
		apiKey: apiKey,
		region: region,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetMatch fetches match details from the Riot Games API
func (c *Client) GetMatch(matchID string) (*types.RiotMatch, error) {
	// Riot API v5 uses regional routing (americas, europe, asia, sea)
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", c.region, matchID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Riot-Token", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("riot API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var match types.RiotMatch
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &match, nil
}

// FormatMatchForAnalysis converts Riot match data into a format suitable for OpenAI analysis
// championFilter and summonerFilter are optional - if provided, detailed analysis will focus on that champion/summoner
func FormatMatchForAnalysis(match *types.RiotMatch, championFilter, summonerFilter string) string {
	if match == nil {
		return ""
	}

	// Create a summary string with key match information
	var summary string
	summary += fmt.Sprintf("Match ID: %s\n", match.Metadata.MatchID)
	summary += fmt.Sprintf("Game Mode: %s\n", match.Info.GameMode)
	summary += fmt.Sprintf("Game Duration: %d seconds (%.2f minutes)\n", match.Info.GameDuration, float64(match.Info.GameDuration)/60.0)
	summary += fmt.Sprintf("Game Version: %s\n\n", match.Info.GameVersion)

	// Team summaries
	summary += "Teams:\n"
	for _, team := range match.Info.Teams {
		teamName := "Blue"
		if team.TeamID == 200 {
			teamName = "Red"
		}
		result := "Lost"
		if team.Win {
			result = "Won"
		}
		summary += fmt.Sprintf("- Team %s (%s): %d turrets destroyed, %d dragons, %d barons\n",
			teamName, result, team.Objectives.Tower.Kills, team.Objectives.Dragon.Kills, team.Objectives.Baron.Kills)
	}

	summary += "\nParticipants:\n"
	var targetParticipant *types.RiotParticipant
	for i := range match.Info.Participants {
		participant := match.Info.Participants[i]
		teamName := "Blue"
		if participant.TeamID == 200 {
			teamName = "Red"
		}
		result := "Lost"
		if participant.Win {
			result = "Won"
		}

		// Check if this is the target participant for deep dive
		isTarget := false
		if championFilter != "" && participant.ChampionName == championFilter {
			isTarget = true
			targetParticipant = &match.Info.Participants[i]
		}
		if summonerFilter != "" && participant.SummonerName == summonerFilter {
			isTarget = true
			targetParticipant = &match.Info.Participants[i]
		}

		marker := ""
		if isTarget {
			marker = " [TARGET FOR DEEP DIVE]"
		}

		summary += fmt.Sprintf("- %s (%s, %s, %s)%s: K/D/A: %d/%d/%d, CS: %d, Gold: %d, Damage: %d\n",
			participant.SummonerName,
			participant.ChampionName,
			teamName,
			result,
			marker,
			participant.Kills,
			participant.Deaths,
			participant.Assists,
			participant.TotalMinionsKilled,
			participant.GoldEarned,
			participant.TotalDamageDealtToChampions,
		)
	}

	// If we have a target participant, add detailed stats
	if targetParticipant != nil {
		summary += "\n=== DETAILED STATS FOR TARGET PLAYER ===\n"
		summary += FormatParticipantDeepDive(targetParticipant, match.Info.GameDuration)
	}

	return summary
}

// FormatParticipantDeepDive creates a detailed analysis string for a specific participant
func FormatParticipantDeepDive(participant *types.RiotParticipant, gameDuration int64) string {
	if participant == nil {
		return ""
	}

	var detail string
	detail += fmt.Sprintf("Summoner: %s (%s#%s)\n", participant.SummonerName, participant.RiotIDName, participant.RiotIDTagline)
	detail += fmt.Sprintf("Champion: %s (Level %d)\n", participant.ChampionName, participant.ChampLevel)
	detail += fmt.Sprintf("Team Position: %s (Lane: %s, Role: %s)\n", participant.TeamPosition, participant.Lane, participant.Role)
	detail += fmt.Sprintf("Result: %s\n\n", map[bool]string{true: "Victory", false: "Defeat"}[participant.Win])

	detail += "Performance Metrics:\n"
	detail += fmt.Sprintf("- K/D/A: %d/%d/%d (KDA Ratio: %.2f)\n",
		participant.Kills, participant.Deaths, participant.Assists,
		float64(participant.Kills+participant.Assists)/float64(max(participant.Deaths, 1)))
	detail += fmt.Sprintf("- CS: %d (%.1f CS/min)\n", participant.TotalMinionsKilled,
		float64(participant.TotalMinionsKilled)/(float64(gameDuration)/60.0))
	detail += fmt.Sprintf("- Gold Earned: %d (Gold/min: %.0f)\n", participant.GoldEarned,
		float64(participant.GoldEarned)/(float64(gameDuration)/60.0))
	detail += fmt.Sprintf("- Gold Spent: %d\n", participant.GoldSpent)

	detail += "\nCombat Stats:\n"
	detail += fmt.Sprintf("- Total Damage to Champions: %d\n", participant.TotalDamageDealtToChampions)
	detail += fmt.Sprintf("- Physical Damage: %d\n", participant.PhysicalDamageDealtToChampions)
	detail += fmt.Sprintf("- Magic Damage: %d\n", participant.MagicDamageDealtToChampions)
	detail += fmt.Sprintf("- True Damage: %d\n", participant.TrueDamageDealtToChampions)
	detail += fmt.Sprintf("- Damage Taken: %d\n", participant.TotalDamageTaken)
	detail += fmt.Sprintf("- Damage Self Mitigated: %d\n", participant.DamageSelfMitigated)
	detail += fmt.Sprintf("- Total Heal: %d\n", participant.TotalHeal)
	detail += fmt.Sprintf("- Total Shields on Teammates: %d\n", participant.TotalDamageShieldedOnTeammates)

	detail += "\nObjective Control:\n"
	detail += fmt.Sprintf("- Turret Kills: %d\n", participant.TurretKills)
	detail += fmt.Sprintf("- Inhibitor Kills: %d\n", participant.InhibitorKills)
	detail += fmt.Sprintf("- Dragon Kills: %d\n", participant.DragonKills)
	detail += fmt.Sprintf("- Baron Kills: %d\n", participant.BaronKills)
	detail += fmt.Sprintf("- First Blood: %s\n", map[bool]string{true: "Yes", false: "No"}[participant.FirstBloodKill])
	detail += fmt.Sprintf("- First Tower: %s\n", map[bool]string{true: "Yes", false: "No"}[participant.FirstTowerKill])

	detail += "\nVision & Map Control:\n"
	detail += fmt.Sprintf("- Vision Score: %d\n", participant.VisionScore)
	detail += fmt.Sprintf("- Wards Placed: %d\n", participant.WardsPlaced)
	detail += fmt.Sprintf("- Wards Killed: %d\n", participant.WardsKilled)
	detail += fmt.Sprintf("- Control Wards Purchased: %d\n", participant.VisionWardsBoughtInGame)
	detail += fmt.Sprintf("- Detector Wards Placed: %d\n", participant.DetectorWardsPlaced)

	detail += "\nSpecial Achievements:\n"
	detail += fmt.Sprintf("- Largest Killing Spree: %d\n", participant.LargestKillingSpree)
	detail += fmt.Sprintf("- Killing Sprees: %d\n", participant.KillingSprees)
	detail += fmt.Sprintf("- Double Kills: %d\n", participant.DoubleKills)
	detail += fmt.Sprintf("- Triple Kills: %d\n", participant.TripleKills)
	detail += fmt.Sprintf("- Quadra Kills: %d\n", participant.QuadraKills)
	detail += fmt.Sprintf("- Penta Kills: %d\n", participant.PentaKills)
	detail += fmt.Sprintf("- Unreal Kills: %d\n", participant.UnrealKills)
	detail += fmt.Sprintf("- Largest Multi Kill: %d\n", participant.LargestMultiKill)

	detail += "\nItem Build:\n"
	items := []int{participant.Item0, participant.Item1, participant.Item2, participant.Item3, participant.Item4, participant.Item5, participant.Item6}
	for i, itemID := range items {
		if itemID != 0 {
			slotName := "Item"
			if i == 6 {
				slotName = "Trinket"
			}
			detail += fmt.Sprintf("- %s %d: %d\n", slotName, i+1, itemID)
		}
	}
	detail += fmt.Sprintf("- Total Items Purchased: %d\n", participant.ItemsPurchased)

	detail += "\nSummoner Spells:\n"
	detail += fmt.Sprintf("- Summoner Spell 1 (ID %d): Used %d times\n", participant.Summoner1ID, participant.Summoner1Casts)
	detail += fmt.Sprintf("- Summoner Spell 2 (ID %d): Used %d times\n", participant.Summoner2ID, participant.Summoner2Casts)

	detail += "\nGame Impact:\n"
	detail += fmt.Sprintf("- Time Spent Dead: %d seconds\n", participant.TotalTimeSpentDead)
	detail += fmt.Sprintf("- Longest Time Spent Living: %d seconds\n", participant.LongestTimeSpentLiving)
	detail += fmt.Sprintf("- Time CC'd Others: %d seconds\n", participant.TimeCCingOthers)
	detail += fmt.Sprintf("- Total Time CC'd: %d seconds\n", participant.TotalTimeCCDealt)

	return detail
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
