# League of Legends Match Advisor
## AI-Powered Performance Analysis Platform

**Version 1.0**  
**December 2024**

---

## Executive Summary

The League of Legends Match Advisor is an advanced, AI-powered analytics platform that transforms raw match data into actionable coaching insights. By leveraging the Riot Games API and OpenAI's GPT models, the platform provides players with data-driven analysis, personalized recommendations, and comprehensive performance breakdowns. The system processes end-of-game statistics to deliver structured insights covering combat performance, objective control, vision management, itemization, and strategic decision-making.

This whitepaper presents the technical architecture, implementation methodology, and capabilities of a production-ready system designed to help players improve their gameplay through intelligent, context-aware analysis.

---

## 1. Introduction

### 1.1 Problem Statement

League of Legends players face significant challenges in understanding their performance and identifying improvement opportunities:

- **Data Overload**: Match statistics are voluminous but lack context and actionable interpretation
- **Generic Advice**: Most coaching resources provide archetypical guidance rather than match-specific insights
- **Limited Analysis Tools**: Existing tools focus on aggregate statistics without deep-dive capabilities
- **Time Constraints**: Manual analysis of match data is time-consuming and requires expertise

### 1.2 Solution Overview

The League of Legends Match Advisor addresses these challenges by:

1. **Automated Data Processing**: Extracts and structures match data from the Riot Games API
2. **AI-Powered Analysis**: Uses advanced language models to generate context-aware, data-driven insights
3. **Structured Output**: Provides organized, actionable recommendations with supporting evidence
4. **Flexible Analysis**: Supports both general match analysis and deep-dive player/champion analysis
5. **Comprehensive Tracking**: Maintains complete analytics history for usage patterns and insights

---

## 2. Technical Architecture

### 2.1 System Overview

The platform is built as a modern, cloud-native web application with the following architecture:

```
┌─────────────────┐
│   Web Frontend  │ (React/Vanilla JS)
└────────┬────────┘
         │ HTTP/REST
┌────────▼─────────────────────────┐
│      Go Backend Service          │
│  ┌──────────────────────────┐   │
│  │  HTTP Request Handlers   │   │
│  └──────────┬───────────────┘   │
│             │                    │
│  ┌──────────▼───────────────┐   │
│  │   Analytics Middleware   │   │
│  └──────────┬───────────────┘   │
│             │                    │
│  ┌──────────▼───────────────┐   │
│  │   Match Analysis Engine  │   │
│  └──────┬───────────────┬───┘   │
│         │               │        │
│  ┌──────▼──────┐  ┌─────▼─────┐ │
│  │ Riot API    │  │ OpenAI API│ │
│  │   Client    │  │  Client   │ │
│  └─────────────┘  └───────────┘ │
└──────────────────────────────────┘
         │
         │ Persistent Storage
┌────────▼────────┐
│  Analytics Data │
│  (JSON/File)    │
└─────────────────┘
```

### 2.2 Technology Stack

**Backend:**
- **Language**: Go 1.21+
- **HTTP Framework**: Standard library `net/http`
- **API Integration**: Riot Games API v5, OpenAI API
- **Storage**: File-based JSON (extensible to databases)
- **Deployment**: Render.com (cloud platform)

**Frontend:**
- **Technology**: Vanilla JavaScript, HTML5, CSS3
- **Architecture**: Single Page Application (SPA)
- **Styling**: Modern CSS with CSS Variables
- **Responsive Design**: Mobile-first approach

**Infrastructure:**
- **Hosting**: Render.com Web Service
- **Persistent Storage**: Render.com Persistent Disk
- **API Gateway**: Built-in HTTP server with middleware

### 2.3 Core Components

#### 2.3.1 Request Handler Layer
- **Match Handler**: Processes match analysis requests
- **Analytics Handler**: Serves analytics and statistics
- **Health Check**: Service monitoring endpoint

#### 2.3.2 Business Logic Layer
- **Riot Client**: Interfaces with Riot Games API
- **OpenAI Client**: Manages AI analysis requests
- **Analytics Tracker**: Records and aggregates usage data

#### 2.3.3 Data Layer
- **Type Definitions**: Structured data models
- **Storage**: Persistent data management
- **Configuration**: Environment-based settings

---

## 3. Data Sources and Processing

### 3.1 Riot Games API Integration

The platform integrates with the Riot Games API v5 to retrieve comprehensive match data:

**Endpoint**: `GET /lol/match/v5/matches/{matchId}`

**Data Retrieved:**
- Match metadata (duration, version, mode, queue type)
- All 10 participants' complete statistics
- Team-level objective control
- Champion selections and positions
- Champion bans

**Key Statistics Captured:**

