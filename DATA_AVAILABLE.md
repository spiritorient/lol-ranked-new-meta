# Exact Data Available to the App

This document explains **exactly** what data is available to your app from the Riot API and what you can use to configure and create the final output.

## 1. Raw Data Received from Riot API

Your app receives a `RiotMatch` object from the Riot Games API v5 endpoint. This contains:

### Match Metadata (`RiotMatchMetadata`)
- `dataVersion`: String - API data version
- `matchId`: String - The match identifier
- `participants`: Array of strings - List of participant PUUIDs

### Match Information (`RiotMatchInfo`)
- `gameCreation`: int64 - Timestamp when game was created
- `gameDuration`: int64 - **Game duration in seconds**
- `gameEndTimestamp`: int64 - Timestamp when game ended
- `gameId`: int64 - Internal game ID
- `gameMode`: String - Game mode (e.g., "CLASSIC", "ARAM")
- `gameName`: String - Game name
- `gameStartTimestamp`: int64 - Timestamp when game started
- `gameType`: String - Game type
- `gameVersion`: String - **Game version/patch number**
- `mapId`: int - Map identifier
- `participants`: Array - **All 10 players' data** (see below)
- `platformId`: String - Platform identifier
- `queueId`: int - Queue type ID
- `teams`: Array - **Team-level data** (see below)
- `tournamentCode`: String - Tournament code (if applicable)

### Participant Data (`RiotParticipant`) - Available for ALL 10 Players

Each participant object contains:

#### Basic Identity
- `summonerName`: String - **Player's summoner name**
- `riotIdName`: String - Riot ID name
- `riotIdTagline`: String - Riot ID tagline
- `puuid`: String - Player's unique identifier
- `summonerId`: String - Summoner ID
- `summonerLevel`: int - Summoner level
- `profileIcon`: int - Profile icon ID

#### Champion Information
- `championId`: int - Champion ID
- `championName`: String - **Champion name** (e.g., "Jinx", "Yasuo")
- `championTransform`: int - Transform state (for Kayn, etc.)
- `champLevel`: int - **Final champion level reached**

#### Team & Position
- `teamId`: int - **100 (Blue) or 200 (Red)**
- `teamPosition`: String - **Position** (e.g., "TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY")
- `individualPosition`: String - Individual position
- `lane`: String - Lane played
- `role`: String - Role played

#### Match Result
- `win`: bool - **Whether this player won or lost**

#### Combat Statistics
- `kills`: int - **Total kills**
- `deaths`: int - **Total deaths**
- `assists`: int - **Total assists**
- `largestKillingSpree`: int - Largest killing spree
- `killingSprees`: int - Number of killing sprees
- `largestMultiKill`: int - Largest multi-kill
- `doubleKills`: int - Double kills
- `tripleKills`: int - Triple kills
- `quadraKills`: int - Quadra kills
- `pentaKills`: int - Penta kills
- `unrealKills`: int - Unreal kills
- `largestCriticalStrike`: int - Largest crit damage

#### Damage Statistics
- `totalDamageDealt`: int - Total damage dealt
- `totalDamageDealtToChampions`: int - **Total damage to champions**
- `physicalDamageDealt`: int - Physical damage dealt
- `physicalDamageDealtToChampions`: int - **Physical damage to champions**
- `magicDamageDealt`: int - Magic damage dealt
- `magicDamageDealtToChampions`: int - **Magic damage to champions**
- `trueDamageDealt`: int - True damage dealt
- `trueDamageDealtToChampions`: int - **True damage to champions**
- `totalDamageTaken`: int - **Total damage taken**
- `physicalDamageTaken`: int - Physical damage taken
- `magicDamageTaken`: int - Magic damage taken
- `trueDamageTaken`: int - True damage taken
- `damageSelfMitigated`: int - Damage self-mitigated
- `totalHeal`: int - Total healing done
- `totalHealsOnTeammates`: int - Healing on teammates
- `totalDamageShieldedOnTeammates`: int - **Shielding on teammates**

#### Economy Statistics
- `goldEarned`: int - **Total gold earned**
- `goldSpent`: int - **Total gold spent**
- `itemsPurchased`: int - **Number of items purchased**

#### Farming Statistics
- `totalMinionsKilled`: int - **Total CS (minions killed)**
- `neutralMinionsKilled`: int - **Jungle camps killed**

