# Testing Your Deployed Service

Your service is now live on Render! Here's how to test it.

## 1. Get a Match ID

To test the service, you need a real League of Legends match ID from the Riot Games API.

### Option A: Get Match ID via Riot API (Recommended)

1. First, get a Summoner's PUUID:
   ```bash
   # Replace YOUR_RIOT_API_KEY and SUMMONER_NAME
   curl "https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/SUMMONER_NAME?api_key=YOUR_RIOT_API_KEY"
   ```

2. Then get their match history:
   ```bash
   # Replace YOUR_RIOT_API_KEY and PUUID
   curl "https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/PUUID/ids?start=0&count=1&api_key=YOUR_RIOT_API_KEY"
   ```

3. Use the first match ID from the response

### Option B: Use a Known Match ID

If you have a match ID from a recent game, you can use it directly.

## 2. Test the Health Endpoint

```bash
# Replace YOUR_SERVICE_URL with your actual Render URL
curl https://YOUR_SERVICE_URL.onrender.com/health
```

Expected response: `OK`

## 3. Test Basic Match Analysis (GET)

```bash
# Replace YOUR_SERVICE_URL and MATCH_ID
curl "https://YOUR_SERVICE_URL.onrender.com/analyze-match-get?match_id=NA1_1234567890"
```

## 4. Test Match Analysis with POST

```bash
curl -X POST https://YOUR_SERVICE_URL.onrender.com/analyze-match \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "NA1_1234567890"
  }'
```

## 5. Test Champion Deep Dive

```bash
# By champion name
curl -X POST https://YOUR_SERVICE_URL.onrender.com/analyze-match \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "NA1_1234567890",
    "champion_name": "Yasuo"
  }'
```

Or via GET:
```bash
curl "https://YOUR_SERVICE_URL.onrender.com/analyze-match-get?match_id=NA1_1234567890&champion_name=Yasuo"
```

## 6. Test from JavaScript/Web Browser

```javascript
// Basic analysis
const response = await fetch('https://YOUR_SERVICE_URL.onrender.com/analyze-match', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    match_id: 'NA1_1234567890'
  })
});

const data = await response.json();
console.log('Analysis:', data.analysis);
console.log('Suggestions:', data.suggestions);
console.log('Coaching Tips:', data.coaching_tips);
if (data.champion_deep_dive) {
  console.log('Deep Dive:', data.champion_deep_dive);
}
```

## Expected Response Format

```json
{
  "match_id": "NA1_1234567890",
  "analysis": "Comprehensive match analysis...",
  "suggestions": [
    "Suggestion 1",
    "Suggestion 2"
  ],
  "coaching_tips": [
    "Tip 1",
    "Tip 2"
  ],
  "champion_deep_dive": "Detailed champion analysis..." // Only if champion_name or summoner_name was provided
}
```

## Troubleshooting

### Service Spinning Down (Free Tier)

On the free tier, Render spins down services after 15 minutes of inactivity. The first request after spin-down may take 30-60 seconds (cold start).

### Rate Limits

Both Riot Games API and OpenAI API have rate limits:
- **Riot API**: Check your rate limits at https://developer.riotgames.com/
- **OpenAI API**: Check your usage at https://platform.openai.com/usage

### Common Errors

1. **"Failed to fetch match data"**: 
   - Check if match ID is correct
   - Verify RIOT_API_KEY is set correctly in Render
   - Ensure the match region matches RIOT_API_REGION

2. **"Failed to analyze match"**:
   - Check OPENAI_API_KEY is set correctly
   - Verify you have OpenAI API credits
   - Check OpenAI rate limits

3. **"RIOT_API_KEY is required"**:
   - Verify environment variables are set in Render dashboard
   - Check variable names are exact (case-sensitive)

## Next Steps

1. ✅ Test all endpoints
2. Build a frontend to consume this API
3. Add caching for match data (reduce API calls)
4. Set up monitoring and logging
5. Consider upgrading Render plan for always-on service
6. Add rate limiting middleware
7. Set up custom domain (optional)


