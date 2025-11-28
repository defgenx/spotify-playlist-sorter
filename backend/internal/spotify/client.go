package spotify

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/domain"
)

type Client struct {
	auth        *spotifyauth.Authenticator
	rateLimiter *rate.Limiter
	mu          sync.Mutex
	clientID    string
	clientSecret string
}

func NewClient(clientID, clientSecret, redirectURL string) *Client {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
		),
	)

	return &Client{
		auth:         auth,
		rateLimiter:  rate.NewLimiter(rate.Limit(2), 5), // 2 requests/sec with burst of 5
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *Client) GetAuthURL(state string) string {
	return c.auth.AuthURL(state)
}

func (c *Client) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.auth.Exchange(ctx, code)
}

func (c *Client) NewSpotifyClient(ctx context.Context, token *oauth2.Token) *spotify.Client {
	httpClient := c.auth.Client(ctx, token)
	return spotify.New(httpClient)
}

func (c *Client) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	// Create OAuth2 config manually
	cfg := &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyauth.AuthURL,
			TokenURL: spotifyauth.TokenURL,
		},
	}
	return cfg.TokenSource(ctx, token)
}

// Rate limited wrapper for API calls
func (c *Client) withRateLimit(ctx context.Context) error {
	return c.rateLimiter.Wait(ctx)
}

// FetchAllLikedSongs fetches all liked songs with pagination
func (c *Client) FetchAllLikedSongs(ctx context.Context, client *spotify.Client, progressFn func(current, total int)) ([]domain.Track, error) {
	var allTracks []domain.Track
	limit := 50
	offset := 0

	// First request to get total
	if err := c.withRateLimit(ctx); err != nil {
		return nil, err
	}

	page, err := client.CurrentUsersTracks(ctx, spotify.Limit(limit), spotify.Offset(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch liked songs: %w", err)
	}

	total := int(page.Total)

	for _, item := range page.Tracks {
		allTracks = append(allTracks, convertSavedTrack(item))
	}

	if progressFn != nil {
		progressFn(len(allTracks), total)
	}

	// Fetch remaining pages
	for offset = limit; offset < total; offset += limit {
		if err := c.withRateLimit(ctx); err != nil {
			return nil, err
		}

		page, err := client.CurrentUsersTracks(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch liked songs at offset %d: %w", offset, err)
		}

		for _, item := range page.Tracks {
			allTracks = append(allTracks, convertSavedTrack(item))
		}

		if progressFn != nil {
			progressFn(len(allTracks), total)
		}
	}

	return allTracks, nil
}

// FetchAllPlaylists fetches all user playlists with pagination
func (c *Client) FetchAllPlaylists(ctx context.Context, client *spotify.Client, userID string) ([]domain.Playlist, error) {
	var allPlaylists []domain.Playlist
	limit := 50
	offset := 0

	for {
		if err := c.withRateLimit(ctx); err != nil {
			return nil, err
		}

		page, err := client.CurrentUsersPlaylists(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlists: %w", err)
		}

		for _, p := range page.Playlists {
			playlist := domain.Playlist{
				ID:          p.ID.String(),
				Name:        p.Name,
				Description: p.Description,
				OwnerID:     p.Owner.ID,
				TrackCount:  int(p.Tracks.Total),
			}

			if len(p.Images) > 0 {
				playlist.ImageURL = p.Images[0].URL
			}

			// Check if managed by our app
			if strings.Contains(p.Description, domain.ManagedTag) {
				playlist.ManagedByApp = true
				// Extract genre from name (assumes format "Genre Name")
				playlist.AssignedGenre = extractGenreFromName(p.Name)
			}

			allPlaylists = append(allPlaylists, playlist)
		}

		if page.Next == "" {
			break
		}
		offset += limit
	}

	return allPlaylists, nil
}

// FetchPlaylistTracks fetches all tracks from a playlist
func (c *Client) FetchPlaylistTracks(ctx context.Context, client *spotify.Client, playlistID string) ([]string, error) {
	var trackIDs []string
	limit := 50
	offset := 0

	for {
		if err := c.withRateLimit(ctx); err != nil {
			return nil, err
		}

		page, err := client.GetPlaylistItems(ctx, spotify.ID(playlistID), spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlist tracks: %w", err)
		}

		for _, item := range page.Items {
			if item.Track.Track != nil {
				trackIDs = append(trackIDs, item.Track.Track.ID.String())
			}
		}

		if page.Next == "" {
			break
		}
		offset += limit
	}

	return trackIDs, nil
}

// BatchFetchArtists fetches artists in batches of 50
func (c *Client) BatchFetchArtists(ctx context.Context, client *spotify.Client, artistIDs []spotify.ID) (map[string]*spotify.FullArtist, error) {
	result := make(map[string]*spotify.FullArtist)

	for i := 0; i < len(artistIDs); i += 50 {
		end := i + 50
		if end > len(artistIDs) {
			end = len(artistIDs)
		}
		batch := artistIDs[i:end]

		if err := c.withRateLimit(ctx); err != nil {
			return nil, err
		}

		artists, err := client.GetArtists(ctx, batch...)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch artists: %w", err)
		}

		for _, artist := range artists {
			if artist != nil {
				result[artist.ID.String()] = artist
			}
		}
	}

	return result, nil
}

