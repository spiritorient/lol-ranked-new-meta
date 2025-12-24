# League of Legends Match Advisor

A Go backend service that analyzes League of Legends matches using the Riot Games API and provides coaching advice through OpenAI's function-calling capabilities.

## Features

- Fetches match data from Riot Games API
- Analyzes matches using OpenAI GPT models
- Provides detailed analysis, suggestions, and coaching tips
- **Champion Deep Dive**: Optional detailed analysis focused on a specific champion or player
- Comprehensive participant statistics including combat, objectives, vision, and itemization
- RESTful API endpoints for easy integration

## Prerequisites

- Go 1.21 or higher
- Riot Games API key ([Get one here](https://developer.riotgames.com/))
- OpenAI API key ([Get one here](https://platform.openai.com/))

## Setup

1. **Clone or navigate to the project directory**

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Configure environment variables:**
   Copy the `.env.example` file to `.env` and fill in your API keys:
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your actual API keys:
   ```
   RIOT_API_KEY=your_riot_api_key_here
   RIOT_API_REGION=americas  # Options: americas, europe, asia, sea
   OPENAI_API_KEY=your_openai_api_key_here
   OPENAI_MODEL=gpt-4o-mini  # Or gpt-4, gpt-3.5-turbo, etc.
   PORT=8080
   ```

   **Note:** You can also set these as environment variables directly without using a `.env` file.

## Running the Server

Start the server:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in your configuration).

## API Endpoints

### POST /analyze-match

Analyzes a match and returns coaching advice.

**Request Body:**
```json
{
  "match_id": "NA1_1234567890",
  "champion_name": "Yasuo",  // Optional: for deep dive analysis on specific champion
  "summoner_name": "PlayerName"  // Optional: for deep dive analysis on specific summoner
}
```

**Note:** You can specify either `champion_name` OR `summoner_name` for a deep dive analysis. If provided, the response will include a `champion_deep_dive` field with detailed analysis focused on that specific player/champion.

**Response:**
```json
{
  "match_id": "NA1_1234567890",
  "analysis": "Detailed analysis of the match...",
  "suggestions": [
    "Suggestion 1",
    "Suggestion 2"
  ],
  "coaching_tips": [
    "Coaching tip 1",
    "Coaching tip 2"
  ],
  "champion_deep_dive": "Detailed deep dive analysis focusing on the specified champion/player..." // Only present if champion_name or summoner_name was provided
}
```

### GET /analyze-match-get?match_id=<match_id>

Convenience GET endpoint for testing.

**Query Parameters:**
- `match_id` (required): The match ID to analyze
- `champion_name` (optional): Champion name for deep dive analysis (e.g., "Yasuo", "Jinx")
- `summoner_name` (optional): Summoner name for deep dive analysis

**Examples:**
```bash
# Basic analysis
curl "http://localhost:8080/analyze-match-get?match_id=NA1_1234567890"

# Deep dive on specific champion
curl "http://localhost:8080/analyze-match-get?match_id=NA1_1234567890&champion_name=Yasuo"

# Deep dive on specific summoner
curl "http://localhost:8080/analyze-match-get?match_id=NA1_1234567890&summoner_name=PlayerName"
```

### GET /health

Health check endpoint.

## Usage Examples

### Using cURL

**POST request (basic analysis):**
```bash
curl -X POST http://localhost:8080/analyze-match \
  -H "Content-Type: application/json" \
  -d '{"match_id": "NA1_1234567890"}'
```

**POST request (with champion deep dive):**
```bash
curl -X POST http://localhost:8080/analyze-match \
  -H "Content-Type: application/json" \
  -d '{"match_id": "NA1_1234567890", "champion_name": "Yasuo"}'
```

**GET request (basic analysis):**
```bash
curl "http://localhost:8080/analyze-match-get?match_id=NA1_1234567890"
```

**GET request (with champion deep dive):**
```bash
curl "http://localhost:8080/analyze-match-get?match_id=NA1_1234567890&champion_name=Yasuo"
```

### Using JavaScript (Fetch API)

**Basic analysis:**
```javascript
const response = await fetch('http://localhost:8080/analyze-match', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    match_id: 'NA1_1234567890'
  })
});

const data = await response.json();
console.log(data);
```

**With champion deep dive:**
```javascript
const response = await fetch('http://localhost:8080/analyze-match', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    match_id: 'NA1_1234567890',
    champion_name: 'Yasuo'  // or summoner_name: 'PlayerName'
  })
});

const data = await response.json();
console.log(data.analysis);           // General match analysis
console.log(data.champion_deep_dive); // Deep dive analysis (if requested)
```

## Project Structure

```
.
├── config/          # Configuration management
├── handlers/        # HTTP request handlers
├── openai/          # OpenAI integration
├── riot/            # Riot Games API client
├── types/           # Shared type definitions
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
└── README.md        # This file
```

## Riot API Regions

The `RIOT_API_REGION` should be set to one of the following regional routing values:
- `americas` - For NA1, BR1, LA1, LA2
- `europe` - For EUW1, EUN1, TR1, RU
- `asia` - For KR, JP1
- `sea` - For PH2, SG2, TH2, TW2, VN2

## Notes

- Rate Limits: Both Riot Games and OpenAI APIs have rate limits. Make sure to handle these appropriately in production.
- Match IDs: Match IDs should be in the format returned by the Riot API (e.g., `NA1_1234567890`).
- API Keys: Never commit your `.env` file to version control. It's already included in `.gitignore`.

## Champion Deep Dive Feature

When you specify a `champion_name` or `summoner_name` in your request, the API will provide:

1. **Enhanced Match Summary**: The target player/champion is marked in the match summary with detailed statistics
2. **Detailed Statistics**: Comprehensive breakdown including:
   - Performance metrics (K/D/A, CS, Gold, Damage)
   - Combat statistics (damage breakdown, damage taken, healing)
   - Objective control (turrets, dragons, barons)
   - Vision control (wards placed/killed, vision score)
   - Item build analysis
   - Special achievements (multi-kills, killing sprees)
3. **AI-Powered Deep Dive Analysis**: OpenAI provides:
   - Champion-specific mechanics and execution analysis
   - Decision-making pattern evaluation
   - Itemization effectiveness review
   - Positioning and map awareness assessment
   - Specific areas for improvement with actionable advice

The deep dive analysis appears in the `champion_deep_dive` field of the response.

## Future Enhancements

- Real-time match analysis for live games
- Caching of match data
- Rate limiting middleware
- Database integration for storing match analyses
- WebSocket support for real-time updates
- Multiple champion comparison analysis

## License

This project is for educational purposes.