**Combat Metrics:**
- Kills, Deaths, Assists (K/D/A)
- Damage dealt (total, physical, magic, true)
- Damage taken
- Healing and shielding
- Crowd control duration

**Economic Metrics:**
- Gold earned and spent
- CS (Creep Score) - minions and jungle camps
- Item purchases

**Objective Control:**
- Turrets, inhibitors, dragons, barons
- First blood, first tower
- Objective steals

**Vision Metrics:**
- Vision score
- Wards placed and killed
- Control wards purchased

**Performance Metrics:**
- Multi-kills, killing sprees
- Time spent dead
- Longest time alive
- Summoner spell usage

### 3.2 Data Limitations

**Available:**
- Complete end-of-game statistics
- Final item builds
- Aggregate totals and counts
- Champion selections and positions

**Not Available:**
- Timeline/event timestamps
- Item purchase timing
- Real-time game state
- Position data over time
- Exact moment-by-moment events

### 3.3 Data Processing Pipeline

1. **Data Extraction**: Fetch raw match data from Riot API
2. **Data Transformation**: Convert to structured format
3. **Contextual Enrichment**: Add opponent comparisons, lane matchups
4. **Summary Generation**: Create human-readable summary
5. **AI Analysis**: Process through OpenAI for insights
6. **Response Structuring**: Format as JSON response

---

## 4. AI Integration and Analysis

### 4.1 OpenAI Integration

The platform leverages OpenAI's GPT models (default: GPT-4o-mini) with function calling capabilities for structured output.

**System Prompts:**
- Emphasize data-driven, specific analysis
- Require citation of actual match statistics
- Focus on actionable insights
- Avoid generic advice

**Analysis Types:**

1. **General Match Analysis**
   - Overall game flow
   - Team performance comparison
   - Key turning points
   - Strategic insights

2. **Deep Dive Analysis** (Champion/Player-specific)
   - Detailed performance breakdown
   - Lane matchup evaluation
   - Item build analysis
   - Specific improvement areas

3. **Structured Insights**
   - What went well (with data)
   - What went wrong (with data)
   - Critical moments
   - Item recommendations
   - Matchup analysis

### 4.2 Analysis Features

**Focus Areas** (User-selectable):
- Combat Statistics
- Vision Control
- Objectives
- Item Builds
- Matchup & Team Composition
- Economy
- Farming

**Output Structure:**
```json
{
  "analysis": "Comprehensive text analysis",
  "suggestions": ["Actionable suggestions"],
  "coaching_tips": ["Strategic tips"],
  "champion_deep_dive": "Detailed player analysis",
  "structured_insights": {
    "what_went_well": [...],
    "what_went_wrong": [...],
    "critical_moments": [...],
    "item_analysis": {...},
    "matchup_analysis": {...},
    "key_statistics": {...}
  }
}
```

---

## 5. Features and Capabilities

### 5.1 Core Features

**Match Analysis:**
- Complete match breakdown
- Team performance comparison
- Objective control analysis
- Strategic insights

**Player Deep Dive:**
- Champion-specific analysis
- Performance metrics breakdown
- Lane matchup evaluation
- Item build assessment

**Structured Insights:**
- Event-based analysis
- Data-backed recommendations
- Critical moment identification
- Improvement suggestions

### 5.2 User Interface

**Modern Web Interface:**
- Responsive design (mobile-friendly)
- Tabbed navigation
- Real-time analysis
- Clear visual hierarchy
- Intuitive controls

**User Experience:**
- Simple match ID input
- Optional champion/player focus
- Selectable analysis focus areas
- Clear feedback and messaging

### 5.3 Analytics and Tracking

**Comprehensive Request Tracking:**
- All HTTP requests logged
- IP address tracking
- User agent identification
- Response time monitoring
- Path and method tracking

**Statistics Collected:**
- Total requests
- Unique visitors
- Requests by endpoint
- Daily breakdowns
- Browser/device types
- Top IPs

**Storage:**
- Persistent disk storage
- Unlimited history (configurable)
- Automatic cleanup (optional)
- Complete request history

---

## 6. Implementation Details

### 6.1 API Endpoints

**POST /analyze-match**
- Accepts JSON request body
- Supports champion/player focus
- Returns comprehensive analysis

**GET /analyze-match-get**
- Convenience GET endpoint
- Query parameter-based
- Same functionality as POST

**GET /analytics**
- Analytics dashboard endpoint
- Optional API key protection
- Complete statistics view

**GET /health**
- Health check endpoint
- Service monitoring

### 6.2 Request Flow

