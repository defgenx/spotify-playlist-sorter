# Architecture

Deep dive into the Spotify Playlist Sorter architecture.

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              User's Browser                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     React Frontend (Port 3000)                       │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │  Login   │  │Dashboard │  │  Genres  │  │ Changes  │            │   │
│  │  │  Page    │  │  Page    │  │   Page   │  │   Page   │            │   │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘            │   │
│  │       │              │              │              │                 │   │
│  │       └──────────────┴──────────────┴──────────────┘                 │   │
│  │                              │                                        │   │
│  │  ┌───────────────────────────┴───────────────────────────┐          │   │
│  │  │              State Management (Zustand)                │          │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌──────────────┐   │          │   │
│  │  │  │ Auth Store  │  │  UI Store   │  │ Query Cache  │   │          │   │
│  │  │  └─────────────┘  └─────────────┘  └──────────────┘   │          │   │
│  │  └───────────────────────────────────────────────────────┘          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└───────────────────────────────────────┬─────────────────────────────────────┘
                                        │
                               Vite Dev Proxy
                                        │
┌───────────────────────────────────────┼─────────────────────────────────────┐
│                                       │                                      │
│  ┌────────────────────────────────────┴────────────────────────────────┐   │
│  │                     Go Backend (Port 3001)                           │   │
│  │                                                                      │   │
│  │  ┌──────────────────────────────────────────────────────────────┐   │   │
│  │  │                    Gin HTTP Router                            │   │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐  │   │   │
│  │  │  │  Auth   │  │ Library │  │  Sort   │  │     Events      │  │   │   │
│  │  │  │ Handler │  │ Handler │  │ Handler │  │(SSE) Handler    │  │   │   │
│  │  │  └────┬────┘  └────┬────┘  └────┬────┘  └────────┬────────┘  │   │   │
│  │  └───────┼────────────┼────────────┼────────────────┼───────────┘   │   │
│  │          │            │            │                │               │   │
│  │  ┌───────┴────────────┴────────────┴────────────────┴───────────┐   │   │
│  │  │                     Service Layer                             │   │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐   │   │   │
│  │  │  │  Library    │  │   Sorter    │  │     Executor        │   │   │   │
│  │  │  │  Service    │  │   Service   │  │     Service         │   │   │   │
│  │  │  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘   │   │   │
│  │  └─────────┼────────────────┼────────────────────┼──────────────┘   │   │
│  │            │                │                    │                  │   │
│  │  ┌─────────┴────────────────┴────────────────────┴──────────────┐   │   │
│  │  │                   Infrastructure Layer                        │   │   │
│  │  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │   │   │
│  │  │  │ Spotify  │  │ Session  │  │   SSE    │  │    Genre     │  │   │   │
│  │  │  │  Client  │  │  Store   │  │Broadcast │  │  Normalizer  │  │   │   │
│  │  │  └────┬─────┘  └──────────┘  └──────────┘  └──────────────┘  │   │   │
│  │  └───────┼──────────────────────────────────────────────────────┘   │   │
│  └──────────┼──────────────────────────────────────────────────────────┘   │
│             │                                                               │
└─────────────┼───────────────────────────────────────────────────────────────┘
              │
     ┌────────┴────────┐
     │   ngrok Tunnel  │
     │    (HTTPS)      │
     └────────┬────────┘
              │
┌─────────────┴────────────────────────────────────────────────────────────────┐
│                              Spotify API                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   OAuth     │  │   Library   │  │   Artists   │  │     Playlists       │  │
│  │  Endpoints  │  │  Endpoints  │  │  Endpoints  │  │     Endpoints       │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Backend Architecture

### Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP request handlers
│   │   │   ├── auth.go          # OAuth & session management
│   │   │   ├── library.go       # Library analysis endpoints
│   │   │   ├── sort.go          # Sort plan & execution
│   │   │   └── events.go        # SSE streaming
│   │   ├── middleware/          # HTTP middleware
│   │   │   ├── auth.go          # Session authentication
│   │   │   └── cors.go          # CORS configuration
│   │   └── router.go            # Route definitions
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── domain/
│   │   ├── track.go             # Track model
│   │   ├── playlist.go          # Playlist model
│   │   └── sortplan.go          # Sort plan model
│   ├── genre/
│   │   └── normalizer.go        # Genre string normalization
│   ├── service/
│   │   ├── library.go           # Library analysis logic
│   │   ├── sorter.go            # Sort plan generation
│   │   └── executor.go          # Plan execution
│   ├── session/
│   │   └── store.go             # In-memory session store
│   ├── spotify/
│   │   └── client.go            # Spotify SDK wrapper
│   └── sse/
│       └── broadcaster.go       # SSE event broadcasting
├── certs/                       # SSL certificates (optional)
├── .env                         # Environment configuration
└── go.mod                       # Go module definition
```

### Component Details

#### Handlers

**AuthHandler** (`handlers/auth.go`)
- Manages OAuth 2.0 flow with Spotify
- Stores OAuth states server-side (not in cookies) for ngrok compatibility
- Uses temp token flow for cross-domain cookie setting
- Handles session creation and token refresh

**LibraryHandler** (`handlers/library.go`)
- Triggers library analysis
- Returns genre-grouped track data
- Broadcasts progress via SSE

**SortHandler** (`handlers/sort.go`)
- Generates sort plans (dry-run or live)
- Executes approved plans
- Reports execution progress

**EventsHandler** (`handlers/events.go`)
- SSE endpoint for real-time updates
- Client connection management
- Event broadcasting

#### Services

**LibraryService** (`service/library.go`)
```go
type LibraryService struct {
    spotifyClient *spotify.Client
    broadcaster   *sse.Broadcaster
}

func (s *LibraryService) AnalyzeLibrary(ctx context.Context, token *oauth2.Token) (*domain.LibraryAnalysis, error)
```
- Fetches all liked songs (paginated, 50/request)
- Collects unique artist IDs
- Batch fetches artist data (50/request)
- Extracts primary genre per track
- Groups tracks by normalized genre

**SorterService** (`service/sorter.go`)
```go
type SorterService struct {
    spotifyClient *spotify.Client
}

func (s *SorterService) GeneratePlan(ctx context.Context, analysis *domain.LibraryAnalysis, playlists []domain.Playlist) *domain.SortPlan
```
- Analyzes current playlist state
- Identifies managed playlists (by description tag)
- Determines tracks to add/remove
- Creates execution plan

**ExecutorService** (`service/executor.go`)
```go
type ExecutorService struct {
    spotifyClient *spotify.Client
    broadcaster   *sse.Broadcaster
}

func (s *ExecutorService) Execute(ctx context.Context, plan *domain.SortPlan, token *oauth2.Token) error
```
- Creates new playlists with managed tag
- Adds tracks in batches (100/request)
- Removes tracks from wrong playlists
- Broadcasts progress updates

#### Infrastructure

**SpotifyClient** (`spotify/client.go`)
- Wraps `zmb3/spotify/v2` SDK
- Configures OAuth 2.0
- Provides rate-limited API access
- Handles token refresh

**SessionStore** (`session/store.go`)
```go
type Session struct {
    UserID    string
    Token     *oauth2.Token
    CreatedAt time.Time
}

type Store struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}
```
- In-memory session storage
- Thread-safe access
- Token management

**SSE Broadcaster** (`sse/broadcaster.go`)
```go
type Broadcaster struct {
    clients map[string]chan Event
    mu      sync.RWMutex
}

func (b *Broadcaster) Subscribe(userID string) <-chan Event
func (b *Broadcaster) Broadcast(userID string, event Event)
```
- Per-user event channels
- Non-blocking broadcast
- Client lifecycle management

### Data Flow

#### Library Analysis Flow

```
1. Frontend: GET /api/library/analysis
   │
2. Handler: Validate session, get token
   │
3. LibraryService.AnalyzeLibrary():
   │
   ├─► Broadcast: "fetching_tracks" progress
   │
   ├─► Spotify API: GET /me/tracks (paginated)
   │   └── 50 tracks per request
   │
   ├─► Broadcast: "fetching_artists" progress
   │
   ├─► Spotify API: GET /artists (batched)
   │   └── 50 artists per request
   │
   ├─► Process: Extract primary genre from each artist
   │
   ├─► GenreNormalizer: Normalize genre strings
   │
   ├─► Broadcast: "complete"
   │
   └─► Return: LibraryAnalysis with grouped tracks
```

#### Sort Plan Flow

```
1. Frontend: POST /api/sort/plan
   │
2. Handler: Get analysis + existing playlists
   │
