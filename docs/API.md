# API Reference

Complete API documentation for the Spotify Playlist Sorter backend.

## Base URL

All endpoints are prefixed with `/api`.

- **Development**: `http://localhost:3001/api`
- **Via ngrok**: `https://YOUR-NGROK-URL/api`

## Authentication

The API uses session-based authentication with HTTP-only cookies. After successful OAuth login, a `session` cookie is set that authenticates subsequent requests.

### Headers

```
Content-Type: application/json
Cookie: session=<session_id>
```

---

## Endpoints

### Health Check

#### `GET /health`

Check if the server is running.

**Auth Required**: No

**Response**:
```json
{
  "status": "ok"
}
```

---

### Authentication

#### `GET /api/auth/login`

Initiate Spotify OAuth flow. Returns the Spotify authorization URL.

**Auth Required**: No

**Response**:
```json
{
  "url": "https://accounts.spotify.com/authorize?client_id=...&redirect_uri=...&response_type=code&scope=...&state=..."
}
```

**Scopes Requested**:
- `user-read-private` - Read user profile
- `user-read-email` - Read user email
- `user-library-read` - Read liked songs
- `playlist-read-private` - Read user playlists
- `playlist-read-collaborative` - Read collaborative playlists
- `playlist-modify-public` - Create/modify public playlists
- `playlist-modify-private` - Create/modify private playlists

---

#### `GET /api/auth/callback`

OAuth callback endpoint. Called by Spotify after user authorization.

**Auth Required**: No

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `code` | string | Yes | Authorization code from Spotify |
| `state` | string | Yes | CSRF state token |

**Response**: Redirects to `{FRONTEND_URL}/callback?token={temp_token}`

**Errors**:
- `400 Bad Request` - Invalid state parameter or no code provided
- `500 Internal Server Error` - Failed to exchange code or get user profile

---

#### `GET /api/auth/complete`

Complete login by exchanging temp token for session cookie.

**Auth Required**: No

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `token` | string | Yes | Temporary token from callback redirect |

**Response**:
```json
{
  "success": true
}
```

**Side Effects**: Sets `session` cookie on response.

**Errors**:
- `400 Bad Request` - No token provided or invalid/expired token

---

#### `GET /api/auth/me`

Get current authenticated user's profile.

**Auth Required**: Yes

**Response**:
```json
{
  "id": "user123",
  "displayName": "John Doe",
  "email": "john@example.com",
  "imageUrl": "https://i.scdn.co/image/...",
  "product": "premium"
}
```

**Errors**:
- `401 Unauthorized` - Not authenticated or session expired

---

#### `POST /api/auth/logout`

Log out and clear session.

**Auth Required**: Yes

**Response**:
```json
{
  "message": "Logged out successfully"
}
```

**Side Effects**: Clears `session` cookie.

---

### Library

#### `GET /api/library/analysis`

Analyze user's liked songs and detect genres.

**Auth Required**: Yes

**Response**:
```json
{
  "totalTracks": 1500,
  "analyzedTracks": 1500,
  "genres": {
    "rock": {
      "name": "rock",
      "count": 250,
      "tracks": [
        {
          "id": "track123",
          "name": "Song Title",
          "artists": [
            {
              "id": "artist123",
              "name": "Artist Name"
            }
          ],
          "album": {
            "id": "album123",
            "name": "Album Name",
            "imageUrl": "https://i.scdn.co/image/..."
          },
          "genre": "rock"
        }
      ]
    },
    "pop": {
      "name": "pop",
      "count": 180,
      "tracks": [...]
    }
  },
  "uncategorized": {
    "name": "uncategorized",
    "count": 45,
    "tracks": [...]
  }
}
```

**Notes**:
- This endpoint streams progress via SSE (subscribe to `/api/events` first)
- Large libraries may take several minutes to analyze
- Tracks are grouped by normalized genre (lowercase, trimmed)

**Errors**:
- `401 Unauthorized` - Not authenticated
- `500 Internal Server Error` - Failed to fetch library

---

### Sort

#### `POST /api/sort/plan`

Generate a sort plan for organizing tracks into genre playlists.

**Auth Required**: Yes

