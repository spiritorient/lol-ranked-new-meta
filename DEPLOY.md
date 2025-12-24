# Deploying to Render.com

This guide will help you deploy the League of Legends Match Advisor backend to Render.com.

## Prerequisites

1. A GitHub account with this repository
2. A Render.com account (free tier available)
3. Riot Games API key
4. OpenAI API key

## Step-by-Step Deployment

### Option 1: Using Render Dashboard (Recommended)

1. **Log in to Render**
   - Go to https://dashboard.render.com
   - Sign up or log in with your GitHub account

2. **Create a New Web Service**
   - Click "New +" → "Web Service"
   - Connect your GitHub account if not already connected
   - Select the repository: `spiritorient/lol-ranked-new-meta`

3. **Configure the Service**
   - **Name**: `lol-ranked-new-meta` (or your preferred name)
   - **Region**: Choose closest to your users
   - **Branch**: `main`
   - **Root Directory**: Leave empty (or `.` if needed)
   - **Runtime**: `Go`
   - **Build Command**: `go mod download && go build -o server`
   - **Start Command**: `./server`
   - **Plan**: Free (or upgrade if needed)

4. **Environment Variables**
   Click "Advanced" and add these environment variables:
   
   | Key | Value | Required |
   |-----|-------|----------|
   | `RIOT_API_KEY` | Your Riot Games API key | ✅ Yes |
   | `OPENAI_API_KEY` | Your OpenAI API key | ✅ Yes |
   | `RIOT_API_REGION` | `americas` (or europe, asia, sea) | No (defaults to americas) |
   | `OPENAI_MODEL` | `gpt-4o-mini` (or gpt-4, etc.) | No (defaults to gpt-4o-mini) |
   | `PORT` | Auto-set by Render | ❌ No (auto-provided) |

5. **Deploy**
   - Click "Create Web Service"
   - Render will automatically build and deploy your service
   - Wait for the build to complete (usually 2-5 minutes)

6. **Access Your Service**
   - Once deployed, your service will be available at: `https://lol-ranked-new-meta.onrender.com`
   - (Or the custom domain you configure)
   - Test the health endpoint: `https://your-service.onrender.com/health`

### Option 2: Using render.yaml (Alternative)

If you prefer using the `render.yaml` file:

1. Follow steps 1-2 from Option 1
2. When creating the service, Render will automatically detect `render.yaml`
3. The configuration will be pre-filled from the YAML file
4. You still need to add the secret environment variables (`RIOT_API_KEY` and `OPENAI_API_KEY`) in the dashboard

## Verifying Deployment

Once deployed, test your service:

```bash
# Health check
curl https://your-service.onrender.com/health

# Test match analysis (replace with real match ID)
curl -X POST https://your-service.onrender.com/analyze-match \
  -H "Content-Type: application/json" \
  -d '{"match_id": "NA1_1234567890"}'
```

## Important Notes

### Free Tier Limitations

- **Spinning Down**: Free services spin down after 15 minutes of inactivity
- **First Request**: The first request after spin-down may take 30-60 seconds (cold start)
- **Upgrade**: Consider upgrading to a paid plan for always-on service

### Environment Variables Security

- Never commit API keys to your repository
- Always set sensitive variables in Render dashboard
- Use Render's environment variable encryption

### Logs and Monitoring

- View logs in the Render dashboard under "Logs" tab
- Monitor service health and uptime
- Set up alerts for service failures

### Custom Domain (Optional)

1. Go to your service settings
2. Click "Custom Domains"
3. Add your domain
4. Follow DNS configuration instructions

## Troubleshooting

### Build Failures

- Check build logs in Render dashboard
- Ensure `go.mod` and `go.sum` are committed
- Verify Go version compatibility (Go 1.21+)

### Runtime Errors

- Check runtime logs
- Verify all environment variables are set
- Test API keys locally first

### Service Not Starting

- Check that `PORT` environment variable is available (Render provides this)
- Verify start command is correct: `./server`
- Check health check endpoint is accessible

## Updating Your Service

After pushing changes to GitHub:

1. Render automatically detects new commits
2. Triggers a new deployment
3. You can also manually deploy from the dashboard

## Cost Estimation

- **Free Tier**: $0/month (with spin-down limitations)
- **Starter Plan**: ~$7/month (always-on, no spin-down)
- **Standard Plan**: ~$25/month (better performance)

Choose based on your usage and requirements.