#### Objective Control
- `turretKills`: int - **Turrets destroyed**
- `turretTakedowns`: int - Turret takedowns
- `turretsLost`: int - Turrets lost
- `inhibitorKills`: int - **Inhibitors destroyed**
- `inhibitorTakedowns`: int - Inhibitor takedowns
- `inhibitorsLost`: int - Inhibitors lost
- `dragonKills`: int - **Dragons killed**
- `baronKills`: int - **Barons killed**
- `nexusKills`: int - Nexus kills
- `nexusTakedowns`: int - Nexus takedowns
- `nexusLost`: int - Nexus lost
- `objectivesStolen`: int - Objectives stolen
- `objectivesStolenAssists`: int - Objective steal assists
- `damageDealtToTurrets`: int - Damage to turrets
- `damageDealtToBuildings`: int - Damage to buildings
- `damageDealtToObjectives`: int - Damage to objectives

#### First Blood/Tower
- `firstBloodKill`: bool - **Got first blood kill**
- `firstBloodAssist`: bool - Got first blood assist
- `firstTowerKill`: bool - **Got first tower kill**
- `firstTowerAssist`: bool - Got first tower assist

#### Vision Statistics
- `visionScore`: int - **Vision score**
- `wardsPlaced`: int - **Wards placed**
- `wardsKilled`: int - **Wards killed**
- `visionWardsBoughtInGame`: int - **Control wards purchased**
- `detectorWardsPlaced`: int - **Control wards placed**
- `sightWardsBoughtInGame`: int - Sight wards purchased

#### Items
- `item0`: int - **Item slot 1 ID**
- `item1`: int - **Item slot 2 ID**
- `item2`: int - **Item slot 3 ID**
- `item3`: int - **Item slot 4 ID**
- `item4`: int - **Item slot 5 ID**
- `item5`: int - **Item slot 6 ID**
- `item6`: int - **Trinket item ID**
- `consumablesPurchased`: int - Consumables purchased

#### Summoner Spells
- `summoner1Id`: int - **Summoner spell 1 ID**
- `summoner1Casts`: int - **Times spell 1 was cast**
- `summoner2Id`: int - **Summoner spell 2 ID**
- `summoner2Casts`: int - **Times spell 2 was cast**

#### Ability Casts
- `spell1Casts`: int - Q ability casts
- `spell2Casts`: int - W ability casts
- `spell3Casts`: int - E ability casts
- `spell4Casts`: int - R ability casts

#### Crowd Control
- `timeCCingOthers`: int - **Time spent CC'ing others (in seconds)**
- `totalTimeCCDealt`: int - **Total CC time dealt**

#### Time Statistics
- `timePlayed`: int - Time played
- `totalTimeSpentDead`: int - **Total time spent dead (seconds)**
- `longestTimeSpentLiving`: int - **Longest time alive (seconds)**

#### Bounty & Challenges
- `bountyLevel`: int - Bounty level
- `challenges`: map[string]interface{} - **Various challenge stats** (complex nested data)

#### Surrender Information
- `gameEndedInEarlySurrender`: bool - Early surrender
- `gameEndedInSurrender`: bool - Surrender
- `teamEarlySurrendered`: bool - Team early surrendered

### Team Data (`RiotTeam`)

For each team (Blue/Red):

- `teamId`: int - **100 (Blue) or 200 (Red)**
- `win`: bool - **Whether team won**
- `bans`: Array of `RiotBan` - **Champions banned**
  - `championId`: int - Banned champion ID
  - `pickTurn`: int - When it was banned
- `objectives`: `RiotObjectives` - **Team objective control**
  - `baron`: `RiotObjective` - Baron stats
    - `first`: bool - **Got first baron**
    - `kills`: int - **Total barons killed**
  - `champion`: `RiotObjective` - Champion objective
  - `dragon`: `RiotObjective` - Dragon stats
    - `first`: bool - **Got first dragon**
    - `kills`: int - **Total dragons killed**
  - `inhibitor`: `RiotObjective` - Inhibitor stats
    - `first`: bool - **Got first inhibitor**
    - `kills`: int - **Total inhibitors destroyed**
  - `riftHerald`: `RiotObjective` - Rift Herald stats
    - `first`: bool - **Got first rift herald**
    - `kills`: int - **Total rift heralds killed**
  - `tower`: `RiotObjective` - Tower stats
    - `first`: bool - **Got first tower**
    - `kills`: int - **Total towers destroyed**

### Runes/Perks (`RiotPerks`)
- `statPerks`: `RiotStatPerks` - Stat runes
  - `defense`: int
  - `flex`: int
  - `offense`: int
- `styles`: Array of `RiotPerkStyle` - Rune trees
  - `description`: String
  - `style`: int - Rune tree ID
  - `selections`: Array of `RiotPerkSelection` - Specific runes
    - `perk`: int - Rune ID
    - `var1`, `var2`, `var3`: int - Rune variables

---

## 2. Data That is NOT Available

**Important limitations:**

1. **No Timeline Data**: The app does NOT receive:
   - Exact timestamps of events (kills, deaths, objectives)
   - Item purchase timestamps
   - Gold/XP over time
   - Position data over time
   - Exact moment-by-moment game state

