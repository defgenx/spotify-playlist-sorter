# Application Flow Diagram

## User Journey Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         INITIAL LOAD                                 │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │  Check Auth     │
                        │  (useAuth hook) │
                        └────────┬────────┘
                                 │
                    ┌────────────┴────────────┐
                    ▼                         ▼
            ┌───────────────┐         ┌──────────────┐
            │ Authenticated │         │ Not Authed   │
            └───────┬───────┘         └──────┬───────┘
                    │                        │
                    ▼                        ▼
            ┌───────────────┐         ┌──────────────┐
            │  Dashboard    │         │ Login Page   │
            └───────────────┘         └──────┬───────┘
                    │                        │
                    │                        ▼
                    │                ┌──────────────────┐
                    │                │ Click "Login"    │
                    │                │ with Spotify     │
                    │                └────────┬─────────┘
                    │                         │
                    │                         ▼
                    │                ┌──────────────────┐
                    │                │ Redirect to      │
                    │                │ Spotify OAuth    │
                    │                └────────┬─────────┘
                    │                         │
                    │                         ▼
                    │                ┌──────────────────┐
                    │                │ User Authorizes  │
                    │                └────────┬─────────┘
                    │                         │
                    │                         ▼
                    │                ┌──────────────────┐
                    │                │ Callback Page    │
                    │                │ (/callback?code) │
                    │                └────────┬─────────┘
                    │                         │
                    │                         ▼
                    │                ┌──────────────────┐
                    │                │ Exchange Code    │
                    │                │ for Session      │
                    │                └────────┬─────────┘
                    │                         │
                    └─────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                    AUTHENTICATED USER FLOW                           │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────┐
│  DASHBOARD  │
└──────┬──────┘
       │
       ├─→ View Library Stats
       │
       └─→ Click "Start Analysis"
              │
              ▼
           ┌──────────────────┐
           │ SSE Connection   │
           │ /api/events      │
           └────────┬─────────┘
                    │
                    ├─→ Progress Events
                    ├─→ Update Progress Bar
                    ├─→ Display Live Log
                    │
                    ▼
           ┌──────────────────┐
           │ Analysis Complete│
           └──────────────────┘

┌─────────────┐
│   GENRES    │
└──────┬──────┘
       │
       ├─→ Click "Load Genres"
       │      │
       │      ▼
       │   ┌──────────────────┐
       │   │ POST /sort/plan  │
       │   └────────┬─────────┘
       │            │
       │            ▼
       └─→ Display Genre Grid
              │
              ├─→ Genre Cards with Stats
              ├─→ New Genre Badges
              └─→ Uncategorized Tracks

┌─────────────┐
│   CHANGES   │
└──────┬──────┘
       │
       ├─→ Click "Generate Plan"
       │      │
       │      ▼
       │   ┌──────────────────┐
       │   │ POST /sort/plan  │
       │   │ (with dry-run)   │
       │   └────────┬─────────┘
       │            │
       │            ▼
       ├─→ Display Summary
       │      │
       │      ├─→ Tracks to Add
       │      ├─→ Tracks to Remove
       │      └─→ Playlists to Create
       │
       ├─→ Display Change Diff
       │      │
       │      ├─→ Additions List
       │      └─→ Removals List
       │
       └─→ Click "Execute Changes"
              │
              ├─→ Check Dry-Run Mode
              │      │
              │      ├─→ ON: Button Disabled
              │      └─→ OFF: Execute Allowed
              │
              ▼
           ┌──────────────────┐
           │ POST             │
           │ /sort/execute    │
           └────────┬─────────┘
                    │
                    ▼
           ┌──────────────────┐
           │ SSE Connection   │
           │ /api/events      │
           └────────┬─────────┘
                    │
                    ├─→ Execution Events
                    ├─→ Update Progress
                    ├─→ Display Live Log
                    │
                    ▼
           ┌──────────────────┐
           │ Execution        │
           │ Complete         │
           └──────────────────┘
```

## Component Hierarchy

```
App
├── BrowserRouter
│   └── Routes
│       ├── /login → Login
│       │
│       ├── /callback → Callback
│       │
│       └── ProtectedRoute
│           └── AppLayout
│               ├── Header
│               │   ├── Logo
│               │   ├── Dry-Run Toggle
│               │   └── User Profile
│               │       └── Logout Button
│               │
│               ├── Sidebar
│               │   ├── NavLink (Dashboard)
│               │   ├── NavLink (Genres)
│               │   └── NavLink (Changes)
│               │
│               └── Main Content
│                   │
│                   ├── / → Dashboard
│                   │   ├── Stat Cards
│                   │   ├── Analysis Section
│                   │   │   ├── Button
│                   │   │   └── ProgressBar
│                   │   └── LiveLog
│                   │
│                   ├── /genres → Genres
│                   │   ├── Summary Cards
│                   │   ├── GenreGrid
│                   │   │   └── GenreCard[]
│                   │   └── Uncategorized Section
│                   │
│                   └── /changes → Changes
│                       ├── Warning Banner (if live)
│                       ├── ChangeSummary
│                       │   └── Stat Cards
│                       ├── New Playlists Section
│                       ├── ChangeDiff
│                       │   ├── Additions
│                       │   │   └── TrackMove[]
│                       │   └── Removals
│                       │       └── TrackMove[]
│                       └── LiveLog
```

## State Management Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         STATE LAYERS                                 │
└─────────────────────────────────────────────────────────────────────┘

┌──────────────────┐
│  Zustand Stores  │
└────────┬─────────┘
         │
         ├─→ authStore
         │   ├─ user: User | null
         │   ├─ isLoading: boolean
         │   ├─ isAuthenticated: boolean
         │   └─ actions: setUser, setLoading, logout
         │
         └─→ uiStore (persisted to localStorage)
             ├─ isDryRun: boolean
             └─ actions: toggleDryRun, setDryRun

┌──────────────────┐
│  React Query     │
└────────┬─────────┘
         │
         ├─→ useLibraryStats
         │   └─ Cache library statistics
         │
         ├─→ useCreateSortPlan (mutation)
         │   └─ Create and cache sort plan
         │
         └─→ useExecuteSortPlan (mutation)
             └─ Execute plan (no cache)

┌──────────────────┐
│  Local State     │
│  (useState)      │
└────────┬─────────┘
         │
         ├─→ SSE endpoint URLs
         ├─→ Collapsible section states
         ├─→ Loading flags
         └─→ Error messages

┌──────────────────┐
│  SSE State       │
│  (useSSE)        │
└────────┬─────────┘
         │
         ├─→ events: ProgressEvent[]
         ├─→ isConnected: boolean
         ├─→ error: string | null
         └─→ Manages EventSource lifecycle
```

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         DATA FLOW                                    │
└─────────────────────────────────────────────────────────────────────┘

