package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/domain"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/genre"
	spotifyClient "github.com/adelvecchio/spotify-playlist-sorter/internal/spotify"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/sse"
)

// LibraryService handles fetching and analyzing user's Spotify library
type LibraryService struct {
	spotifyClient *spotifyClient.Client
	broadcaster   *sse.Broadcaster
}

// NewLibraryService creates a new library service
func NewLibraryService(client *spotifyClient.Client, broadcaster *sse.Broadcaster) *LibraryService {
	return &LibraryService{
		spotifyClient: client,
		broadcaster:   broadcaster,
	}
}

// LibraryAnalysis contains the complete analysis of user's library
type LibraryAnalysis struct {
	Tracks             []domain.Track            `json:"tracks"`
	Playlists          []domain.Playlist         `json:"playlists"`
	GenreDistribution  map[string]int            `json:"genreDistribution"`
	TotalLikedSongs    int                       `json:"totalLikedSongs"`
	TracksWithGenre    int                       `json:"tracksWithGenre"`
	TracksWithoutGenre int                       `json:"tracksWithoutGenre"`
	GroupingSuggestions []genre.GroupSuggestion  `json:"groupingSuggestions"`
	GenreGroups        map[string]*genre.GenreGroup `json:"genreGroups"`
}

// AnalyzeLibrary fetches all liked songs, playlists, and analyzes genres
func (s *LibraryService) AnalyzeLibrary(ctx context.Context, client *spotify.Client, userID string) (*LibraryAnalysis, error) {
	log.Info().Str("userID", userID).Msg("Starting library analysis")

	// Fetch liked songs
	s.broadcaster.SendProgress(userID, sse.PhaseFetchingLikedSongs, 0, 0, "Fetching your liked songs...")
	tracks, err := s.spotifyClient.FetchAllLikedSongs(ctx, client, func(current, total int) {
		s.broadcaster.SendProgress(userID, sse.PhaseFetchingLikedSongs, current, total,
			fmt.Sprintf("Fetching liked songs: %d/%d", current, total))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch liked songs: %w", err)
	}

	log.Info().Int("count", len(tracks)).Msg("Fetched liked songs")

	// Fetch playlists
	s.broadcaster.SendProgress(userID, sse.PhaseFetchingPlaylists, 0, 0, "Fetching your playlists...")
	playlists, err := s.spotifyClient.FetchAllPlaylists(ctx, client, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}

	log.Info().Int("count", len(playlists)).Msg("Fetched playlists")

	// Fetch playlist tracks for managed playlists
	s.broadcaster.SendInfo(userID, "Loading managed playlists...")
	for i := range playlists {
		if playlists[i].ManagedByApp && playlists[i].OwnerID == userID {
			trackIDs, err := s.spotifyClient.FetchPlaylistTracks(ctx, client, playlists[i].ID)
			if err != nil {
				log.Warn().Err(err).Str("playlistID", playlists[i].ID).Msg("Failed to fetch playlist tracks")
				continue
			}
			playlists[i].TrackIDs = trackIDs
		}
	}

	// Build track to playlist mapping
	trackToPlaylists := make(map[string][]string)
	for _, playlist := range playlists {
		if playlist.ManagedByApp && playlist.OwnerID == userID {
			for _, trackID := range playlist.TrackIDs {
				trackToPlaylists[trackID] = append(trackToPlaylists[trackID], playlist.ID)
			}
		}
	}

	// Update tracks with playlist membership
	for i := range tracks {
		if playlists, ok := trackToPlaylists[tracks[i].ID]; ok {
			tracks[i].InPlaylists = playlists
		}
	}

	// Fetch artist genres
	s.broadcaster.SendProgress(userID, sse.PhaseFetchingArtists, 0, len(tracks), "Fetching artist information...")
	tracks, err = s.enrichTracksWithGenres(ctx, client, tracks, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich tracks with genres: %w", err)
	}

	// Analyze genre distribution
	s.broadcaster.SendProgress(userID, sse.PhaseAnalyzing, 0, 0, "Analyzing your music library...")
	genreDistribution := make(map[string]int)
	tracksWithGenre := 0
	tracksWithoutGenre := 0

	for _, track := range tracks {
		if track.PrimaryGenre != "" {
			genreDistribution[track.PrimaryGenre]++
			tracksWithGenre++
		} else {
			tracksWithoutGenre++
		}
	}

	// Generate grouping suggestions (min 10 tracks per genre to suggest grouping)
	groupingSuggestions := genre.SuggestGroupings(genreDistribution, 10)
	genreGroups := genre.GroupGenres(genreDistribution)

	log.Info().
		Int("total", len(tracks)).
		Int("withGenre", tracksWithGenre).
		Int("withoutGenre", tracksWithoutGenre).
		Int("uniqueGenres", len(genreDistribution)).
		Int("groupingSuggestions", len(groupingSuggestions)).
		Msg("Library analysis complete")

	return &LibraryAnalysis{
		Tracks:              tracks,
		Playlists:           playlists,
		GenreDistribution:   genreDistribution,
		TotalLikedSongs:     len(tracks),
		TracksWithGenre:     tracksWithGenre,
		TracksWithoutGenre:  tracksWithoutGenre,
		GroupingSuggestions: groupingSuggestions,
		GenreGroups:         genreGroups,
	}, nil
}

// enrichTracksWithGenres fetches artist information and assigns genres to tracks
func (s *LibraryService) enrichTracksWithGenres(ctx context.Context, client *spotify.Client, tracks []domain.Track, userID string) ([]domain.Track, error) {
	// Collect unique artist IDs
	artistIDMap := make(map[string]bool)
	for _, track := range tracks {
		for _, artist := range track.Artists {
			artistIDMap[artist.ID] = true
		}
	}

	// Convert to slice
	artistIDs := make([]spotify.ID, 0, len(artistIDMap))
	for id := range artistIDMap {
		artistIDs = append(artistIDs, spotify.ID(id))
	}

	log.Info().Int("count", len(artistIDs)).Msg("Fetching artist genres")

	// Batch fetch artists
	artistsMap, err := s.spotifyClient.BatchFetchArtists(ctx, client, artistIDs)
	if err != nil {
		return nil, err
	}

	// Update tracks with artist genres
	for i := range tracks {
		for j := range tracks[i].Artists {
			artistID := tracks[i].Artists[j].ID
			if artist, ok := artistsMap[artistID]; ok && artist != nil {
				tracks[i].Artists[j].Genres = artist.Genres
			}
		}

		// Assign primary genre
		tracks[i].PrimaryGenre = s.determinePrimaryGenre(tracks[i])

		// Update progress periodically
		if (i+1)%100 == 0 || i == len(tracks)-1 {
			s.broadcaster.SendProgress(userID, sse.PhaseFetchingArtists, i+1, len(tracks),
				fmt.Sprintf("Processing artist genres: %d/%d", i+1, len(tracks)))
		}
	}

	return tracks, nil
}

// determinePrimaryGenre determines the primary genre for a track based on its artists
func (s *LibraryService) determinePrimaryGenre(track domain.Track) string {
	// Collect all genres from all artists
	var allGenres []string
	genreCounts := make(map[string]int)

	for _, artist := range track.Artists {
		for _, g := range artist.Genres {
			normalized := genre.NormalizeGenre(g)
			if normalized != "" {
				allGenres = append(allGenres, g)
				genreCounts[normalized]++
			}
		}
	}

	if len(allGenres) == 0 {
		return ""
	}

	// Find the most common genre
	maxCount := 0
	var primaryGenre string

	for g, count := range genreCounts {
		if count > maxCount {
			maxCount = count
			// Find original (non-normalized) genre name
			for _, original := range allGenres {
				if genre.NormalizeGenre(original) == g {
					primaryGenre = original
					break
				}
			}
		}
	}

	// If still no primary genre, use the extraction logic
	if primaryGenre == "" {
		primaryGenre = genre.ExtractPrimaryGenre(allGenres)
	}

	return primaryGenre
}

// GetManagedPlaylists returns only playlists managed by the app
func (s *LibraryService) GetManagedPlaylists(playlists []domain.Playlist, userID string) []domain.Playlist {
	var managed []domain.Playlist
	for _, p := range playlists {
		if p.ManagedByApp && p.OwnerID == userID {
			managed = append(managed, p)
		}
	}
	return managed
}

// BuildGenreToPlaylistMap creates a mapping from normalized genre to playlist
// When enabledGroups is provided, the map will be used with effective genres (after grouping)
func (s *LibraryService) BuildGenreToPlaylistMap(playlists []domain.Playlist, userID string, enabledGroups map[string]bool) map[string]*domain.Playlist {
	result := make(map[string]*domain.Playlist)

	for i := range playlists {
		if playlists[i].ManagedByApp && playlists[i].OwnerID == userID {
			playlistGenre := playlists[i].AssignedGenre
			normalized := genre.NormalizeGenre(playlistGenre)
			if normalized != "" {
				// Map the playlist's genre directly
				result[normalized] = &playlists[i]
			}
		}
	}

	return result
}
