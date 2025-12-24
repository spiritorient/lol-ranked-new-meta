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

	systemPrompt := `You are an expert League of Legends coach and analyst. Your task is to analyze match data and provide:
1. A comprehensive analysis of the match, highlighting key moments, strengths, and weaknesses
2. Specific, actionable suggestions for improvement
3. Coaching tips for future matches

Focus on practical advice that can help players improve their gameplay. Consider factors like:
- Team composition and synergy
- Objective control (dragons, barons, towers)
- Individual performance (K/D/A, CS, gold, damage)
- Game timing and decision-making
- Vision control and map awareness`

	userPrompt := fmt.Sprintf(`Please analyze this League of Legends match data and provide detailed coaching advice:

%s

Provide your analysis, suggestions, and coaching tips.`, matchSummary)

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

		// If champion/summoner filter is specified, generate deep dive analysis
		if championFilter != "" || summonerFilter != "" {
			deepDive, err := c.AnalyzeChampionDeepDive(ctx, matchSummary, championFilter, summonerFilter)
			if err != nil {
				// Log error but don't fail the whole request
				response.ChampionDeepDive = "Failed to generate deep dive analysis: " + err.Error()
			} else {
				response.ChampionDeepDive = deepDive
			}
		}

		return response, nil
	}

	// Fallback: extract from content if function calling didn't work as expected
	response := c.extractFromContent(choice.Message.Content)

	// If champion/summoner filter is specified, generate deep dive analysis
	if championFilter != "" || summonerFilter != "" {
		deepDive, err := c.AnalyzeChampionDeepDive(ctx, matchSummary, championFilter, summonerFilter)
		if err != nil {
			// Log error but don't fail the whole request
			response.ChampionDeepDive = "Failed to generate deep dive analysis: " + err.Error()
		} else {
			response.ChampionDeepDive = deepDive
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

	systemPrompt := `You are an expert League of Legends coach specializing in detailed champion performance analysis. 
Your task is to provide a comprehensive, in-depth analysis of a specific player's performance in a match.
Focus on:
- Champion-specific mechanics and execution
- Decision-making patterns throughout the game
- Itemization choices and build path effectiveness
- Positioning and map awareness
- Team fight participation and impact
- Farming patterns and resource management
- Vision control and warding patterns
- Comparison with expected performance for that champion/role
- Specific areas for improvement with actionable advice`

	userPrompt := fmt.Sprintf(`Please provide a detailed deep dive analysis for the player/champion marked as "[TARGET FOR DEEP DIVE]" in the following match data:

%s

Focus specifically on %s's performance. Provide insights on:
1. What they did well
2. Critical mistakes or missed opportunities
3. Champion-specific mechanics and combos
4. Item build analysis
5. Specific coaching points for improvement
6. Role-specific recommendations`, matchSummary, targetName)

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

// extractFromContent is a fallback method to parse analysis from content
func (c *Client) extractFromContent(content string) *types.MatchResponse {
	// This is a simple fallback - in production, you might want more sophisticated parsing
	return &types.MatchResponse{
		Analysis:     content,
		Suggestions:  []string{"Review the detailed analysis above for specific suggestions"},
		CoachingTips: []string{"Focus on the key areas mentioned in the analysis"},
	}
}