User Action
    │
    ▼
React Component
    │
    ├─→ useAuth hook
    │   └─→ authStore (Zustand)
    │       └─→ api.ts
    │           └─→ Backend API
    │
    ├─→ useApi hook
    │   └─→ React Query
    │       └─→ api.ts
    │           └─→ Backend API
    │               │
    │               ▼
    │           Response
    │               │
    │               ▼
    │           React Query Cache
    │               │
    │               ▼
    │           Component Re-render
    │
    └─→ useSSE hook
        └─→ EventSource
            └─→ Backend SSE Endpoint
                │
                ├─→ Event Stream
                │       │
                │       ▼
                │   Event Buffer (state)
                │       │
                │       ▼
                │   Component Re-render
                │       │
                │       ▼
                │   LiveLog Component
                │
                └─→ Connection closed
                    └─→ Cleanup
```

## API Call Patterns

```
┌─────────────────────────────────────────────────────────────────────┐
│                      API INTEGRATION                                 │
└─────────────────────────────────────────────────────────────────────┘

REST API Calls:
──────────────

Component → useApi hook → React Query → api.ts → fetch()
                                                      │
                                                      ▼
                                            Backend (localhost:8080)
                                                      │
                                                      ▼
                                                  Response
                                                      │
                                                      ▼
                                            React Query Cache
                                                      │
                                                      ▼
                                            Component Update

Server-Sent Events (SSE):
─────────────────────────

Component → useSSE hook → EventSource → Backend SSE endpoint
                              │
                              ├─→ onopen → setConnected(true)
                              │
                              ├─→ onmessage → buffer event → trigger callback
                              │
                              ├─→ onerror → setError → close connection
                              │
                              └─→ cleanup → close connection

Vite Dev Proxy:
───────────────

Frontend (localhost:3000) → /api/* → Vite Proxy → Backend (localhost:8080)
```

## Authentication Flow Detail

```
┌─────────────────────────────────────────────────────────────────────┐
│                    AUTHENTICATION FLOW                               │
└─────────────────────────────────────────────────────────────────────┘

1. Initial Load
   │
   ├─→ App.tsx mounts
   │   └─→ useAuth() called
   │       └─→ useEffect runs
   │           └─→ checkAuth()
   │               └─→ api.getCurrentUser()
   │                   │
   │                   ├─→ Success: setUser(user)
   │                   └─→ Error: setUser(null)
   │
   ├─→ isAuthenticated = !!user
   │
   └─→ Route decision
       ├─→ Authenticated → Dashboard
       └─→ Not authenticated → Login

2. Login Process
   │
   ├─→ User clicks "Login with Spotify"
   │   └─→ login() called
   │       └─→ api.getLoginUrl()
   │           └─→ Backend returns OAuth URL
   │               └─→ window.location.href = url
   │                   └─→ Redirect to Spotify
   │
   ├─→ User authorizes on Spotify
   │   └─→ Spotify redirects to /callback?code=xxx
   │
   ├─→ Callback component
   │   └─→ Extract code from URL
   │       └─→ api.handleCallback(code)
   │           └─→ Backend sets HTTP-only cookie
   │               └─→ checkAuth() to load user
   │                   └─→ navigate('/') to Dashboard
   │
   └─→ Authenticated state set
       └─→ User can access protected routes

3. Logout Process
   │
   └─→ User clicks logout
       └─→ logout() called
           ├─→ api.logout() → clear backend session
           └─→ logoutStore() → clear frontend state
               └─→ Redirect to /login
```

## Error Handling Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                     ERROR HANDLING                                   │
└─────────────────────────────────────────────────────────────────────┘

API Error:
──────────
fetch() → Error
    │
    ├─→ React Query error state
    │   └─→ Component receives error
    │       └─→ Display error message
    │
    └─→ useAuth error
        └─→ setUser(null)
            └─→ Redirect to login

SSE Error:
──────────
EventSource → onerror
    │
    ├─→ setError('Connection lost')
    ├─→ setConnected(false)
    ├─→ Close connection
    │
    └─→ Component displays error
        └─→ User can retry

Network Error:
──────────────
Backend unavailable
    │
    ├─→ fetch() throws
    │   └─→ "Failed to fetch"
    │       └─→ Display user-friendly message
    │
    └─→ SSE connection fails
        └─→ "Connection lost"
            └─→ Retry button available
```

This comprehensive flow diagram shows how all parts of the application work together!