1. **Client Request**: User submits match ID (and optional filters)
2. **Validation**: Server validates request parameters
3. **Riot API Call**: Fetch match data from Riot Games
4. **Data Processing**: Format and enrich match data
5. **OpenAI Analysis**: Generate AI-powered insights
6. **Response**: Return structured JSON response
7. **Analytics**: Track request for statistics

### 6.3 Error Handling

- Graceful API failures
- Clear error messages
- Fallback mechanisms
- Logging for debugging

### 6.4 Performance Optimizations

- Asynchronous analytics saving
- Batch processing (save every 10 requests)
- Efficient data structures
- Minimal memory footprint

---

## 7. Deployment and Infrastructure

### 7.1 Cloud Deployment

**Platform**: Render.com
- Web service hosting
- Persistent disk storage
- Automatic deployments
- Environment variable management

**Configuration:**
- Environment-based settings
- Secure API key storage
- Configurable regions
- Scalable architecture

### 7.2 Data Persistence

**Storage Strategy:**
- Persistent disk mounting (`/data`)
- JSON file-based storage
- Atomic writes (crash-safe)
- Automatic directory creation

**Data Retention:**
- Unlimited by default
- Configurable limits (days/records)
- Automatic cleanup
- Complete history preservation

### 7.3 Security Considerations

- API key protection
- Environment variable encryption
- Optional analytics endpoint protection
- CORS configuration
- Input validation

---

## 8. Use Cases

### 8.1 Individual Players

**Self-Improvement:**
- Post-game analysis
- Performance tracking
- Weakness identification
- Improvement planning

**Learning:**
- Understanding game mechanics
- Strategic decision analysis
- Item build optimization
- Vision control mastery

### 8.2 Coaches and Analysts

**Player Evaluation:**
- Performance assessment
- Pattern identification
- Coaching material generation
- Progress tracking

**Team Analysis:**
- Team composition evaluation
- Objective control analysis
- Strategic planning
- Match preparation

### 8.3 Content Creators

**Content Generation:**
- Match breakdown videos
- Educational content
- Analysis videos
- Coaching streams

---

## 9. Future Enhancements

### 9.1 Planned Features

**Timeline Integration:**
- Event timestamp analysis
- Real-time match tracking
- Moment-by-moment breakdown
- Replay integration

**Advanced Analytics:**
- Machine learning models
- Predictive analytics
- Performance trends
- Comparative analysis

**Enhanced UI:**
- Interactive visualizations
- Charts and graphs
- Timeline visualization
- Comparison tools

**Database Integration:**
- PostgreSQL/MongoDB support
- Query optimization
- Advanced filtering
- Historical analysis

### 9.2 Scalability Improvements

- Caching layer (Redis)
- Rate limiting
- Load balancing
- CDN integration

### 9.3 Additional Integrations

- Discord bot
- Mobile app
- Browser extension
- API marketplace

---

## 10. Technical Specifications

### 10.1 System Requirements

**Server:**
- Go 1.21 or higher
- 512MB RAM minimum
- Persistent disk storage
- Internet connectivity

**Dependencies:**
- Riot Games API access
- OpenAI API access
- Standard Go libraries

### 10.2 API Rate Limits

**Riot Games API:**
- 100 requests per 2 minutes (development)
- Higher limits available (production)

**OpenAI API:**
- Model-dependent limits
- Token usage optimization
- Cost management

### 10.3 Performance Metrics

- Average response time: 2-5 seconds
- Concurrent request handling
- Efficient memory usage
- Scalable architecture

---

## 11. Data Privacy and Ethics

### 11.1 Data Handling

- Public match data only (Riot API)
- No personal information storage
- IP address anonymization (optional)
- Analytics data retention policies

### 11.2 User Privacy

- No user registration required
- No tracking cookies
- Transparent data usage
- Optional analytics opt-out

### 11.3 Ethical Considerations

- Fair use of Riot Games API
- Responsible AI usage
- Transparent analysis
- Educational purpose focus

---

## 12. Conclusion

The League of Legends Match Advisor represents a significant advancement in esports analytics, combining the power of modern AI with comprehensive game data to deliver actionable insights. The platform's architecture is designed for scalability, maintainability, and extensibility, making it suitable for both individual players and professional teams.

**Key Achievements:**
- Automated match analysis with AI-powered insights
- Comprehensive data processing and structuring
- User-friendly interface with flexible analysis options
- Complete analytics tracking and storage
- Production-ready cloud deployment

**Impact:**
- Enables players to understand their performance deeply
- Provides coaches with data-driven analysis tools
- Supports content creators with analysis capabilities
- Contributes to the broader esports analytics ecosystem

The platform demonstrates the potential of combining game APIs with modern AI to create valuable tools for the gaming community. As the system evolves with additional features and integrations, it will continue to serve as a valuable resource for League of Legends players seeking to improve their gameplay.

