# Spotify Playlist Sorter - Frontend

A modern React frontend for automatically organizing Spotify playlists by genre.

## Tech Stack

- **Framework**: Vite + React 18 + TypeScript
- **Styling**: Tailwind CSS (dark theme, Spotify-inspired)
- **State Management**: Zustand
- **Data Fetching**: @tanstack/react-query
- **Routing**: react-router-dom
- **Icons**: lucide-react

## Features

- Spotify OAuth authentication
- Real-time progress tracking via Server-Sent Events (SSE)
- Dry-run mode for previewing changes before applying
- Genre-based playlist organization
- Before/after diff view for track changes
- Responsive design with dark theme

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend server running at `http://localhost:8080`

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
src/
├── components/
│   ├── ui/           # Reusable UI components
│   ├── layout/       # Header, Sidebar
│   ├── songs/        # Song-related components
│   ├── genres/       # Genre display components
│   ├── changes/      # Change preview components
│   └── progress/     # Progress tracking components
├── hooks/
│   ├── useAuth.ts    # Authentication hook
│   ├── useSSE.ts     # Server-Sent Events hook
│   └── useApi.ts     # API query hooks
├── stores/
│   ├── authStore.ts  # Auth state management
│   └── uiStore.ts    # UI state (dry-run toggle)
├── lib/
│   ├── api.ts        # API client
│   └── types.ts      # TypeScript type definitions
├── pages/
│   ├── Login.tsx
│   ├── Callback.tsx
│   ├── Dashboard.tsx
│   ├── Genres.tsx
│   └── Changes.tsx
├── App.tsx
└── main.tsx
```

## Environment

The frontend expects the backend API to be available at `http://localhost:8080/api`.
This is configured in `vite.config.ts` as a proxy.

## Key Features Explained

### Dry-Run Mode

The dry-run toggle in the header allows users to preview changes without applying them.
This state is persisted in localStorage via Zustand.

### Server-Sent Events

The `useSSE` hook connects to the backend's SSE endpoint for real-time progress updates
during library analysis and playlist execution.

### Route Protection

All routes except `/login` and `/callback` are protected and require authentication.
The `ProtectedRoute` component handles this logic.

## License

MIT
