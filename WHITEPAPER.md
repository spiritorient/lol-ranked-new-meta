# OldBrickFactory's Match Advisor

## AI-Powered Post-Game Analysis for League of Legends

**Version 1.0**  
**December 2025**

---

## Executive Summary

OldBrickFactory's Match Advisor transforms your League of Legends post-game statistics into actionable coaching insights using artificial intelligence. By combining official Riot Games match data with OpenAI's language models, the platform delivers personalized, data-driven analysis that helps players understand what happened in their games and how to improve.

Enter a match ID, optionally select a champion or summoner to focus on, and receive comprehensive analysis including what went well, what went wrong, critical game moments, item build evaluation, and matchup insights—all backed by your actual match statistics.

---

## The Problem

After a League of Legends match ends, players are left with a scoreboard full of numbers but little context:

- **Raw stats lack meaning**: K/D/A, CS, and damage numbers don't explain *why* you won or lost
- **Self-analysis is biased**: It's hard to objectively evaluate your own performance
- **Generic advice doesn't help**: "Ward more" or "CS better" isn't actionable without specific context
- **Time constraints**: Manually reviewing match data takes significant effort

---

## The Solution

OldBrickFactory's Match Advisor bridges the gap between raw statistics and actionable understanding.

**How it works:**

1. **Enter your Match ID** (e.g., `EUN1_3879610338`)
2. **Optionally specify a champion or summoner** for detailed analysis
3. **Receive AI-generated insights** based on your actual match data

The AI analyzes your specific match—not generic patterns—and provides coaching that references real numbers from your game.

---

## Features

### Match Analysis
Get a comprehensive breakdown of your match including:
- Overall game flow and outcome analysis
- Team performance comparison
- Objective control evaluation
- Strategic insights based on actual events

### Structured Insights

The analysis is organized into clear, actionable sections:

**✅ What Went Well**  
Specific achievements with supporting statistics. Example: "Achieved 8.2 CS/min, outfarming your lane opponent by 45 CS at 20 minutes."

**❌ What Went Wrong**  
Concrete areas for improvement with data. Example: "Died 4 times before 15 minutes, giving the enemy mid laner 1,200 gold advantage."

**⚡ Critical Moments**  
Game-changing events that determined the outcome. Example: "Baron fight at 28 minutes resulted in a team wipe, leading to the game-winning push."

### Deep Dive Analysis

When you specify a champion or summoner name, the platform provides:

- **Performance Metrics**: K/D/A ratio, CS/min, gold/min, damage output
- **Combat Statistics**: Damage breakdown (physical, magic, true), damage taken, healing
- **Objective Control**: Turrets, inhibitors, dragons, barons, first blood/tower
- **Vision Control**: Vision score, wards placed/killed, control wards purchased
- **Lane Matchup**: Direct comparison with your lane opponent
- **Item Build Analysis**: Evaluation of your items against the enemy team composition

### Item & Matchup Analysis

- **Item Timing**: How your build adapted to the game state
- **Opponent Matchup**: How your items countered (or failed to counter) enemy champions
- **Recommendations**: Suggestions for alternative item choices

### Focus Areas

Optionally emphasize specific aspects of analysis:
- Combat statistics
- Vision control
- Objectives
- Item builds
- Matchup & team composition
- Economy
- Farming

---

## Data Source

All analysis is based on official end-of-game statistics from the **Riot Games API v5**.

**Data Retrieved:**
- Match metadata (duration, game mode, version)
- All 10 participants' complete statistics
- Team-level objective counts
- Champion selections and positions
- Final item builds

**Data Limitations:**
- Timeline/event timestamps are not available in the current implementation
- Real-time or in-game analysis is not supported
- Only completed matches can be analyzed

---

## Technology

**Backend**: Go web service  
**AI**: OpenAI GPT-4o-mini (configurable)  
**Data**: Riot Games API v5  
**Frontend**: Web application (HTML/CSS/JavaScript)

The platform uses OpenAI's function calling to generate structured, consistent responses that can be displayed in an interactive tabbed interface.

---

## How to Use

### Web Interface

1. Navigate to the application
2. Enter a Match ID in the format `REGION_NUMBER` (e.g., `EUN1_3879610338`)
3. Optionally enter a Champion Name (e.g., `Nautilus`) or Summoner Name
4. Click "Analyze Match"
5. Browse results across tabs: Overview, What Went Well, What Went Wrong, Critical Moments, Item Analysis, Matchup, and Deep Dive

### API

**Endpoint**: `POST /analyze-match` or `GET /analyze-match-get`

**Request Parameters:**
- `match_id` (required): The match ID to analyze
- `champion_name` (optional): Champion name for focused analysis
- `summoner_name` (optional): Summoner name for focused analysis
- `focus_areas` (optional): Comma-separated list of areas to emphasize

**Example:**
```
GET /analyze-match-get?match_id=EUN1_3879610338&champion_name=Nautilus
```

---

## Response Structure

```json
{
  "match_id": "EUN1_3879610338",
  "analysis": "Comprehensive text analysis of the match...",
  "suggestions": ["Actionable improvement suggestions..."],
  "coaching_tips": ["Strategic tips for future games..."],
  "champion_deep_dive": "Detailed analysis for the specified champion...",
  "structured_insights": {
    "what_went_well": [...],
    "what_went_wrong": [...],
    "critical_moments": [...],
    "item_analysis": {...},
    "matchup_analysis": {...}
  }
}
```

---

## Regions

The platform supports all Riot Games regional routing values:
- `americas` — NA1, BR1, LA1, LA2
- `europe` — EUW1, EUN1, TR1, RU
- `asia` — KR, JP1
- `sea` — PH2, SG2, TH2, TW2, VN2

---

## Privacy

- **No account required**: Just enter a match ID
- **Public data only**: All match data comes from Riot's public API
- **No personal data stored**: Match analysis is generated on-demand

---

## Limitations

- **Post-game only**: Cannot analyze live or ongoing matches
- **End-game statistics**: Analysis is based on final totals, not moment-by-moment events
- **AI interpretation**: Insights are AI-generated and should be considered as coaching suggestions, not absolute truths
- **API rate limits**: Subject to Riot Games and OpenAI API rate limits

---

## Future Considerations

Potential enhancements include:
- Match timeline integration for event-by-event analysis
- Historical trend tracking across multiple matches
- Champion-specific coaching recommendations
- Interactive visualizations and charts

---

## Credits

Built by OldBrickFactory.

**Powered by:**
- [Riot Games API](https://developer.riotgames.com/) — Official League of Legends match data
- [OpenAI](https://openai.com/) — AI-powered analysis generation

---

*OldBrickFactory's Match Advisor isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games, and all associated properties are trademarks or registered trademarks of Riot Games, Inc.*