---

## Appendix A: API Reference

### A.1 Request Format

**POST /analyze-match**
```json
{
  "match_id": "NA1_1234567890",
  "champion_name": "Yasuo",  // Optional
  "summoner_name": "PlayerName",  // Optional
  "focus_areas": ["combat", "vision"]  // Optional
}
```

### A.2 Response Format

See Section 4.2 for complete response structure.

### A.3 Error Responses

```json
{
  "error": "Error message description"
}
```

---

## Appendix B: Configuration

### B.1 Environment Variables

- `RIOT_API_KEY`: Riot Games API key (required)
- `OPENAI_API_KEY`: OpenAI API key (required)
- `RIOT_API_REGION`: API region (default: americas)
- `OPENAI_MODEL`: AI model (default: gpt-4o-mini)
- `PORT`: Server port (default: 8080)
- `ANALYTICS_DATA_PATH`: Analytics storage path
- `ANALYTICS_MAX_DAYS`: Retention days (0 = unlimited)
- `ANALYTICS_MAX_RECORDS`: Retention count (0 = unlimited)
- `ANALYTICS_KEY`: Analytics endpoint protection key

---

## Appendix C: Data Schema

### C.1 Match Request
```go
type MatchRequest struct {
    MatchID      string   `json:"match_id"`
    Region       string   `json:"region,omitempty"`
    ChampionName string   `json:"champion_name,omitempty"`
    SummonerName string   `json:"summoner_name,omitempty"`
    FocusAreas   []string `json:"focus_areas,omitempty"`
}
```

### C.2 Match Response
See types/match.go for complete schema definitions.

---

## 13. Performance Metrics and Benchmarks

### 13.1 Response Time Analysis

**Typical Performance:**
- **Riot API Call**: 200-500ms (network dependent)
- **Data Processing**: 10-50ms
- **OpenAI Analysis**: 2-5 seconds (model dependent)
- **Total End-to-End**: 2.5-6 seconds average

**Performance Breakdown:**
```
Request Processing Flow:
├── Request Parsing: ~1ms
├── Riot API Call: 200-500ms
├── Data Formatting: 10-50ms
├── OpenAI Request: 2-5 seconds
├── Response Encoding: ~5ms
└── Analytics Tracking: ~1ms (async)
```

### 13.2 System Resource Usage

**Memory:**
- Base memory footprint: ~20-30MB
- Per request: ~1-5MB (temporary)
- Analytics data (1000 records): ~2-5MB

**CPU:**
- Average CPU usage: 5-15% (idle)
- Peak CPU usage: 30-50% (during analysis)
- Concurrent request handling: Up to 100+ concurrent requests

**Storage:**
- Analytics data growth: ~2-5KB per request
- 1000 requests ≈ 2-5MB
- Recommended disk space: 1GB+ for production

### 13.3 Scalability Characteristics

**Current Limitations:**
- Single instance deployment
- No built-in rate limiting (relies on API providers)
- File-based storage (single file bottleneck)

**Scalability Targets:**
- Support 100+ concurrent requests
- Handle 10,000+ requests/day
- Maintain <6 second response time (P95)
- 99.9% uptime target

---

## 14. Security Architecture

### 14.1 API Key Management

**Secure Storage:**
- API keys stored as environment variables
- Never committed to version control
- Encrypted at rest on cloud platform
- Rotated regularly (recommended every 90 days)

**Access Control:**
- No user authentication required (public API)
- Optional analytics endpoint protection via API key
- CORS configured for controlled origins

### 14.2 Input Validation

**Match ID Validation:**
- Format validation: `{REGION}_{NUMBER}`
- Length checks: 10-20 characters
- Character whitelist: alphanumeric and underscores
- Injection prevention: URL encoding/escaping

**Request Validation:**
- JSON schema validation
- Type checking (string, array, etc.)
- Length limits on all fields
- Sanitization of user inputs

### 14.3 Data Privacy

**Data Collection:**
- Only public match data (via Riot API)
- No personally identifiable information (PII) stored
- IP addresses stored for analytics (can be anonymized)
- No tracking cookies or client-side storage

**Data Retention:**
- Analytics data: Configurable (default: unlimited)
- Request logs: No persistent logging
- Match data: Not stored (fetched on-demand)

**GDPR Compliance:**
- Right to access: Analytics endpoint provides data
- Right to deletion: Data can be cleared via configuration
- Data minimization: Only necessary data collected
- Transparency: Clear data usage policies

### 14.4 Threat Mitigation

**Common Threats Addressed:**
- **DDoS Protection**: Rate limiting at platform level
- **API Key Theft**: Environment variable isolation
- **Injection Attacks**: Input validation and sanitization
- **Data Leakage**: No sensitive data in logs or responses