// CreatePlaylist creates a new playlist
func (c *Client) CreatePlaylist(ctx context.Context, client *spotify.Client, userID, name, description string, public bool) (*spotify.FullPlaylist, error) {
	if err := c.withRateLimit(ctx); err != nil {
		return nil, err
	}

	fullDescription := description + " " + domain.ManagedTag
	return client.CreatePlaylistForUser(ctx, userID, name, fullDescription, public, false)
}

// AddTracksToPlaylist adds tracks in batches of 100
func (c *Client) AddTracksToPlaylist(ctx context.Context, client *spotify.Client, playlistID string, trackIDs []spotify.ID) error {
	for i := 0; i < len(trackIDs); i += 100 {
		end := i + 100
		if end > len(trackIDs) {
			end = len(trackIDs)
		}
		batch := trackIDs[i:end]

		if err := c.withRateLimit(ctx); err != nil {
			return err
		}

		_, err := client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), batch...)
		if err != nil {
			return fmt.Errorf("failed to add tracks to playlist: %w", err)
		}
	}

	return nil
}

// RemoveTracksFromPlaylist removes tracks from a playlist
func (c *Client) RemoveTracksFromPlaylist(ctx context.Context, client *spotify.Client, playlistID string, trackIDs []spotify.ID) error {
	if err := c.withRateLimit(ctx); err != nil {
		return err
	}

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), trackIDs...)
	if err != nil {
		return fmt.Errorf("failed to remove tracks from playlist: %w", err)
	}

	return nil
}

// GetCurrentUser returns the current user's profile
func (c *Client) GetCurrentUser(ctx context.Context, client *spotify.Client) (*spotify.PrivateUser, error) {
	if err := c.withRateLimit(ctx); err != nil {
		return nil, err
	}

	return client.CurrentUser(ctx)
}

// Helper functions

func convertSavedTrack(st spotify.SavedTrack) domain.Track {
	track := domain.Track{
		ID:       st.ID.String(),
		Name:     st.Name,
		Duration: int(st.Duration),
	}

	if st.Album.Name != "" {
		track.AlbumName = st.Album.Name
	}

	if len(st.Album.Images) > 0 {
		track.AlbumImage = st.Album.Images[0].URL
	}

	for _, a := range st.Artists {
		track.Artists = append(track.Artists, domain.Artist{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return track
}

func extractGenreFromName(name string) string {
	// For now, assume the playlist name IS the genre
	// Could add more sophisticated matching later
	return strings.ToLower(strings.TrimSpace(name))
}

// HandleRateLimitError checks if an error is a rate limit error and waits
func HandleRateLimitError(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}

	// Check for 429 status
	if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "rate limit") {
		// Default retry after 30 seconds
		log.Warn().Msg("Rate limited by Spotify API, waiting 30 seconds")
		return true, 30 * time.Second
	}

	return false, 0
}

// WithRetry wraps an operation with retry logic for rate limiting
func WithRetry(ctx context.Context, maxRetries int, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		isRateLimited, waitTime := HandleRateLimitError(err)
		if !isRateLimited {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			continue
		}
	}

	return lastErr
}