3. SorterService.GeneratePlan():
   │
   ├─► Filter: Identify managed playlists
   │
   ├─► Build: Track-to-playlist mapping
   │
   ├─► Compare: Current state vs desired state
   │
   ├─► Calculate: Required operations
   │   ├── Playlists to create
   │   ├── Tracks to add
   │   └── Tracks to remove
   │
   └─► Return: SortPlan with operations
```

#### Execution Flow

```
1. Frontend: POST /api/sort/execute
   │
2. Handler: Validate plan, get token
   │
3. ExecutorService.Execute():
   │
   ├─► Broadcast: "creating_playlists"
   │
   ├─► For each new genre:
   │   └── Spotify API: POST /users/{id}/playlists
   │
   ├─► Broadcast: "adding_tracks"
   │
   ├─► For each playlist:
   │   └── Spotify API: POST /playlists/{id}/tracks (batched)
   │
   ├─► Broadcast: "removing_tracks"
   │
   ├─► For each removal:
   │   └── Spotify API: DELETE /playlists/{id}/tracks
   │
   └─► Broadcast: "complete"
```

---

## Frontend Architecture

### Directory Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/                  # Base UI components
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Badge.tsx
│   │   │   ├── Toggle.tsx
│   │   │   └── Spinner.tsx
│   │   ├── layout/              # Layout components
│   │   │   ├── Header.tsx
│   │   │   └── Sidebar.tsx
│   │   ├── songs/               # Song display components
│   │   │   ├── SongCard.tsx
│   │   │   └── SongList.tsx
│   │   ├── genres/              # Genre display components
│   │   │   ├── GenreCard.tsx
│   │   │   └── GenreGrid.tsx
│   │   ├── changes/             # Change preview components
│   │   │   ├── ChangeSummary.tsx
│   │   │   └── ChangeDiff.tsx
│   │   └── progress/            # Progress display
│   │       ├── ProgressBar.tsx
│   │       └── LiveLog.tsx
│   ├── hooks/
│   │   ├── useAuth.ts           # Authentication hook
│   │   └── useSSE.ts            # SSE subscription hook
│   ├── stores/
│   │   ├── authStore.ts         # Auth state (Zustand)
│   │   └── uiStore.ts           # UI state (dry-run toggle)
│   ├── lib/
│   │   ├── api.ts               # API client
│   │   └── types.ts             # TypeScript types
│   ├── pages/
│   │   ├── Login.tsx            # Login page
│   │   ├── Callback.tsx         # OAuth callback handler
│   │   ├── Dashboard.tsx        # Main dashboard
│   │   ├── Genres.tsx           # Genre view page
│   │   └── Changes.tsx          # Changes preview page
│   ├── App.tsx                  # Root component
│   └── main.tsx                 # Entry point
├── index.html
├── tailwind.config.js
├── tsconfig.json
├── vite.config.ts
└── package.json
```

### State Management

#### Auth Store (Zustand)
```typescript
interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  checkAuth: () => Promise<void>;
  logout: () => Promise<void>;
}
```

#### UI Store (Zustand)
```typescript
interface UIState {
  isDryRun: boolean;
  setDryRun: (value: boolean) => void;
}
```

#### Server State (TanStack Query)
- Library analysis data
- Sort plan data
- Automatic caching and invalidation

### Component Hierarchy

```
App
├── Header
│   ├── Logo
│   └── UserMenu
├── Sidebar
│   └── Navigation
└── Routes
    ├── Login
    │   └── SpotifyLoginButton
    ├── Callback
    │   └── LoadingSpinner
    ├── Dashboard
    │   ├── StatsCards
    │   ├── ProgressBar
    │   └── ActionButtons
    ├── Genres
    │   └── GenreGrid
    │       └── GenreCard[]
    │           └── SongList
    │               └── SongCard[]
    └── Changes
        ├── ChangeSummary
        ├── DryRunToggle
        ├── ChangeDiff
        │   └── OperationList
        └── ExecuteButton
```

---

## OAuth Flow

The OAuth flow handles cross-domain authentication with ngrok:

