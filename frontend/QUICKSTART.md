# Quick Start Guide

## Prerequisites Check

Before starting, ensure you have:
- ✅ Node.js 18+ installed (`node --version`)
- ✅ npm installed (`npm --version`)
- ✅ Backend server running on `http://localhost:8080`

## Installation (First Time)

```bash
# Navigate to the frontend directory
cd /Users/adelvecchio/Documents/workspace/spotify-playlist-sorter/frontend

# Install dependencies
npm install
```

## Running the Application

### Option 1: Quick Start Script (Recommended)
```bash
./start.sh
```

This script will:
- Check if dependencies are installed (installs if needed)
- Verify backend is running
- Start the development server

### Option 2: Manual Start
```bash
npm run dev
```

### Access the Application
Open your browser and go to: **http://localhost:3000**

## First-Time User Flow

1. **Login**: Click "Login with Spotify" on the login page
2. **Authorize**: Grant permissions on Spotify's authorization page
3. **Redirect**: You'll be redirected back to the app
4. **Dashboard**: Start by clicking "Start Analysis" to analyze your library
5. **Genres**: View detected genres in the Genres tab
6. **Changes**: Generate a sort plan and preview changes in the Changes tab
7. **Execute**: Toggle off "Dry Run" and click "Execute Changes" to apply

## Available Commands

```bash
npm run dev      # Start development server (port 3000)
npm run build    # Build for production
npm run preview  # Preview production build
npm run lint     # Run ESLint
```

## Features Overview

### Dry-Run Mode (Preview Mode)
- **Location**: Toggle in the header (top right)
- **Default**: ON (green badge)
- **Purpose**: Preview changes without applying them
- **Saved**: Your preference is saved in localStorage

### Three Main Pages

1. **Dashboard** (`/`)
   - View library statistics
   - Run analysis on your Spotify library
   - See real-time progress

2. **Genres** (`/genres`)
   - View all detected genres
   - See track counts per genre
   - Identify new genres

3. **Changes** (`/changes`)
   - Generate sort plan
   - Preview track additions/removals
   - Execute changes (when dry-run is OFF)

## Troubleshooting

### "Connection refused" or API errors
**Problem**: Backend is not running
**Solution**: Start the backend server first on port 8080

### Port 3000 already in use
**Problem**: Another app is using port 3000
**Solution**:
1. Stop the other app, or
2. Edit `vite.config.ts` and change the port:
```ts
server: {
  port: 3001, // Change to any available port
}
```

### Build errors after pulling changes
**Problem**: Stale dependencies
**Solution**:
```bash
rm -rf node_modules package-lock.json
npm install
```

### Styles not loading correctly
**Problem**: Tailwind CSS not compiled
**Solution**: Restart the dev server
```bash
# Press Ctrl+C to stop
npm run dev
```

## Project Structure (Quick Reference)

```
frontend/
├── src/
│   ├── pages/          # Login, Dashboard, Genres, Changes
│   ├── components/     # Reusable React components
│   ├── hooks/          # useAuth, useSSE, useApi
│   ├── stores/         # Zustand state stores
│   └── lib/            # API client and types
├── public/             # Static assets
├── package.json        # Dependencies and scripts
├── vite.config.ts      # Vite configuration (API proxy)
└── tailwind.config.js  # Tailwind theme configuration
```

## Getting Help

- **Setup Issues**: See `SETUP.md`
- **Component Usage**: See `COMPONENTS.md`
- **Architecture**: See `PROJECT_OVERVIEW.md`
- **General Info**: See `README.md`

## Next Steps

1. ✅ Install dependencies: `npm install`
2. ✅ Start dev server: `npm run dev`
3. ✅ Open browser: `http://localhost:3000`
4. ✅ Login with Spotify
5. ✅ Run library analysis
6. ✅ Explore genres
7. ✅ Generate sort plan
8. ✅ Preview changes
9. ✅ Execute (toggle dry-run OFF first!)

## Common Workflows

### Developing a New Feature
1. Create component in appropriate directory
2. Import and use in page component
3. Update types in `src/lib/types.ts` if needed
4. Test in browser with hot reload
5. Build to verify: `npm run build`

### Adding a New Page
1. Create page component in `src/pages/`
2. Add route in `src/App.tsx`
3. Add navigation link in `src/components/layout/Sidebar.tsx`
4. Wrap in `ProtectedRoute` if auth required

### Calling a New API Endpoint
1. Add method to `src/lib/api.ts`
2. Add type to `src/lib/types.ts`
3. Create hook in `src/hooks/useApi.ts` (optional)
4. Use in component

## Performance Tips

- React Query caches API responses automatically
- SSE connections are managed by the useSSE hook
- Zustand state updates only re-render affected components
- Production build is optimized and minified

## Security Notes

- Authentication uses HTTP-only cookies (set by backend)
- No tokens stored in localStorage
- CORS configured via backend
- OAuth flow handled securely by backend

## Development Tips

1. **Use TypeScript**: Leverage type safety - check `src/lib/types.ts`
2. **Check Console**: Development errors appear in browser console
3. **Hot Reload**: Changes auto-reload, no need to refresh
4. **Component Dev**: Test components individually before integrating
5. **API Testing**: Use browser DevTools Network tab to debug API calls

## Production Deployment

1. Build the application:
   ```bash
   npm run build
   ```

2. Test the build locally:
   ```bash
   npm run preview
   ```

3. Deploy the `dist/` folder to your hosting service

4. Configure environment:
   - Set API base URL if not using proxy
   - Update OAuth redirect URLs in backend config

Enjoy building with Spotify Playlist Sorter! 🎵