2. **No Real-Time Events**: You only get:
   - Final statistics
   - Aggregate totals
   - Final item builds (not purchase order)

3. **No Match History**: The app only processes ONE match at a time

---

## 3. How Data is Processed

### Step 1: Raw Data → Formatted Summary
The `FormatMatchForAnalysis()` function in `riot/client.go` converts the raw `RiotMatch` into a text summary that includes:

- Match metadata (ID, duration, mode, version)
- Team summaries (which team won, objectives secured)
- All 10 participants with basic stats (K/D/A, CS, Gold, Damage)
- **If champion/summoner filter is provided**: Detailed deep dive including:
  - All combat stats
  - Objective control
  - Vision stats
  - Item build (final items only, no timestamps)
  - Opponent composition analysis
  - Lane matchup comparison

### Step 2: Summary → OpenAI Analysis
The formatted summary is sent to OpenAI with prompts requesting:
- Data-driven analysis
- Specific event identification
- Comparison with opponents
- Actionable suggestions

### Step 3: OpenAI → Structured Response
OpenAI returns:
- `analysis`: String - General match analysis
- `suggestions`: Array of strings - Improvement suggestions
- `coaching_tips`: Array of strings - Coaching tips
- `champion_deep_dive`: String (optional) - Deep dive for specific champion/summoner
- `structured_insights`: `StructuredInsights` (optional) - Structured data with:
  - `what_went_well`: Array of specific events with data
  - `what_went_wrong`: Array of specific events with data
  - `critical_moments`: Array of key moments
  - `item_analysis`: Item build analysis
  - `matchup_analysis`: Champion/team composition analysis
  - `key_statistics`: Organized stats by category

---

## 4. Final Output Structure

Your app returns a `MatchResponse` JSON with:

```json
{
  "match_id": "string",
  "analysis": "string - comprehensive analysis text",
  "suggestions": ["string", "string", ...],
  "coaching_tips": ["string", "string", ...],
  "champion_deep_dive": "string (optional)",
  "structured_insights": {
    "what_went_well": [
      {
        "title": "string",
        "description": "string",
        "impact": "string",
        "data": ["string", ...],
        "category": "string"
      }
    ],
    "what_went_wrong": [...],
    "critical_moments": [...],
    "item_analysis": {
      "timing_analysis": "string",
      "opponent_matchup": "string",
      "recommendations": ["string", ...]
    },
    "matchup_analysis": {
      "lane_matchup": "string",
      "team_composition": "string",
      "synergies": ["string", ...],
      "counters": ["string", ...],
      "win_conditions": ["string", ...]
    },
    "key_statistics": {
      "combat": [{"label": "string", "value": "string", "context": "string"}],
      "objectives": [...],
      "economy": [...],
      "vision": [...]
    }
  },
  "error": "string (if error occurred)"
}
```

---

## 5. Key Data Points You Can Use

### For General Analysis:
- All 10 players' K/D/A, CS, Gold, Damage
- Team objective control (dragons, barons, towers)
- Game duration and version
- Win/loss for each player

### For Deep Dive Analysis (when champion/summoner filter provided):
- **Exact lane matchup**: Opponent champion, position, stats comparison
- **Detailed combat stats**: All damage types, damage taken, healing, shielding
- **Vision control**: Vision score, wards placed/killed, control wards
- **Objective participation**: Dragons, barons, towers, inhibitors
- **Item build**: Final 6 items + trinket (item IDs)
- **Summoner spells**: Which spells and how many times used
- **Time management**: Time spent dead, longest time alive
- **Special achievements**: Multi-kills, killing sprees, first blood/tower

### For Team Composition Analysis:
- All champions on both teams
- Team positions
- Bans
- Team objective control

---

## 6. Data Limitations for Analysis

Since you don't have timeline data, you **cannot**:
- Know exactly when events happened (only that they happened)
- Analyze item purchase timing (only final build)
- Track gold/XP curves over time
- Identify specific game phases (early/mid/late) with precision
- Analyze positioning or movement patterns

You **can**:
- Compare final stats between players
- Analyze final item builds vs opponent composition
- Calculate rates (CS/min, gold/min) using game duration
- Identify standout performances using aggregate stats
- Compare vision control between players
- Analyze objective control patterns

---

## Summary

**You have access to:**
- Complete end-of-game statistics for all 10 players
- Team-level objective control
- Final item builds
- Champion selections and positions
- All combat, farming, vision, and objective statistics
- Game metadata (duration, version, mode)

**You do NOT have:**
- Timeline/event timestamps
- Item purchase order/timing
- Real-time game state
- Position data over time

The app processes this data through OpenAI to generate coaching insights, but all analysis is based on these aggregate statistics, not minute-by-minute game events.

