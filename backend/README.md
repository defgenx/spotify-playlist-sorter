# Spotify Playlist Sorter - Backend

A Go backend service that automatically organizes your Spotify liked songs into genre-based playlists.

## Features

- OAuth authentication with Spotify
- Fetch and analyze liked songs and playlists
- Automatic genre detection from artist metadata
- Intelligent playlist organization with fuzzy genre matching
- Real-time progress updates via Server-Sent Events (SSE)
- Dry-run mode for previewing changes
- Managed playlists with automatic tagging

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go            # Authentication endpoints
│   │   │   ├── library.go         # Library analysis endpoints
│   │   │   ├── sort.go            # Sort plan endpoints
│   │   │   └── events.go          # SSE event streaming
│   │   ├── middleware/
│   │   │   ├── auth.go            # Session validation
│   │   │   └── cors.go            # CORS handling
│   │   └── router.go              # Route configuration
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── domain/
│   │   ├── track.go               # Track domain model
│   │   ├── playlist.go            # Playlist domain model
│   │   └── sortplan.go            # Sort plan models
│   ├── genre/
│   │   └── normalizer.go          # Genre normalization & matching
│   ├── service/
│   │   ├── library.go             # Library analysis service
│   │   ├── sorter.go              # Sort plan generation
│   │   └── executor.go            # Sort plan execution
│   ├── session/
│   │   └── store.go               # In-memory session store
│   ├── spotify/
│   │   └── client.go              # Spotify API client wrapper
│   └── sse/
│       └── broadcaster.go         # SSE event broadcaster
├── .env.example                    # Example environment variables
├── go.mod                          # Go module definition
└── README.md                       # This file
```

## Prerequisites

- Go 1.24.1 or later
- Spotify Developer Account
- Spotify App credentials (Client ID & Secret)

## Setup

1. **Clone the repository**

```bash
cd backend
```

2. **Install dependencies**

```bash
go mod download
```

3. **Configure environment variables**

Copy `.env.example` to `.env` and fill in your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your Spotify app credentials:

```env
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here
SESSION_SECRET=generate_a_random_secret_here
```

4. **Get Spotify Credentials**

- Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
- Create a new app
- Add `http://localhost:8080/api/auth/callback` to Redirect URIs
- Copy Client ID and Client Secret

## Running the Server

### Development

```bash
go run cmd/server/main.go
```

### Production Build

```bash
go build -o spotify-sorter cmd/server/main.go
./spotify-sorter
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### Authentication

- `GET /api/auth/login` - Get Spotify OAuth URL
- `GET /api/auth/callback` - OAuth callback handler
- `GET /api/auth/me` - Get current user profile
- `POST /api/auth/logout` - Logout current user

### Library

- `GET /api/library/analysis` - Analyze user's library (tracks, playlists, genres)

### Sort

- `POST /api/sort/plan` - Generate a sort plan
  - Body: `{"dryRun": true}`
- `POST /api/sort/execute` - Execute sort plan
  - Body: `{"dryRun": false}`

### Events

- `GET /api/events` - SSE stream for real-time progress updates

### Health

- `GET /health` - Health check endpoint

## How It Works

1. **Authentication**: Users authenticate via Spotify OAuth
2. **Library Analysis**:
   - Fetches all liked songs
   - Fetches all user playlists
   - Retrieves artist information and genres
   - Assigns primary genre to each track
3. **Sort Plan Generation**:
   - Identifies tracks not in genre playlists
   - Identifies tracks in wrong playlists (only managed ones)
   - Determines which playlists need to be created
   - Groups uncategorized tracks
4. **Execution**:
   - Creates new genre playlists with managed tag
   - Adds tracks to correct playlists
   - Removes tracks from incorrect managed playlists
   - Creates "Uncategorized" playlist for tracks without genres

## Managed Playlists

Playlists created by this app are tagged with `[Managed by SpotifyPlaylistSorter]` in their description. Only managed playlists will be modified by the sorter (tracks can be removed). This prevents accidentally modifying user-curated playlists.

## Development

### Code Organization

- **cmd/**: Application entry points
- **internal/api/**: HTTP handlers and routing
- **internal/domain/**: Domain models
- **internal/service/**: Business logic
- **internal/spotify/**: Spotify API wrapper with rate limiting
- **internal/genre/**: Genre normalization and matching
- **internal/session/**: Session management
- **internal/sse/**: Server-Sent Events broadcaster

### Key Dependencies

- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/zmb3/spotify/v2` - Spotify API client
- `github.com/rs/zerolog` - Structured logging
- `github.com/caarlos0/env/v10` - Environment variable parsing
- `golang.org/x/oauth2` - OAuth2 client

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `FRONTEND_URL` | Frontend URL for redirects | `http://localhost:5173` |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | `http://localhost:5173` |
| `SPOTIFY_CLIENT_ID` | Spotify app client ID | *required* |
| `SPOTIFY_CLIENT_SECRET` | Spotify app client secret | *required* |
| `SPOTIFY_REDIRECT_URL` | OAuth callback URL | `http://localhost:8080/api/auth/callback` |
| `SESSION_SECRET` | Secret for session encryption | *required* |

## License

MIT
