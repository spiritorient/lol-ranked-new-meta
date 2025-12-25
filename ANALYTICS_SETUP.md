# Analytics Setup Guide

## Render.com Persistent Disk Configuration

To enable analytics tracking with persistent storage on Render.com:

### Step 1: Add a Disk in Render.com Dashboard

1. Go to your Render.com dashboard
2. Select your service (`lol-ranked-new-meta`)
3. Go to **Settings** → **Disks**
4. Click **Add Disk**
5. Configure:
   - **Name**: `analytics-data` (or any name you prefer)
   - **Mount Path**: `/data` (this is the absolute path where the disk will be mounted)
   - **Size**: `1 GB` (minimum, increase if needed)

### Step 2: Set Environment Variable (Optional)

The `render.yaml` file already includes the configuration, but you can also set it manually:

- **Key**: `ANALYTICS_DATA_PATH`
- **Value**: `/data/analytics.json`

This tells the app to store analytics data on the persistent disk.

### Step 3: Deploy

After adding the disk and setting the environment variable, redeploy your service. The analytics data will now persist across deployments.

## Local Development

For local development, the app will try to use `/data/analytics.json` by default. To use a local path instead, set:

```bash
export ANALYTICS_DATA_PATH=./data/analytics.json
```

Or create a `.env` file:

```
ANALYTICS_DATA_PATH=./data/analytics.json
```

## Viewing Analytics

Once deployed, you can view analytics at:

```
https://your-app.onrender.com/analytics
```

### Response Format

The default response includes:
- Summary statistics
- Requests by path, method, day
- User agents breakdown
- Top IPs
- Recent requests (last 100)
- Sample of all stored requests (last 50)

### View All Stored Requests

To get **all stored requests** (can be large!):

```
https://your-app.onrender.com/analytics?all=true
```

**Warning:** If you have many requests stored, this can be a very large JSON response. Use with caution!

### Optional: Protect Analytics Endpoint

To protect the analytics endpoint with an API key:

1. Set environment variable in Render.com:
   - **Key**: `ANALYTICS_KEY`
   - **Value**: (your secret key)

2. Access analytics with:
   ```
   https://your-app.onrender.com/analytics?key=YOUR_SECRET_KEY
   ```

## What Gets Tracked

- **Total requests** - Count of all HTTP requests
- **Unique IPs** - Number of unique visitors
- **Requests by path** - Which endpoints are accessed most
- **Requests by method** - GET, POST, etc.
- **Requests by day** - Daily breakdown
- **User agents** - Browser/device types
- **Top IPs** - Most frequent visitors
- **Recent requests** - Last 100 requests in memory (for quick access)
- **All requests** - **ALL requests are stored on disk** (unlimited by default!)

## Storage Configuration

By default, **ALL requests are stored permanently** on the persistent disk. This means you have a complete history of every request ever made to your app.

### Optional Limits

If you want to limit storage, you can set these environment variables:

- **`ANALYTICS_MAX_DAYS`** - Maximum number of days to keep requests (e.g., `30` = keep last 30 days, `0` = unlimited)
- **`ANALYTICS_MAX_RECORDS`** - Maximum total number of records to keep (e.g., `10000` = keep last 10,000 requests, `0` = unlimited)

**Examples:**

```bash
# Keep only last 30 days of requests
ANALYTICS_MAX_DAYS=30

# Keep only last 10,000 requests
ANALYTICS_MAX_RECORDS=10000

# Keep last 30 days OR 10,000 requests (whichever limit is hit first)
ANALYTICS_MAX_DAYS=30
ANALYTICS_MAX_RECORDS=10000
```

**Note:** If you don't set these variables, all requests are stored indefinitely (unlimited).

## Data Storage

- Analytics data is stored in JSON format
- File location: `/data/analytics.json` (on persistent disk)
- Data persists across deployments and restarts
- The file is automatically created on first request
- **All requests are stored** (not just recent ones)
- Data is saved every 10 requests to balance performance and data safety
- Old records are automatically cleaned up if limits are set

## Troubleshooting

### Analytics not working?

1. **Check disk is mounted**: Verify the disk is added in Render.com dashboard
2. **Check mount path**: Ensure mount path is `/data` (or update `ANALYTICS_DATA_PATH` accordingly)
3. **Check permissions**: The app should have write permissions to the mount path
4. **Check logs**: Look for analytics-related messages in Render.com logs

### Disk size issues?

If you run out of space:
1. Go to Render.com dashboard → Your service → Disks
2. Increase the disk size
3. Redeploy if needed

### Want to reset analytics?

Simply delete the analytics file:
```bash
rm /data/analytics.json
```

The file will be recreated on the next request.

