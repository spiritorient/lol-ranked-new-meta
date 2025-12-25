# Running the App Locally

Quick guide to run OldBrickFactory's Match Advisor on your local machine.

## Prerequisites

- Go 1.21 or higher installed ([Download Go](https://golang.org/dl/))
- Riot Games API key ([Get one here](https://developer.riotgames.com/))
- OpenAI API key ([Get one here](https://platform.openai.com/))

## Step 1: Install Dependencies

```bash
cd "/Users/neven/Desktop/LOL RANKED NEW META"
go mod download
```

## Step 2: Set Up Environment Variables

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Edit `.env` and add your API keys:

```env
RIOT_API_KEY=your_riot_api_key_here
RIOT_API_REGION=americas
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-4o-mini
PORT=8080
```

**Note:** Never commit your `.env` file (it's already in `.gitignore`).

## Step 3: Run the Server

```bash
go run main.go
```

You should see:
```
Server starting on port 8080
Endpoints available:
  POST /analyze-match - Analyze a match (requires JSON body with match_id)
  GET  /analyze-match-get?match_id=<match_id> - Analyze a match (convenience endpoint)
  GET  /health - Health check
Server starting on :8080
```

## Step 4: Access the App

Open your web browser and go to:
**http://localhost:8080**

You'll see the frontend interface where you can:
- Select a region
- Enter a game ID
- Optionally specify a champion name
- Click "Analyze Match" to get insights

## Testing the API Directly

You can also test the API endpoints directly:

### Health Check
```bash
curl http://localhost:8080/health
```

### Analyze Match (GET)
```bash
curl "http://localhost:8080/analyze-match-get?match_id=EUN1_3879610338&champion_name=Nautilus"
```

### Analyze Match (POST)
```bash
curl -X POST http://localhost:8080/analyze-match \
  -H "Content-Type: application/json" \
  -d '{"match_id": "EUN1_3879610338", "champion_name": "Nautilus"}'
```

## Troubleshooting

### Port Already in Use
If port 8080 is already in use, change the `PORT` in your `.env` file to a different port (e.g., `8081`).

### Missing API Keys
Make sure both `RIOT_API_KEY` and `OPENAI_API_KEY` are set in your `.env` file.

### Go Version Issues
Check your Go version:
```bash
go version
```
Should be 1.21 or higher. If not, update Go.

### Build Errors
If you get build errors, try:
```bash
go mod tidy
go build
```

## Development Tips

- The frontend files are in `frontend/` directory
- Backend code is in the root directory organized by package
- Hot reload: Use tools like `air` or `fresh` for automatic restarts during development

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

