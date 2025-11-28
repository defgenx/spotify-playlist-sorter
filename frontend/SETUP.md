# Setup Guide

## Prerequisites

- Node.js 18+ and npm
- Backend server running at `http://localhost:8080`

## Installation Steps

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Start the development server:**
   ```bash
   npm run dev
   ```

3. **Access the application:**
   Open your browser and navigate to `http://localhost:3000`

## Available Scripts

- `npm run dev` - Start development server (port 3000)
- `npm run build` - Build for production
- `npm run preview` - Preview production build locally
- `npm run lint` - Run ESLint

## Configuration

### API Proxy

The frontend is configured to proxy API requests to the backend:
- Development: `http://localhost:8080/api`
- Configured in `vite.config.ts`

### Tailwind Colors

Custom Spotify-themed colors are defined in `tailwind.config.js`:
- `spotify-green`: #1DB954 (Spotify brand color)
- `spotify-black`: #121212 (Background)
- `spotify-darkgray`: #181818
- `spotify-gray`: #282828
- `spotify-lightgray`: #B3B3B3

## Features

### Authentication Flow

1. User clicks "Login with Spotify" on `/login`
2. App redirects to Spotify OAuth page
3. After authorization, Spotify redirects to `/callback?code=xxx`
4. Callback page exchanges code for session
5. User is redirected to dashboard

### Dry-Run Mode

- Toggle in the header (always visible when authenticated)
- Persisted to localStorage via Zustand
- When enabled (green): Preview mode - no changes applied
- When disabled (red): Live mode - changes applied to Spotify

### Real-Time Updates

The app uses Server-Sent Events (SSE) for real-time progress updates:
- Library analysis progress
- Playlist execution progress
- Live logs displayed in the UI

## Project Structure Overview

```
src/
├── components/        # React components
│   ├── ui/           # Base UI components (Button, Card, Badge, etc.)
│   ├── layout/       # Layout components (Header, Sidebar)
│   ├── songs/        # Song display components
│   ├── genres/       # Genre display components
│   ├── changes/      # Change preview components
│   └── progress/     # Progress tracking components
├── hooks/            # Custom React hooks
├── lib/              # Utilities and types
├── pages/            # Page components
└── stores/           # Zustand state stores
```

## Troubleshooting

### Port Already in Use

If port 3000 is already in use, you can change it in `vite.config.ts`:
```ts
server: {
  port: 3001, // Change to any available port
}
```

### Backend Connection Issues

Ensure the backend is running on port 8080. Check the proxy configuration in `vite.config.ts`.

### Build Errors

Clear node_modules and reinstall:
```bash
rm -rf node_modules package-lock.json
npm install
```

## Development Tips

1. **Hot Module Replacement**: Changes are automatically reflected in the browser
2. **TypeScript**: Use TypeScript for type safety - check `src/lib/types.ts` for available types
3. **Component Structure**: Follow the existing component patterns for consistency
4. **State Management**: Use Zustand stores for global state, React hooks for local state

## Production Build

1. Build the application:
   ```bash
   npm run build
   ```

2. The build output will be in the `dist/` folder

3. Preview the production build:
   ```bash
   npm run preview
   ```

4. Deploy the `dist/` folder to your hosting service

## Environment Variables

This app doesn't require environment variables for development. The API proxy is configured in `vite.config.ts`.

For production deployments, you may need to configure:
- API base URL (if not using proxy)
- OAuth redirect URLs (configured in backend)
