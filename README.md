# Spotify Playlist Sorter

> **Warning**
> This project is a **work in progress** and was heavily vibe-coded with AI assistance. Expect rough edges, potential bugs, and incomplete features. Use at your own risk and don't blame the robots if your playlists get weird.

Automatically organize your Spotify liked songs into genre-based playlists.

## Features

- **OAuth Authentication**: Secure login with Spotify via HTTPS (ngrok)
- **Genre Detection**: Analyzes songs and detects genres from artist data
- **Smart Sorting**: Organizes songs into genre-based playlists
- **Dry-Run Mode**: Preview all changes before applying them
- **Real-time Progress**: SSE-powered progress updates during analysis and execution
- **Playlist Reorganization**: Moves songs out of wrong genre playlists (only app-managed ones)

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│  React Frontend │────▶│   Go Backend    │────▶│   Spotify API   │
│  (Port 3000)    │     │   (Port 3001)   │     │                 │
│                 │◀────│                 │◀────│                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │
        │    ┌─────────────┐    │
        └───▶│    ngrok    │◀───┘
             │   (HTTPS)   │
             └─────────────┘
                   │
                   ▼
            Spotify OAuth
              Callback
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- ngrok account (free tier works)
- Spotify Developer account

### 1. Clone and Setup

```bash
cd spotify-playlist-sorter

# Backend setup
cd backend
cp ../.env.example .env
go mod tidy

# Frontend setup
cd ../frontend
npm install
```

### 2. Configure ngrok

```bash
# Login to ngrok (one-time)
ngrok config add-authtoken YOUR_AUTHTOKEN

# Start tunnel
ngrok http 3001
```

Copy the HTTPS URL (e.g., `https://abc123.ngrok-free.dev`)

### 3. Configure Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app (or use existing)
3. Add Redirect URI: `https://YOUR-NGROK-URL/api/auth/callback`
4. Copy Client ID and Client Secret

### 4. Update Configuration

Edit `backend/.env`:
```bash
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
SPOTIFY_REDIRECT_URL=https://YOUR-NGROK-URL/api/auth/callback
CORS_ORIGINS=http://localhost:3000,https://YOUR-NGROK-URL
```

### 5. Start the Application

```bash
# Terminal 1: ngrok (keep running)
ngrok http 3001

# Terminal 2: Backend
cd backend
go run ./cmd/server

# Terminal 3: Frontend
cd frontend
npm run dev
```

### 6. Use the App

1. Open `http://localhost:3000`
2. Click "Login with Spotify"
3. Authorize the app
4. Click "Start Analysis"
5. Review detected genres
6. Go to "Changes" to preview
7. Toggle off "Preview Mode" and click "Execute"

## Project Structure

```
spotify-playlist-sorter/
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── internal/
│   │   ├── api/                    # HTTP handlers & routes
│   │   │   ├── handlers/           # Request handlers
│   │   │   │   ├── auth.go         # OAuth + session
│   │   │   │   ├── library.go      # Library analysis
│   │   │   │   ├── sort.go         # Sort operations
│   │   │   │   └── events.go       # SSE streaming
│   │   │   ├── middleware/         # Auth, CORS
│   │   │   └── router.go           # Route definitions
│   │   ├── config/                 # Environment config
│   │   ├── domain/                 # Domain models
│   │   │   ├── track.go
│   │   │   ├── playlist.go
│   │   │   └── sortplan.go
│   │   ├── genre/                  # Genre normalization
│   │   ├── service/                # Business logic
│   │   │   ├── library.go          # Fetch & analyze
│   │   │   ├── sorter.go           # Plan generation
│   │   │   └── executor.go         # Plan execution
│   │   ├── session/                # Session store
│   │   ├── spotify/                # Spotify SDK wrapper
│   │   └── sse/                    # Event broadcasting
│   ├── certs/                      # SSL certs (optional)
│   └── .env
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/                 # Button, Card, Badge, etc.
│   │   │   ├── layout/             # Header, Sidebar
│   │   │   ├── songs/              # SongCard, SongList
│   │   │   ├── genres/             # GenreCard, GenreGrid
│   │   │   ├── changes/            # ChangeSummary, ChangeDiff
│   │   │   └── progress/           # ProgressBar, LiveLog
│   │   ├── hooks/
│   │   │   ├── useAuth.ts
│   │   │   └── useSSE.ts
│   │   ├── stores/
│   │   │   ├── authStore.ts
│   │   │   └── uiStore.ts          # Dry-run toggle
│   │   ├── lib/
│   │   │   ├── api.ts
│   │   │   └── types.ts
│   │   └── pages/
│   │       ├── Login.tsx
│   │       ├── Callback.tsx
│   │       ├── Dashboard.tsx
│   │       ├── Genres.tsx
│   │       └── Changes.tsx
│   └── package.json
├── docs/
│   ├── API.md
│   ├── ARCHITECTURE.md
│   └── SETUP.md
└── README.md
```

## API Reference

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/auth/login` | No | Get Spotify OAuth URL |
| GET | `/api/auth/callback` | No | OAuth callback (via ngrok) |
| GET | `/api/auth/complete` | No | Complete login, set cookie |
| GET | `/api/auth/me` | Yes | Get current user profile |
| POST | `/api/auth/logout` | Yes | Logout, clear session |

### Library

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/library/analysis` | Yes | Analyze liked songs |