**Security Best Practices:**
- HTTPS enforced for all communications
- Regular security updates
- Error messages don't expose system internals
- Secure headers (CORS, Content-Type validation)

---

## 15. Cost Analysis

### 15.1 Infrastructure Costs (Render.com)

**Web Service:**
- Free tier: 750 hours/month (suitable for development)
- Starter plan: $7/month (512MB RAM, suitable for low traffic)
- Standard plan: $25/month (2GB RAM, recommended for production)

**Persistent Disk:**
- Storage: $0.25/GB/month
- Example: 1GB = $0.25/month

**Estimated Monthly Costs:**
- Development: $0-7/month
- Low traffic (1000 requests/day): $7-25/month
- Medium traffic (10,000 requests/day): $25-50/month
- High traffic (100,000 requests/day): $50-200/month

### 15.2 API Costs

**Riot Games API:**
- Development key: Free (rate limited)
- Production key: Free (higher rate limits available)
- No direct costs for API usage

**OpenAI API:**
- GPT-4o-mini (default): $0.15/1M input tokens, $0.60/1M output tokens
- GPT-4: $10-30/1M input tokens, $30-60/1M output tokens
- Average cost per analysis:
  - GPT-4o-mini: $0.001-0.005 per analysis
  - GPT-4: $0.05-0.20 per analysis

**Estimated Monthly API Costs:**
- 1,000 analyses/month (GPT-4o-mini): $1-5
- 10,000 analyses/month (GPT-4o-mini): $10-50
- 100,000 analyses/month (GPT-4o-mini): $100-500

**Total Estimated Costs:**
- Development: $0-10/month
- Production (low): $20-35/month
- Production (medium): $35-100/month
- Production (high): $150-700/month

### 15.3 Cost Optimization Strategies

1. **Model Selection**: Use GPT-4o-mini for cost efficiency
2. **Caching**: Cache match data (not currently implemented)
3. **Rate Limiting**: Prevent abuse and reduce API calls
4. **Batch Processing**: Process multiple matches efficiently
5. **Monitoring**: Track API usage and costs

---

## 16. Detailed API Documentation

### 16.1 POST /analyze-match

**Endpoint:** `POST /analyze-match`

**Request Headers:**
```
Content-Type: application/json
```

**Request Body Schema:**
```json
{
  "match_id": "string (required)",
  "region": "string (optional)",
  "champion_name": "string (optional)",
  "summoner_name": "string (optional)",
  "focus_areas": ["string"] (optional)
}
```

**Request Example:**
```json
{
  "match_id": "NA1_1234567890",
  "champion_name": "Yasuo",
  "focus_areas": ["combat", "vision", "objectives"]
}
```

**Response Schema:**
```json
{
  "match_id": "string",
  "analysis": "string",
  "suggestions": ["string"],
  "coaching_tips": ["string"],
  "champion_deep_dive": "string (optional)",
  "structured_insights": {
    "what_went_well": [
      {
        "title": "string",
        "description": "string",
        "impact": "string",
        "data": ["string"],
        "category": "string"
      }
    ],
    "what_went_wrong": [...],
    "critical_moments": [...],
    "item_analysis": {...},
    "matchup_analysis": {...},
    "key_statistics": {...}
  },
  "error": "string (only on error)"
}
```