```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│     Frontend    │         │     Backend     │         │    Spotify      │
│  (localhost)    │         │  (via ngrok)    │         │                 │
└────────┬────────┘         └────────┬────────┘         └────────┬────────┘
         │                           │                           │
         │  1. GET /api/auth/login   │                           │
         │─────────────────────────► │                           │
         │                           │                           │
         │  2. Return auth URL       │                           │
         │◄───────────────────────── │                           │
         │                           │                           │
         │  3. Open Spotify auth     │                           │
         │───────────────────────────┼─────────────────────────► │
         │                           │                           │
         │                           │  4. User authorizes       │
         │                           │◄───────────────────────── │
         │                           │                           │
         │                           │  5. Redirect to ngrok URL │
         │                           │◄───────────────────────── │
         │                           │     /api/auth/callback    │
         │                           │                           │
         │                           │  6. Exchange code for     │
         │                           │     token, create session │
         │                           │                           │
         │  7. Redirect to frontend  │                           │
         │◄───────────────────────── │                           │
         │     /callback?token=xxx   │                           │
         │                           │                           │
         │  8. GET /api/auth/complete│                           │
         │─────────────────────────► │                           │
         │                           │                           │
         │  9. Set session cookie    │                           │
         │◄───────────────────────── │                           │
         │     (on localhost)        │                           │
         │                           │                           │
         │  10. Navigate to /        │                           │
         │                           │                           │
```

**Key Design Decisions:**

1. **Server-side state storage**: OAuth state is stored in server memory, not cookies, to handle cross-domain redirect through ngrok.

2. **Temp token exchange**: After Spotify callback, backend creates a temp token and redirects to frontend. Frontend then calls `/api/auth/complete` to set the session cookie on localhost.

3. **Session cookies on localhost**: Final session cookie is set on localhost domain, enabling API calls through the Vite proxy.

---

## Rate Limiting Strategy

### Spotify API Limits
- 180 requests per minute (soft limit)
- Varies by endpoint

### Implementation

**Request Rate Limiter**
```go
limiter := rate.NewLimiter(rate.Limit(2), 5) // 2 req/sec, burst of 5
```

**Batch Sizes**
- Artists: 50 per request (API limit)
- Tracks add/remove: 100 per request (API limit)
- Liked songs: 50 per page (API limit)

**Retry Strategy**
```go
func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err := c.httpClient.Do(req)
        if resp.StatusCode == 429 {
            retryAfter := parseRetryAfter(resp.Header)
            time.Sleep(retryAfter)
            continue
        }
        return resp, err
    }
}
```

---

## Managed Playlist Identification

Playlists are identified as "managed" by checking the description:

```go
const ManagedTag = "[Managed by SpotifyPlaylistSorter]"

func isManaged(playlist *spotify.SimplePlaylist) bool {
    return strings.Contains(playlist.Description, ManagedTag)
}
```

**New Playlist Creation**
```go
func createPlaylist(ctx context.Context, client *spotify.Client, userID, genre string) (*spotify.FullPlaylist, error) {
    return client.CreatePlaylistForUser(ctx, userID, genre, fmt.Sprintf(
        "%s playlist. %s",
        genre,
        ManagedTag,
    ), false, false)
}
```

**Behavior**
- Only managed playlists are modified during sorting
- User's personal playlists are never touched
- Tracks are only removed from managed playlists
- New playlists are always created with the managed tag

---

## Error Handling

### Backend Errors

```go
// Structured error response
type ErrorResponse struct {
    Error string `json:"error"`
}

// Handler error
func (h *Handler) SomeEndpoint(c *gin.Context) {
    if err != nil {
        log.Error().Err(err).Msg("Operation failed")
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: "User-friendly error message",
        })
        return
    }
}
```

### Frontend Errors

```typescript
// API client error handling
async fetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE}${endpoint}`, options);

    if (!response.ok) {
        const error = await response.text();
        throw new Error(error || `HTTP ${response.status}`);
    }

    return response.json();
}

// Component error handling with TanStack Query
const { data, error, isError } = useQuery({
    queryKey: ['library'],
    queryFn: () => api.getLibraryStats(),
});

if (isError) {
    return <ErrorMessage message={error.message} />;
}
```

---

## Security Considerations

1. **CSRF Protection**: OAuth state parameter prevents CSRF attacks
2. **HttpOnly Cookies**: Session cookies are HttpOnly to prevent XSS theft
3. **Token Storage**: Tokens stored server-side only, never exposed to frontend
4. **Managed Playlists**: Only app-created playlists can be modified
5. **Input Validation**: All user inputs validated before processing
6. **CORS**: Strict origin checking for API endpoints
