# Spotify Playlist Sorter - Frontend Project Overview

## Summary

A modern, fully-featured React frontend for the Spotify Playlist Sorter application. Built with Vite, React 18, TypeScript, and Tailwind CSS, featuring real-time updates, a dark Spotify-inspired theme, and a comprehensive component library.

## Key Features

### 1. Authentication
- Spotify OAuth 2.0 integration
- Protected routes with automatic redirects
- Session management via cookies
- User profile display in header

### 2. Real-Time Updates
- Server-Sent Events (SSE) for live progress tracking
- Progress bars with percentage and current/total counts
- Activity logs with status indicators
- Used for library analysis and playlist execution

### 3. Dry-Run Mode
- Toggle between Preview and Live modes
- Persisted to localStorage
- Visual indicators (green for preview, red for live)
- Always visible in header
- Prevents accidental changes in preview mode

### 4. Three Main Views

#### Dashboard (`/`)
- Library statistics (tracks, playlists, genres)
- Library analysis trigger with live progress
- Activity log display

#### Genres (`/genres`)
- Grid view of all detected genres
- Track counts per genre
- New genre indicators
- Uncategorized tracks section

#### Changes (`/changes`)
- Sort plan generation
- Change summary (additions, removals, new playlists)
- Before/after diff view
- Execute button (disabled in dry-run mode)
- Live execution progress

## Tech Stack

### Core
- **Vite** - Fast build tool and dev server
- **React 18** - UI library
- **TypeScript** - Type safety
- **React Router DOM** - Client-side routing

### State Management
- **Zustand** - Lightweight state management
  - Auth store (user, authentication state)
  - UI store (dry-run toggle, persisted)

### Data Fetching
- **@tanstack/react-query** - Server state management
  - Query caching
  - Automatic refetching
  - Loading and error states

### Styling
- **Tailwind CSS** - Utility-first CSS framework
- Custom Spotify color palette
- Dark theme optimized
- Responsive design

### Icons
- **lucide-react** - Modern icon library

## Project Structure

```
frontend/
├── public/              # Static assets
├── src/
│   ├── components/      # React components
│   │   ├── ui/         # Base UI components (17 components)
│   │   ├── layout/     # Header, Sidebar
│   │   ├── songs/      # SongCard, SongList
│   │   ├── genres/     # GenreCard, GenreGrid
│   │   ├── changes/    # ChangeSummary, ChangeDiff, TrackMove
│   │   └── progress/   # ProgressBar, LiveLog
│   ├── hooks/          # Custom React hooks (3)
│   │   ├── useAuth.ts
│   │   ├── useSSE.ts
│   │   └── useApi.ts
│   ├── lib/            # Utilities and types
│   │   ├── api.ts      # API client with fetch wrapper
│   │   └── types.ts    # TypeScript type definitions
│   ├── pages/          # Page components (5)
│   │   ├── Login.tsx
│   │   ├── Callback.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Genres.tsx
│   │   └── Changes.tsx
│   ├── stores/         # Zustand stores (2)
│   │   ├── authStore.ts
│   │   └── uiStore.ts
│   ├── App.tsx         # Main app component with routing
│   ├── main.tsx        # Entry point
│   └── index.css       # Global styles
├── index.html          # HTML template
├── package.json        # Dependencies and scripts
├── vite.config.ts      # Vite configuration
├── tailwind.config.js  # Tailwind configuration
├── tsconfig.json       # TypeScript configuration
├── .eslintrc.cjs       # ESLint configuration
├── .gitignore
├── README.md           # User documentation
├── SETUP.md            # Setup guide
├── COMPONENTS.md       # Component API reference
├── PROJECT_OVERVIEW.md # This file
└── start.sh           # Quick start script
```

## Component Architecture

### UI Components (Reusable)
All UI components support:
- TypeScript props with type safety
- Variants for different styles
- Consistent API patterns
- Forward refs where appropriate
- Tailwind CSS utility classes

Components:
- **Button** - Primary, secondary, danger, ghost variants
- **Card** - Default and highlighted variants, with Header/Title/Content subcomponents
- **Badge** - Success, warning, danger, info variants
- **Progress** - Percentage bar with optional label
- **Toggle** - Switch component for dry-run mode

### Feature Components
Built on top of UI components:
- Domain-specific (songs, genres, changes)
- Accept domain objects as props
- Handle complex layouts and interactions

### Layout Components
- **Header** - Global navigation bar with user info and dry-run toggle
- **Sidebar** - Left navigation with active route highlighting

## State Management Strategy

### Local State (useState)
- Component-specific UI state
- Form inputs
- Temporary flags

### Zustand Stores
- **authStore** - User authentication and profile
- **uiStore** - Global UI preferences (dry-run mode)