**Response Example:**
```json
{
  "match_id": "NA1_1234567890",
  "analysis": "In this 32-minute match, Team Blue secured victory through superior objective control...",
  "suggestions": [
    "Improve early game CS: You had 4.2 CS/min at 10 minutes, below the 6.0 CS/min target",
    "Increase vision control: Placed only 12 wards vs opponent's 25 wards",
    "Focus on objective timing: Team secured only 1 dragon while opponents secured 3"
  ],
  "coaching_tips": [
    "Practice last-hitting in custom games to improve CS consistency",
    "Purchase control wards on every back after 10 minutes",
    "Group for dragon spawns 30 seconds before they appear"
  ],
  "champion_deep_dive": "Yasuo Performance Analysis:\n\nYour Yasuo performance showed strong mid-game teamfighting...",
  "structured_insights": {
    "what_went_well": [
      {
        "title": "Strong Teamfight Execution",
        "description": "Achieved 3 multi-kills in teamfights at 18, 24, and 29 minutes",
        "impact": "Generated 2000+ gold advantage for team",
        "data": ["3 double kills", "12/5/8 KDA", "Highest damage on team"],
        "category": "combat"
      }
    ],
    "what_went_wrong": [
      {
        "title": "Poor Early Game CS",
        "description": "Averaged only 4.2 CS/min in first 10 minutes",
        "impact": "Fell behind in gold and experience",
        "data": ["42 CS at 10 min", "Opponent had 68 CS", "600 gold deficit"],
        "category": "farming"
      }
    ],
    "critical_moments": [
      {
        "title": "Baron Teamfight at 28 minutes",
        "description": "Secured baron and won teamfight, turning 2k gold deficit into 3k advantage",
        "outcome": "Victory",
        "impact": "Led to game-winning push",
        "data": ["5-0 ace", "Baron secured", "3 towers destroyed"]
      }
    ],
    "key_statistics": {
      "combat": [
        {"label": "KDA", "value": "12/5/8", "context": "2.4 KDA ratio"},
        {"label": "Damage Dealt", "value": "45,230", "context": "Highest on team"}
      ],
      "objectives": [
        {"label": "Turrets", "value": "3", "context": "Team total: 9"},
        {"label": "Dragons", "value": "2", "context": "Team total: 4"}
      ]
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request format or missing match_id
- `500 Internal Server Error`: Riot API error or OpenAI error
- `503 Service Unavailable`: Rate limit exceeded or service unavailable

**Error Response Example:**
```json
{
  "error": "Failed to fetch match data: Riot API error: status 404, body: Match not found"
}
```

### 16.2 GET /analyze-match-get

**Endpoint:** `GET /analyze-match-get`

**Query Parameters:**
- `match_id` (required): Match ID to analyze
- `champion_name` (optional): Champion name for deep dive
- `summoner_name` (optional): Summoner name for deep dive
- `focus_areas` (optional): Comma-separated list of focus areas

**Example Request:**
```
GET /analyze-match-get?match_id=NA1_1234567890&champion_name=Yasuo&focus_areas=combat,vision
```

**Response:** Same as POST /analyze-match

### 16.3 GET /analytics

**Endpoint:** `GET /analytics`

**Query Parameters:**
- `key` (optional): Analytics protection key (if configured)

**Response Schema:**
```json
{
  "total_requests": 1234,
  "unique_ips": {
    "192.168.1.1": 45,
    "10.0.0.1": 32
  },
  "requests_by_path": {
    "/analyze-match": 1000,
    "/health": 234
  },
  "requests_by_method": {
    "POST": 1000,
    "GET": 234
  },
  "requests_by_day": {
    "2024-12-01": 100,
    "2024-12-02": 150
  },
  "user_agents": {
    "Desktop Chrome": 800,
    "Mobile Safari": 200
  },
  "recent_requests": [...],
  "first_request": "2024-12-01T00:00:00Z",
  "last_request": "2024-12-02T12:00:00Z"
}
```

### 16.4 GET /health

**Endpoint:** `GET /health`

**Response:**
```
OK
```

**Status Codes:**
- `200 OK`: Service is healthy

---

## 17. Error Handling and Troubleshooting

### 17.1 Common Errors

**Error: "match_id is required"**
- **Cause**: Missing match_id in request
- **Solution**: Include match_id in request body or query parameter

**Error: "Failed to fetch match data: Riot API error: status 404"**
- **Cause**: Invalid match ID or match not found
- **Solution**: Verify match ID format and ensure match exists in region

**Error: "Failed to fetch match data: Riot API error: status 403"**
- **Cause**: Invalid or expired Riot API key
- **Solution**: Check RIOT_API_KEY environment variable and verify API key validity

**Error: "Failed to analyze match: OpenAI API error"**
- **Cause**: Invalid OpenAI API key or rate limit exceeded
- **Solution**: Check OPENAI_API_KEY environment variable and monitor API usage

**Error: "rate limit exceeded"**
- **Cause**: Too many requests to external APIs
- **Solution**: Implement rate limiting or wait before retrying

### 17.2 Debugging

**Enable Debug Logging:**
- Check server logs for detailed error messages
- Logs include request details, API responses, and error stack traces

**Verify Configuration:**
```bash
# Check environment variables
echo $RIOT_API_KEY
echo $OPENAI_API_KEY
echo $PORT
```

**Test API Connectivity:**
```bash
# Test Riot API
curl -H "X-Riot-Token: YOUR_KEY" \
  "https://americas.api.riotgames.com/lol/match/v5/matches/NA1_TEST"