### Sort

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/sort/plan` | Yes | Generate sort plan |
| POST | `/api/sort/execute` | Yes | Execute sort plan |

### Events (SSE)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/events` | Yes | Real-time progress stream |

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SPOTIFY_CLIENT_ID` | Yes | - | Spotify app client ID |
| `SPOTIFY_CLIENT_SECRET` | Yes | - | Spotify app client secret |
| `SPOTIFY_REDIRECT_URL` | Yes | - | OAuth callback URL (ngrok HTTPS) |
| `PORT` | No | `3001` | Backend server port |
| `FRONTEND_URL` | No | `http://localhost:3000` | Frontend URL |
| `CORS_ORIGINS` | No | `http://localhost:3000` | Allowed CORS origins |
| `SESSION_SECRET` | Yes | - | Session encryption key |

### Example `.env`

```bash
# Server
PORT=3001
FRONTEND_URL=http://localhost:3000
CORS_ORIGINS=http://localhost:3000,https://abc123.ngrok-free.dev

# Spotify
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here
SPOTIFY_REDIRECT_URL=https://abc123.ngrok-free.dev/api/auth/callback

# Session
SESSION_SECRET=your_random_secret_here
```

## How It Works

### 1. Genre Detection

```
Liked Songs → Extract Artist IDs → Batch Fetch Artists → Get Genres → Normalize
```

- Fetches all liked songs (paginated, 50/request)
- Collects unique artist IDs from first artist of each track
- Batch fetches artist data (50 artists/request)
- Extracts primary genre (first genre from artist)
- Normalizes: lowercase, remove special chars, trim whitespace

### 2. Sort Plan Generation

```
Liked Songs + Existing Playlists → Build Mapping → Identify Changes → Generate Plan
```

- Builds track-to-playlist mapping
- Identifies songs not in any genre playlist → add
- Identifies songs in wrong genre playlists → remove
- Groups changes by genre
- Determines new playlists to create

### 3. Execution

```
Plan → Create Playlists → Add Tracks → Remove Tracks → Report Progress
```

- Creates new playlists with managed tag
- Adds tracks in batches (100/request)
- Removes tracks from wrong playlists
- Streams progress via SSE

## Managed Playlists

The app identifies playlists it manages by checking for:
```
[Managed by SpotifyPlaylistSorter]
```
in the playlist description.

**Behavior:**
- Only adds songs to managed playlists
- Only removes songs from managed playlists
- Never modifies your personal playlists
- Creates new playlists with this tag

## OAuth Flow

```
1. Frontend: GET /api/auth/login
   ↓ Returns Spotify auth URL

2. User: Redirects to Spotify
   ↓ User authorizes

3. Spotify: Redirects to ngrok callback
   ↓ https://xxx.ngrok-free.dev/api/auth/callback

4. Backend: Exchanges code for token
   ↓ Creates session, generates temp token

5. Backend: Redirects to frontend
   ↓ http://localhost:3000/callback?token=xxx

6. Frontend: GET /api/auth/complete?token=xxx
   ↓ Sets session cookie on localhost

7. Frontend: Redirects to dashboard
```

## Rate Limiting

- **Request Rate**: 2 requests/second with burst of 5
- **Batch Sizes**:
  - Artists: 50 per request
  - Tracks (add/remove): 100 per request
  - Liked songs: 50 per page
- **Retry**: Automatic retry on 429 errors with backoff

## Troubleshooting

### "Invalid state parameter"
- OAuth state is stored server-side
- Ensure backend didn't restart between login click and callback
- Try logging in again

### "INVALID_CLIENT: Insecure redirect URI"
- Spotify requires HTTPS for callbacks
- Ensure ngrok is running and URL is correct
- Verify redirect URI in Spotify Dashboard matches exactly

### Session not persisting after login
- Cookie is set via `/api/auth/complete` on localhost
- Ensure frontend properly calls this endpoint
- Check browser dev tools for cookie

### "No songs found" but I have liked songs
- Check if Spotify scopes are correct
- Re-authorize the app
- Check backend logs for errors

### Rate limit errors
- Built-in rate limiter should handle this
- For very large libraries, analysis may take several minutes
- Check progress via SSE stream

## Tech Stack

### Backend
| Package | Purpose |
|---------|---------|
| `github.com/zmb3/spotify/v2` | Spotify SDK |
| `github.com/gin-gonic/gin` | HTTP framework |
| `golang.org/x/oauth2` | OAuth 2.0 |
| `golang.org/x/time/rate` | Rate limiting |
| `github.com/google/uuid` | UUID generation |
| `github.com/joho/godotenv` | Env file loading |
| `github.com/rs/zerolog` | Structured logging |

### Frontend
| Package | Purpose |
|---------|---------|
| `react` | UI framework |
| `react-router-dom` | Routing |
| `@tanstack/react-query` | Server state |
| `zustand` | Client state |
| `tailwindcss` | Styling |
| `lucide-react` | Icons |
| `vite` | Build tool |

## License

MIT
