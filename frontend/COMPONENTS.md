# Component Reference

## UI Components

### Button
```tsx
<Button
  variant="primary" | "secondary" | "danger" | "ghost"
  size="sm" | "md" | "lg"
  isLoading={boolean}
  onClick={handler}
>
  Click me
</Button>
```

### Card
```tsx
<Card variant="default" | "highlighted">
  <CardHeader>
    <CardTitle>Title</CardTitle>
  </CardHeader>
  <CardContent>
    Content here
  </CardContent>
</Card>
```

### Badge
```tsx
<Badge variant="default" | "success" | "warning" | "danger" | "info">
  Label
</Badge>
```

### Progress
```tsx
<Progress
  value={50}
  max={100}
  showLabel={true}
/>
```

### Toggle
```tsx
<Toggle
  checked={isDryRun}
  onChange={toggleDryRun}
  label="Dry Run"
/>
```

## Layout Components

### Header
Global header with user info and dry-run toggle. Always visible when authenticated.

### Sidebar
Navigation sidebar with routes to Dashboard, Genres, and Changes pages.

## Feature Components

### SongCard
```tsx
<SongCard track={trackObject} />
```

### SongList
```tsx
<SongList
  tracks={trackArray}
  emptyMessage="No tracks found"
/>
```

### GenreCard
```tsx
<GenreCard genre={genreStatObject} />
```

### GenreGrid
```tsx
<GenreGrid genres={genreStatsArray} />
```

### TrackMove
```tsx
<TrackMove
  move={trackMoveObject}
  variant="add" | "remove"
/>
```

### ChangeSummary
```tsx
<ChangeSummary
  tracksToAdd={number}
  tracksToRemove={number}
  playlistsToCreate={number}
/>
```

### ChangeDiff
```tsx
<ChangeDiff
  tracksToAdd={trackMoveArray}
  tracksToRemove={trackMoveArray}
/>
```

### ProgressBar
```tsx
<ProgressBar
  current={50}
  total={100}
  message="Processing..."
/>
```

### LiveLog
```tsx
<LiveLog events={progressEventsArray} />
```

## Hooks

### useAuth
```tsx
const {
  user,
  isLoading,
  isAuthenticated,
  login,
  logout,
  checkAuth
} = useAuth();
```

### useSSE
```tsx
const {
  events,
  isConnected,
  error,
  clear
} = useSSE(endpoint, {
  onMessage: (event) => {},
  onError: (error) => {},
  onComplete: () => {}
});
```

### useApi
```tsx
// Library stats
const { data, isLoading, error } = useLibraryStats();

// Create sort plan
const createPlan = useCreateSortPlan();
await createPlan.mutateAsync();

// Execute sort plan
const executePlan = useExecuteSortPlan();
await executePlan.mutateAsync(planId);
```

## Stores

### useAuthStore
```tsx
const {
  user,
  isLoading,
  isAuthenticated,
  setUser,
  setLoading,
  logout
} = useAuthStore();
```

### useUIStore
```tsx
const {
  isDryRun,
  toggleDryRun,
  setDryRun
} = useUIStore();
```

## Type Definitions

All types are defined in `src/lib/types.ts`:

- `Artist`
- `Track`
- `TrackMove`
- `GenreStat`
- `SortPlan`
- `User`
- `LibraryStats`
- `ProgressEvent`
- `AuthResponse`

## Styling Utilities

### Tailwind Custom Classes

```css
/* Colors */
bg-spotify-green      /* #1DB954 - Spotify brand green */
bg-spotify-black      /* #121212 - Background */
bg-spotify-darkgray   /* #181818 */
bg-spotify-gray       /* #282828 */
text-spotify-lightgray /* #B3B3B3 */

/* Responsive breakpoints */
sm:  /* 640px */
md:  /* 768px */
lg:  /* 1024px */
xl:  /* 1280px */
```

## Routing

```tsx
/ - Dashboard
/genres - Genres view
/changes - Changes/Sort Plan view
/login - Login page
/callback - OAuth callback handler
```

Protected routes require authentication. Non-authenticated users are redirected to `/login`.