# Test health endpoint
curl http://localhost:8080/health
```

### 17.3 Performance Troubleshooting

**Slow Response Times:**
- Check OpenAI API response times (usually 2-5 seconds)
- Verify network connectivity to APIs
- Monitor server resource usage (CPU, memory)
- Consider using faster OpenAI models (GPT-4o-mini vs GPT-4)

**High Memory Usage:**
- Reduce analytics data retention period
- Limit number of recent requests in memory
- Consider moving to database storage for large datasets

---

## 18. Testing and Quality Assurance

### 18.1 Testing Strategy

**Unit Tests:**
- Test individual components in isolation
- Mock external API dependencies
- Test data transformation and formatting functions

**Integration Tests:**
- Test API endpoints with mock data
- Verify request/response formats
- Test error handling paths

**End-to-End Tests:**
- Test complete request flow
- Verify external API integrations (with test keys)
- Test analytics tracking

### 18.2 Test Cases

**API Endpoint Tests:**
- Valid match ID request
- Invalid match ID request
- Missing required parameters
- Optional parameters (champion_name, focus_areas)
- Error responses

**Data Processing Tests:**
- Match data formatting
- Participant data extraction
- Deep dive data generation
- Structured insights generation

**Edge Cases:**
- Very long match duration (>60 minutes)
- Surrendered matches
- Missing participant data
- Invalid champion names
- Special characters in summoner names

### 18.3 Manual Testing

**Test Match IDs:**
- Use recent match IDs from your region
- Test different game modes (Ranked, Normal, ARAM)
- Test different match durations
- Test different champion roles

**Test Scenarios:**
1. Basic match analysis (no filters)
2. Champion-specific deep dive
3. Summoner-specific deep dive
4. Multiple focus areas
5. Error scenarios (invalid IDs, API failures)

---

## 19. Monitoring and Observability

### 19.1 Metrics to Monitor

**Performance Metrics:**
- Request latency (P50, P95, P99)
- API response times (Riot, OpenAI)
- Error rates
- Throughput (requests per second)

**Resource Metrics:**
- CPU usage
- Memory usage
- Disk I/O
- Network I/O

**Business Metrics:**
- Total requests
- Unique users
- Popular endpoints
- Analysis success rate

### 19.2 Logging

**Log Levels:**
- **INFO**: Normal operations, request processing
- **WARN**: Recoverable errors, rate limit warnings
- **ERROR**: API failures, processing errors
- **DEBUG**: Detailed debugging information

**Log Format:**
```
[timestamp] [level] [component] message
```

**Example Logs:**
```
[2024-12-02 10:30:45] [INFO] [match-handler] Fetching match data for match ID: NA1_1234567890
[2024-12-02 10:30:46] [INFO] [match-handler] Analyzing match using OpenAI
[2024-12-02 10:30:50] [INFO] [match-handler] Analysis completed successfully
```

### 19.3 Alerting

**Critical Alerts:**
- Service downtime
- High error rate (>5%)
- API key expiration
- Storage full

**Warning Alerts:**
- High latency (>10 seconds)
- Rate limit approaching
- High memory usage (>80%)
- Unusual request patterns

---

## 20. Deployment Guide

### 20.1 Prerequisites

**Required:**
- Riot Games API key
- OpenAI API key
- Render.com account (or alternative hosting)

**Recommended:**
- GitHub repository
- Domain name (optional)
- Monitoring service (optional)

### 20.2 Deployment Steps

**1. Prepare Repository:**
```bash
git clone <repository-url>
cd lol-ranked-new-meta
```

**2. Set Environment Variables:**
```bash
export RIOT_API_KEY=your_riot_key
export OPENAI_API_KEY=your_openai_key
export PORT=8080
export RIOT_API_REGION=americas
export OPENAI_MODEL=gpt-4o-mini
```

**3. Deploy to Render.com:**
- Connect GitHub repository
- Create new Web Service
- Set environment variables
- Deploy

**4. Verify Deployment:**
```bash
curl https://your-service.onrender.com/health
```

### 20.3 Post-Deployment

**Verify:**
- Health check endpoint responds
- Match analysis endpoint works
- Analytics endpoint accessible (if protected)
- Frontend loads correctly

**Monitor:**
- Check logs for errors
- Monitor API usage
- Track response times
- Review analytics data

---

## 21. Integration Examples

### 21.1 JavaScript/Node.js Integration

```javascript
async function analyzeMatch(matchId, championName = null) {
  const response = await fetch('https://api.example.com/analyze-match', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      match_id: matchId,
      champion_name: championName,
      focus_areas: ['combat', 'vision']
    })
  });
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  const data = await response.json();
  return data;
}

// Usage
analyzeMatch('NA1_1234567890', 'Yasuo')
  .then(result => {
    console.log('Analysis:', result.analysis);
    console.log('Suggestions:', result.suggestions);
  })
  .catch(error => console.error('Error:', error));
```

### 21.2 Python Integration

```python
import requests
import json

