# Frontend - League of Legends Match Advisor

A modern, interactive web frontend for the League of Legends Match Advisor API.

## Features

- 🎯 **Interactive Tabs**: Click through different analysis sections
- ✅ **What Went Well**: Highlight specific achievements and positive events
- ❌ **What Went Wrong**: Identify specific mistakes with supporting data
- ⚡ **Critical Moments**: Game-changing events with context
- 🛡️ **Item Analysis**: Detailed item build analysis with timing and matchup considerations
- ⚔️ **Matchup Analysis**: Champion matchups and team composition insights
- 🔍 **Deep Dive**: Comprehensive champion-specific analysis

## Setup

1. **Update API URL**: Edit `script.js` and update the `API_URL` constant:
   ```javascript
   const API_URL = 'https://your-service.onrender.com';
   ```

2. **Serve the files**: You can serve these files using:
   - Any web server (nginx, Apache)
   - GitHub Pages
   - Vercel / Netlify (static hosting)
   - Simple Python HTTP server: `python3 -m http.server 8000`

## Usage

1. Open `index.html` in a web browser
2. Enter a Match ID (required)
3. Optionally specify a Champion Name or Summoner Name for deep dive analysis
4. Click "Analyze Match"
5. Navigate through the tabs to view different insights

## Deployment Options

### Option 1: GitHub Pages
1. Push the `frontend/` folder to a GitHub repository
2. Enable GitHub Pages in repository settings
3. Select the branch/folder containing the frontend files

### Option 2: Vercel/Netlify
1. Connect your GitHub repository
2. Set build directory to `frontend/`
3. Deploy (no build step needed for static files)

### Option 3: Render Static Site
1. Create a new Static Site in Render
2. Point to your repository
3. Set build command: `echo "No build needed"`
4. Set publish directory: `frontend`

## Customization

- **Colors**: Edit CSS variables in `style.css` (`:root` section)
- **API URL**: Update `API_URL` in `script.js`
- **Styling**: Modify `style.css` to match your branding