**Request Body**:
```json
{
  "dryRun": true
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `dryRun` | boolean | No | If true, only preview changes (default: false) |

**Response**:
```json
{
  "id": "plan-uuid-123",
  "dryRun": true,
  "summary": {
    "totalTracks": 500,
    "tracksToAdd": 450,
    "tracksToRemove": 30,
    "playlistsToCreate": 5,
    "playlistsToUpdate": 12
  },
  "operations": [
    {
      "type": "create_playlist",
      "genre": "indie rock",
      "playlistName": "indie rock"
    },
    {
      "type": "add_tracks",
      "genre": "rock",
      "playlistId": "playlist123",
      "playlistName": "rock",
      "trackIds": ["track1", "track2", "track3"],
      "trackCount": 3
    },
    {
      "type": "remove_tracks",
      "genre": "pop",
      "playlistId": "playlist456",
      "playlistName": "pop",
      "trackIds": ["track4"],
      "trackCount": 1
    }
  ]
}
```

**Operation Types**:
- `create_playlist` - Create a new genre playlist
- `add_tracks` - Add tracks to a playlist
- `remove_tracks` - Remove tracks from a playlist (only managed playlists)

**Errors**:
- `401 Unauthorized` - Not authenticated
- `500 Internal Server Error` - Failed to generate plan

---

#### `POST /api/sort/execute`

Execute a previously generated sort plan.

**Auth Required**: Yes

**Request Body**:
```json
{
  "planId": "plan-uuid-123"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `planId` | string | Yes | ID of the plan to execute |

**Response**:
```json
{
  "success": true,
  "executed": {
    "playlistsCreated": 5,
    "tracksAdded": 450,
    "tracksRemoved": 30
  }
}
```

**Notes**:
- Progress is streamed via SSE (subscribe to `/api/events` first)
- Only executes plans that are not in dry-run mode
- Tracks are added/removed in batches of 100

**Errors**:
- `400 Bad Request` - Invalid plan ID or dry-run plan
- `401 Unauthorized` - Not authenticated
- `500 Internal Server Error` - Execution failed

---

### Events (Server-Sent Events)

#### `GET /api/events`

Subscribe to real-time progress updates.

**Auth Required**: Yes

**Response**: SSE stream

**Event Types**:

##### `analysis_progress`
```json
{
  "type": "analysis_progress",
  "data": {
    "phase": "fetching_tracks",
    "current": 500,
    "total": 1500,
    "message": "Fetching liked songs..."
  }
}
```

Phases:
- `fetching_tracks` - Loading liked songs from Spotify
- `fetching_artists` - Loading artist data for genre detection
- `analyzing` - Processing and categorizing tracks
- `complete` - Analysis finished

##### `sort_progress`
```json
{
  "type": "sort_progress",
  "data": {
    "phase": "adding_tracks",
    "current": 100,
    "total": 450,
    "genre": "rock",
    "message": "Adding tracks to rock playlist..."
  }
}
```

Phases:
- `creating_playlists` - Creating new genre playlists
- `adding_tracks` - Adding tracks to playlists
- `removing_tracks` - Removing misplaced tracks
- `complete` - Sort finished

##### `error`
```json
{
  "type": "error",
  "data": {
    "message": "Rate limit exceeded, retrying..."
  }
}
```

**Usage Example** (JavaScript):
```javascript
const eventSource = new EventSource('/api/events', {
  withCredentials: true
});

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
};

eventSource.onerror = (error) => {
  console.error('SSE Error:', error);
  eventSource.close();
};
```

---

#### `POST /api/events/test`

Send a test event (development only).

**Auth Required**: Yes

**Request Body**:
```json
{
  "message": "Test message"
}
```

**Response**:
```json
{
  "sent": true
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `302` | Redirect (OAuth callback) |
| `400` | Bad Request - Invalid parameters |
| `401` | Unauthorized - Not authenticated |
| `429` | Too Many Requests - Rate limited |
| `500` | Internal Server Error |

---

## Rate Limiting

The API implements rate limiting to comply with Spotify's API limits:

- **Request Rate**: 2 requests/second with burst of 5
- **Automatic Retry**: 429 errors are retried with exponential backoff

---

## Data Models

### Track
```typescript
interface Track {
  id: string;
  name: string;
  artists: Artist[];
  album: Album;
  genre: string;
}
```

### Artist
```typescript
interface Artist {
  id: string;
  name: string;
}
```

### Album
```typescript
interface Album {
  id: string;
  name: string;
  imageUrl: string;
}
```

### Genre
```typescript
interface Genre {
  name: string;
  count: number;
  tracks: Track[];
}
```

### SortPlan
```typescript
interface SortPlan {
  id: string;
  dryRun: boolean;
  summary: {
    totalTracks: number;
    tracksToAdd: number;
    tracksToRemove: number;
    playlistsToCreate: number;
    playlistsToUpdate: number;
  };
  operations: SortOperation[];
}
```

### SortOperation
```typescript
interface SortOperation {
  type: 'create_playlist' | 'add_tracks' | 'remove_tracks';
  genre: string;
  playlistId?: string;
  playlistName: string;
  trackIds?: string[];
  trackCount?: number;
}
```