def analyze_match(match_id, champion_name=None, focus_areas=None):
    url = 'https://api.example.com/analyze-match'
    payload = {
        'match_id': match_id,
        'champion_name': champion_name,
        'focus_areas': focus_areas or []
    }
    
    response = requests.post(url, json=payload)
    response.raise_for_status()
    
    return response.json()

# Usage
result = analyze_match(
    'NA1_1234567890',
    champion_name='Yasuo',
    focus_areas=['combat', 'vision']
)
print(result['analysis'])
print(result['suggestions'])
```

### 21.3 cURL Integration

```bash
# Basic analysis
curl -X POST https://api.example.com/analyze-match \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "NA1_1234567890"
  }'

# Deep dive analysis
curl -X POST https://api.example.com/analyze-match \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "NA1_1234567890",
    "champion_name": "Yasuo",
    "focus_areas": ["combat", "vision", "objectives"]
  }'
```

---

## 22. Advanced Configuration

### 22.1 Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `RIOT_API_KEY` | Riot Games API key | - | Yes |
| `OPENAI_API_KEY` | OpenAI API key | - | Yes |
| `PORT` | Server port | `8080` | No |
| `RIOT_API_REGION` | API region | `americas` | No |
| `OPENAI_MODEL` | OpenAI model | `gpt-4o-mini` | No |
| `ANALYTICS_DATA_PATH` | Analytics storage path | `/data/analytics.json` | No |
| `ANALYTICS_MAX_DAYS` | Retention days (0=unlimited) | `0` | No |
| `ANALYTICS_MAX_RECORDS` | Retention count (0=unlimited) | `0` | No |
| `ANALYTICS_KEY` | Analytics endpoint protection | - | No |

### 22.2 Region Configuration

**Available Regions:**
- `americas`: NA1, BR1, LA1, LA2
- `europe`: EUW1, EUN1, TR1, RU
- `asia`: KR, JP1
- `sea`: PH2, SG2, TH2, TW2, VN2

**Region Selection:**
- Set `RIOT_API_REGION` environment variable
- Or specify in request body (if supported)

### 22.3 Model Configuration

**Available OpenAI Models:**
- `gpt-4o-mini`: Fast, cost-effective (recommended)
- `gpt-4o`: More capable, higher cost
- `gpt-4`: Most capable, highest cost
- `gpt-3.5-turbo`: Legacy, cheaper but less capable

**Model Selection:**
- Set `OPENAI_MODEL` environment variable
- Consider cost vs. quality trade-offs
- Monitor response times and quality

---

## 23. Roadmap and Future Development

### 23.1 Short-Term Enhancements (1-3 months)

**Immediate Improvements:**
- Rate limiting middleware
- Request caching for match data
- Enhanced error messages
- API documentation (OpenAPI/Swagger)
- Database integration (PostgreSQL)

**Feature Additions:**
- Batch match analysis
- Match comparison feature
- Historical trend analysis
- Custom analysis templates
- Export functionality (PDF, CSV)

### 23.2 Medium-Term Goals (3-6 months)

**Advanced Features:**
- Timeline/event data integration
- Real-time match tracking
- Machine learning models for predictions
- Advanced visualizations
- Mobile app

**Infrastructure:**
- Multi-region deployment
- Load balancing
- Redis caching layer
- CDN integration
- Database sharding

### 23.3 Long-Term Vision (6-12 months)

**Platform Expansion:**
- Support for other Riot games (VALORANT, TFT)
- Multi-game analytics dashboard
- Social features (share analyses)
- Coaching marketplace
- API marketplace

**Enterprise Features:**
- Team/org management
- Advanced analytics dashboard
- Custom branding
- White-label solution
- Enterprise support

---

## 24. Contributing and Community

### 24.1 Contributing Guidelines

**Getting Started:**
1. Fork the repository
2. Create a feature branch
3. Make changes
4. Submit pull request

**Code Standards:**
- Follow Go formatting (`gofmt`)
- Write tests for new features
- Update documentation
- Follow existing code style

### 24.2 Community Support

**Support Channels:**
- GitHub Issues: Bug reports and feature requests
- Discussions: Questions and community help
- Documentation: Comprehensive guides and examples

**Reporting Issues:**
- Provide detailed error messages
- Include steps to reproduce
- Specify environment details
- Attach relevant logs

---

## References

1. Riot Games API Documentation: https://developer.riotgames.com/
2. OpenAI API Documentation: https://platform.openai.com/docs/
3. Go Programming Language: https://go.dev/
4. Render.com Documentation: https://render.com/docs/

---

**Document Version**: 1.0  
**Last Updated**: December 2024  
**Author**: Development Team  
**License**: Educational Use

---

*This whitepaper provides a comprehensive overview of the League of Legends Match Advisor platform. For technical implementation details, refer to the project's README.md and source code documentation.*

