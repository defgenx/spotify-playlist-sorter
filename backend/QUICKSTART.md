# Quick Start Guide

## Prerequisites

1. **Spotify Developer Account**: Get your credentials at [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. **Go 1.24+**: Install from [golang.org](https://golang.org/)

## Setup (5 minutes)

### 1. Create Spotify App

1. Go to https://developer.spotify.com/dashboard
2. Click "Create App"
3. Fill in:
   - App name: "Spotify Playlist Sorter"
   - App description: "Organize liked songs by genre"
   - Redirect URI: `http://localhost:8080/api/auth/callback`
4. Click "Save"
5. Copy your **Client ID** and **Client Secret**

### 2. Configure Environment

```bash
cd backend
cp .env.example .env
```

Edit `.env` and add your Spotify credentials:

```env
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here
SESSION_SECRET=$(openssl rand -base64 32)
```

### 3. Run the Server

```bash
# Install dependencies (first time only)
go mod download

# Run the server
go run cmd/server/main.go
```

You should see:
```
INFO Starting Spotify Playlist Sorter API
INFO Configuration loaded port=8080
INFO Starting HTTP server addr=:8080
```

## Testing the API

### 1. Check Health

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

### 2. Test Authentication Flow

Open your browser to:
```
http://localhost:8080/api/auth/login
```

This will return a JSON with the Spotify login URL:
```json
{
  "url": "https://accounts.spotify.com/authorize?..."
}
```

### 3. Complete OAuth Flow

Visit the URL from step 2 in your browser. After authorizing, you'll be redirected to the frontend URL (which may not exist yet, but the session will be created).

## API Endpoints

### Authentication
- `GET /api/auth/login` - Get Spotify OAuth URL
- `GET /api/auth/callback` - OAuth callback (automatic)
- `GET /api/auth/me` - Get current user profile
- `POST /api/auth/logout` - Logout

### Library Analysis
- `GET /api/library/analysis` - Analyze your music library

### Sorting
- `POST /api/sort/plan` - Generate sort plan
  ```bash
  curl -X POST http://localhost:8080/api/sort/plan \
    -H "Content-Type: application/json" \
    -d '{"dryRun": true}' \
    --cookie "spotify_session=YOUR_SESSION_ID"
  ```

- `POST /api/sort/execute` - Execute sort plan
  ```bash
  curl -X POST http://localhost:8080/api/sort/execute \
    -H "Content-Type: application/json" \
    -d '{"dryRun": false}' \
    --cookie "spotify_session=YOUR_SESSION_ID"
  ```

### Real-time Progress
- `GET /api/events` - SSE stream for progress updates

## Typical Workflow

1. **Login** via `/api/auth/login`
2. **Analyze Library** via `/api/library/analysis`
   - Fetches all liked songs
   - Gets genre information
   - Identifies existing playlists
3. **Generate Plan** via `/api/sort/plan` (dry run)
   - Preview what will be changed
4. **Execute Plan** via `/api/sort/execute`
   - Creates genre playlists
   - Moves songs to correct playlists
   - Removes songs from wrong playlists

## What Gets Created

The sorter creates playlists like:
- "indie rock" - All indie rock songs
- "electronic" - All electronic songs
- "jazz" - All jazz songs
- "Uncategorized" - Songs without clear genre

All created playlists are marked with `[Managed by SpotifyPlaylistSorter]` in the description.

## Troubleshooting

### "Failed to load configuration"
- Check your `.env` file exists
- Verify all required variables are set

### "Failed to authenticate"
- Verify Spotify credentials are correct
- Check redirect URI matches in Spotify Dashboard

### "Rate limited"
- The app includes automatic rate limiting
- Wait a moment and try again

### "Session expired"
- Re-authenticate via `/api/auth/login`

## Development

### Run with Hot Reload
```bash
# Install air (if not installed)
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

### Run Tests
```bash
go test ./...
```

### Build for Production
```bash
go build -o spotify-sorter cmd/server/main.go
./spotify-sorter
```

## Next Steps

1. Build the frontend (React/Vue/etc) to consume this API
2. Add authentication state management
3. Create UI for viewing sort plans
4. Add SSE event handling for real-time progress
5. Deploy to production (Heroku, Railway, Fly.io, etc)

## Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8080` | Server port |
| `FRONTEND_URL` | No | `http://localhost:5173` | Frontend URL for redirects |
| `CORS_ORIGINS` | No | `http://localhost:5173` | Allowed CORS origins |
| `SPOTIFY_CLIENT_ID` | Yes | - | Your Spotify app client ID |
| `SPOTIFY_CLIENT_SECRET` | Yes | - | Your Spotify app client secret |
| `SPOTIFY_REDIRECT_URL` | No | `http://localhost:8080/api/auth/callback` | OAuth callback URL |
| `SESSION_SECRET` | Yes | - | Random secret for session encryption |

## Support

For issues or questions, check:
- The main README.md
- Spotify Web API docs: https://developer.spotify.com/documentation/web-api
- Go documentation: https://pkg.go.dev