### React Query
- Server data caching
- API request states
- Mutations for POST requests

### Server-Sent Events
- Real-time updates from backend
- Custom useSSE hook for connection management
- Event buffering and display

## API Integration

### REST Endpoints
```
GET  /api/auth/login      - Get Spotify OAuth URL
GET  /api/auth/callback   - Exchange code for session
GET  /api/auth/me         - Get current user
POST /api/auth/logout     - Logout
GET  /api/library/stats   - Get library statistics
POST /api/sort/plan       - Create sort plan
POST /api/sort/execute    - Execute sort plan
```

### SSE Endpoints
```
GET /api/events?type=analysis&...  - Library analysis progress
GET /api/events?type=execution&... - Execution progress
```

### API Client (`src/lib/api.ts`)
- Centralized fetch wrapper
- Automatic JSON handling
- Cookie-based auth (credentials: 'include')
- Error handling
- EventSource factory for SSE

## Styling System

### Tailwind Configuration
Custom color palette:
```js
spotify: {
  green: '#1DB954',      // Brand color, CTAs
  black: '#121212',      // Background
  darkgray: '#181818',   // Cards
  gray: '#282828',       // Borders, hover states
  lightgray: '#B3B3B3',  // Secondary text
}
```

### Design Patterns
- Dark theme throughout
- Rounded corners (rounded-lg, rounded-full)
- Subtle borders and shadows
- Hover states for interactivity
- Responsive grid layouts
- Custom scrollbar styling

## Development Workflow

### Starting Development
```bash
npm install          # Install dependencies
npm run dev          # Start dev server (port 3000)
# or
./start.sh          # Run start script (checks backend)
```

### Building
```bash
npm run build       # TypeScript check + Vite build
npm run preview     # Preview production build
```

### Linting
```bash
npm run lint        # ESLint check
```

## Key Implementation Details

### Authentication Flow
1. User visits protected route → redirected to `/login`
2. Click "Login with Spotify" → call `GET /api/auth/login`
3. Redirect to Spotify OAuth page
4. User authorizes → Spotify redirects to `/callback?code=xxx`
5. Callback page calls `GET /api/auth/callback?code=xxx`
6. Backend sets session cookie
7. Frontend calls `GET /api/auth/me` to get user
8. User stored in Zustand → authenticated
9. Redirect to dashboard

### SSE Connection Management
- `useSSE` hook manages EventSource lifecycle
- Accepts endpoint URL or null (for no connection)
- Buffers events in state array
- Provides callbacks for message, error, complete
- Auto-closes on complete/error
- Returns connection status and events

### Dry-Run Mode
- State stored in Zustand with localStorage persistence
- Toggle in header always visible
- Changes page checks mode before executing
- Execute button disabled when dry-run is ON
- Visual warnings when dry-run is OFF

### Route Protection
- `ProtectedRoute` wrapper component
- Checks `useAuth().isAuthenticated`
- Shows loading spinner while checking
- Redirects to `/login` if not authenticated
- Wraps all pages except login and callback

## Future Enhancements

Potential improvements:
- [ ] Playlist detail view
- [ ] Search and filter functionality
- [ ] Bulk genre editing
- [ ] Export sort plan to JSON
- [ ] Undo functionality
- [ ] Dark/light theme toggle
- [ ] Mobile-optimized layout
- [ ] Keyboard shortcuts
- [ ] Analytics dashboard
- [ ] Settings page

## Testing Recommendations

### Unit Tests
- Component rendering
- Hook behavior
- API client methods
- Store actions

### Integration Tests
- Authentication flow
- API integration
- SSE connections
- Route protection

### E2E Tests
- Full user journey
- Login → analyze → sort → execute
- Error scenarios

## Performance Considerations

### Optimizations Applied
- React Query caching
- Lazy loading for routes (can be added)
- Optimized bundle size
- Efficient re-renders with proper hooks
- Debounced SSE event updates

### Bundle Size
Production build:
- CSS: ~15.65 KB (gzipped: 3.74 KB)
- JS: ~235.67 KB (gzipped: 73.06 KB)

## Browser Support

Targets modern browsers with ES2020 support:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## Deployment

### Static Hosting
Build output can be deployed to:
- Vercel
- Netlify
- AWS S3 + CloudFront
- GitHub Pages
- Any static hosting service

### Configuration
Update API URL for production (if not using proxy):
- Modify `src/lib/api.ts` to use `VITE_API_URL` environment variable
- Or configure reverse proxy on hosting platform

## Credits

Built with modern web technologies:
- React, Vite, TypeScript, Tailwind CSS
- Icons by Lucide
- Color palette inspired by Spotify design system
