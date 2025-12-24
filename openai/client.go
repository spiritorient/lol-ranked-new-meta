package openai

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"lol-ranked-new-meta/types"
)

type Client struct {
	client *openai.Client
	model  string
}

// NewClient creates a new OpenAI client
func NewClient(apiKey, model string) *Client {
	return &Client{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

// AnalyzeMatch analyzes a League of Legends match and provides coaching advice
// championFilter and summonerFilter are optional - if provided, will generate a deep dive analysis
func (c *Client) AnalyzeMatch(ctx context.Context, matchSummary string, championFilter, summonerFilter string) (*types.MatchResponse, error) {
	// Define the function schema for structured output
	analyzeMatchFunction := openai.FunctionDefinition{
		Name:        "analyze_match",
		Description: "Analyzes a League of Legends match and provides detailed coaching advice, suggestions, and tips",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"analysis": map[string]interface{}{
					"type":        "string",
					"description": "A comprehensive analysis of the match performance, key moments, and overall game flow",
				},
				"suggestions": map[string]interface{}{
					"type":        "array",
					"description": "List of actionable suggestions for improvement based on match data",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"coaching_tips": map[string]interface{}{
					"type":        "array",
					"description": "List of coaching tips and strategies for future matches",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			"required": []string{"analysis", "suggestions", "coaching_tips"},
		},
	}

	systemPrompt := `You are an expert League of Legends coach providing DATA-DRIVEN, SPECIFIC analysis.
CRITICAL: Analyze the ACTUAL events from this EXACT match using the provided data.

Your analysis must:
- Reference specific numbers, stats, and achievements from this match
- Identify what ACTUALLY happened, not generic patterns
- Compare actual performance vs opponents using real data
- Explain WHY specific events mattered based on the match outcome
- Focus on concrete, measurable events that occurred in this game

Avoid generic advice. Instead, cite specific stats like "Team secured 3 dragons at 15, 22, and 28 minutes" or "ADC had 300 CS at 30 minutes vs opponent's 250".`

	userPrompt := fmt.Sprintf(`Analyze this EXACT League of Legends match using the specific data provided:

%s

Provide analysis that references SPECIFIC NUMBERS, EVENTS, and STATS from this match. 
Focus on what actually happened, not generic coaching advice.`, matchSummary)

	// Create the chat completion request with function calling
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Functions: []openai.FunctionDefinition{analyzeMatchFunction},
		FunctionCall: map[string]interface{}{
			"name": "analyze_match",
		},
		Temperature: 0.7,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := resp.Choices[0]

	// Parse the function call arguments
	var functionArgs map[string]interface{}
	if choice.Message.FunctionCall != nil && choice.Message.FunctionCall.Arguments != "" {
		// Parse the JSON arguments
		if err := json.Unmarshal([]byte(choice.Message.FunctionCall.Arguments), &functionArgs); err != nil {
			// If parsing fails, try to extract information from the message content instead
			return c.extractFromContent(choice.Message.Content), nil
		}

		response := &types.MatchResponse{}

		if analysis, ok := functionArgs["analysis"].(string); ok {
			response.Analysis = analysis
		}

		if suggestions, ok := functionArgs["suggestions"].([]interface{}); ok {
			response.Suggestions = make([]string, len(suggestions))
			for i, s := range suggestions {
				if str, ok := s.(string); ok {
					response.Suggestions[i] = str
				}
			}
		}

		if tips, ok := functionArgs["coaching_tips"].([]interface{}); ok {
			response.CoachingTips = make([]string, len(tips))
			for i, t := range tips {
				if str, ok := t.(string); ok {
					response.CoachingTips[i] = str
				}
			}
		}

		// If champion/summoner filter is specified, generate deep dive analysis and structured insights
		if championFilter != "" || summonerFilter != "" {
			deepDive, err := c.AnalyzeChampionDeepDive(ctx, matchSummary, championFilter, summonerFilter)
			if err != nil {
				// Log error but don't fail the whole request
				response.ChampionDeepDive = "Failed to generate deep dive analysis: " + err.Error()
			} else {
				response.ChampionDeepDive = deepDive
			}
			
			// Generate structured insights for interactive frontend
			structuredInsights, err := c.GenerateStructuredInsights(ctx, matchSummary, championFilter, summonerFilter)
			if err == nil {
				response.StructuredInsights = structuredInsights
			}
		}

		return response, nil
	}

	// Fallback: extract from content if function calling didn't work as expected
	response := c.extractFromContent(choice.Message.Content)

	// If champion/summoner filter is specified, generate deep dive analysis and structured insights
	if championFilter != "" || summonerFilter != "" {
		deepDive, err := c.AnalyzeChampionDeepDive(ctx, matchSummary, championFilter, summonerFilter)
		if err != nil {
			// Log error but don't fail the whole request
			response.ChampionDeepDive = "Failed to generate deep dive analysis: " + err.Error()
		} else {
			response.ChampionDeepDive = deepDive
		}
		
		// Generate structured insights for interactive frontend
		structuredInsights, err := c.GenerateStructuredInsights(ctx, matchSummary, championFilter, summonerFilter)
		if err == nil {
			response.StructuredInsights = structuredInsights
		}
	}

	return response, nil
}

// AnalyzeChampionDeepDive provides a detailed analysis focused on a specific champion
func (c *Client) AnalyzeChampionDeepDive(ctx context.Context, matchSummary, championFilter, summonerFilter string) (string, error) {
	targetName := championFilter
	if summonerFilter != "" {
		targetName = summonerFilter
	}

	systemPrompt := `You are an expert League of Legends coach specializing in data-driven, specific match analysis. 
CRITICAL: Focus on ACTUAL EVENTS and SPECIFIC DATA from this exact match, not generic archetypical advice.

Your analysis must:
- Reference specific numbers, stats, and events from the match data provided
- Explain what ACTUALLY happened, not what "usually" happens
- Compare actual performance to opponent's actual performance using the data
- Analyze item builds in context of the actual opponent champions faced
- Identify concrete mistakes using specific match statistics
- Highlight specific good plays using actual numbers and achievements

Avoid generic advice like "ward more" - instead say "placed only X wards compared to opponent's Y" with specific impact.`

	userPrompt := fmt.Sprintf(`Analyze the performance of %s in this EXACT match. Use the actual data provided.

MATCH DATA:
%s

Provide analysis focusing on SPECIFIC EVENTS AND NUMBERS from this match:
1. What went well - cite specific stats (e.g., "Achieved 8.5 CS/min at 15 minutes, above average")
2. What went wrong - cite specific failures (e.g., "Died 5 times before 10 minutes, giving enemy ADC 1500 gold")
3. Critical moments - identify specific game-changing events using the data
4. Item build analysis - evaluate items purchased in context of actual opponent champions
5. Matchup performance - compare actual stats vs lane opponent (provided in data)
6. Specific, actionable improvements based on this exact match's data`, targetName, matchSummary)

	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: 0.7,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create deep dive analysis: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateStructuredInsights creates structured, data-driven insights for interactive frontend
func (c *Client) GenerateStructuredInsights(ctx context.Context, matchSummary, championFilter, summonerFilter string) (*types.StructuredInsights, error) {
	targetName := championFilter
	if summonerFilter != "" {
		targetName = summonerFilter
	}

	// Define function schema for structured insights
	structuredFunction := openai.FunctionDefinition{
		Name:        "generate_structured_insights",
		Description: "Generates structured, data-driven insights about a League of Legends match with specific events and statistics",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"what_went_well": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"title":       map[string]interface{}{"type": "string"},
							"description": map[string]interface{}{"type": "string"},
							"impact":      map[string]interface{}{"type": "string"},
							"data":        map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
							"category":    map[string]interface{}{"type": "string"},
						},
					},
				},
				"what_went_wrong": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"title":       map[string]interface{}{"type": "string"},
							"description": map[string]interface{}{"type": "string"},
							"impact":      map[string]interface{}{"type": "string"},
							"data":        map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
							"category":    map[string]interface{}{"type": "string"},
						},
					},
				},
				"critical_moments": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"title":       map[string]interface{}{"type": "string"},
							"description": map[string]interface{}{"type": "string"},
							"outcome":     map[string]interface{}{"type": "string"},
							"impact":      map[string]interface{}{"type": "string"},
							"data":        map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						},
					},
				},
				"item_analysis": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"timing_analysis":  map[string]interface{}{"type": "string"},
						"opponent_matchup": map[string]interface{}{"type": "string"},
						"recommendations":  map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					},
				},
				"matchup_analysis": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"lane_matchup":    map[string]interface{}{"type": "string"},
						"team_composition": map[string]interface{}{"type": "string"},
						"synergies":       map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"counters":        map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"win_conditions":  map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					},
				},
			},
			"required": []string{"what_went_well", "what_went_wrong", "critical_moments"},
		},
	}

	systemPrompt := `You are an expert League of Legends analyst. Generate STRUCTURED insights based on ACTUAL match data.
CRITICAL: Only reference specific numbers, stats, and events from the provided match data.
Each insight must cite actual data (e.g., "Died 3 times before 10 minutes" not "died early").`

	userPrompt := fmt.Sprintf(`Generate structured insights for %s in this match. Use ONLY the actual data provided:

%s

Return structured data with:
1. What went well - specific achievements with numbers
2. What went wrong - specific failures with supporting data
3. Critical moments - game-changing events with context
4. Item analysis - evaluate items vs actual opponent champions
5. Matchup analysis - compare actual performance vs lane opponent`, targetName, matchSummary)

	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		Functions: []openai.FunctionDefinition{structuredFunction},
		FunctionCall: map[string]interface{}{"name": "generate_structured_insights"},
		Temperature: 0.3, // Lower temperature for more consistent structured output
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate structured insights: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := resp.Choices[0]
	if choice.Message.FunctionCall != nil && choice.Message.FunctionCall.Arguments != "" {
		var insights types.StructuredInsights
		if err := json.Unmarshal([]byte(choice.Message.FunctionCall.Arguments), &insights); err == nil {
			return &insights, nil
		}
	}

	return nil, fmt.Errorf("failed to parse structured insights")
}

// extractFromContent is a fallback method to parse analysis from content
func (c *Client) extractFromContent(content string) *types.MatchResponse {
	// This is a simple fallback - in production, you might want more sophisticated parsing
	return &types.MatchResponse{
		Analysis:     content,
		Suggestions:  []string{"Review the detailed analysis above for specific suggestions"},
		CoachingTips: []string{"Focus on the key areas mentioned in the analysis"},
	}
}

