package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/domain"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/genre"
)

// SorterService generates sort plans for organizing tracks into playlists
type SorterService struct {
	libraryService *LibraryService
}

// NewSorterService creates a new sorter service
func NewSorterService(libraryService *LibraryService) *SorterService {
	return &SorterService{
		libraryService: libraryService,
	}
}

// GenerateSortPlan creates a sort plan based on library analysis
func (s *SorterService) GenerateSortPlan(ctx context.Context, analysis *LibraryAnalysis, userID string, dryRun bool, enabledGroups map[string]bool) (*domain.SortPlan, error) {
	log.Info().Str("userID", userID).Bool("dryRun", dryRun).Int("enabledGroups", len(enabledGroups)).Msg("Generating sort plan")

	plan := &domain.SortPlan{
		ID:                  uuid.New().String(),
		CreatedAt:           time.Now(),
		DryRun:              dryRun,
		TotalLikedTracks:    len(analysis.Tracks),
		TracksToAdd:         []domain.TrackMove{},
		TracksToRemove:      []domain.TrackMove{},
		PlaylistsToCreate:   []string{},
		UncategorizedTracks: []domain.Track{},
		GenreStats:          []domain.GenreStat{},
		EnabledGroups:       enabledGroups,
	}

	// Build genre to playlist mapping (with grouping awareness)
	genreToPlaylist := s.libraryService.BuildGenreToPlaylistMap(analysis.Playlists, userID, enabledGroups)

	// Track which genres need new playlists
	neededGenres := make(map[string]bool)

	// Build playlist name lookup
	playlistNames := make(map[string]string)
	for _, p := range analysis.Playlists {
		playlistNames[p.ID] = p.Name
	}

	// Process each track
	for _, track := range analysis.Tracks {
		if track.PrimaryGenre == "" {
			// No genre found - goes to uncategorized
			plan.UncategorizedTracks = append(plan.UncategorizedTracks, track)
			continue
		}

		// Apply grouping if enabled for this genre's parent
		effectiveGenre := genre.ApplyGrouping(track.PrimaryGenre, enabledGroups)
		normalizedGenre := genre.NormalizeGenre(effectiveGenre)
		targetPlaylist, exists := genreToPlaylist[normalizedGenre]

		if !exists {
			// Need to create new playlist for this genre (use effective genre, not original)
			neededGenres[effectiveGenre] = true
		}

		// Check if track is already in the correct playlist
		inCorrectPlaylist := false
		if exists {
			for _, playlistID := range track.InPlaylists {
				if playlistID == targetPlaylist.ID {
					inCorrectPlaylist = true
					break
				}
			}
		}

		// If not in correct playlist, need to add
		if !inCorrectPlaylist {
			artistName := ""
			if len(track.Artists) > 0 {
				artistName = track.Artists[0].Name
			}

			toPlaylistID := ""
			toPlaylistName := effectiveGenre
			if exists {
				toPlaylistID = targetPlaylist.ID
				toPlaylistName = targetPlaylist.Name
			}

			reason := "Song belongs to this genre"
			if effectiveGenre != track.PrimaryGenre {
				reason = fmt.Sprintf("Song genre '%s' grouped into '%s'", track.PrimaryGenre, effectiveGenre)
			}

			plan.TracksToAdd = append(plan.TracksToAdd, domain.TrackMove{
				TrackID:          track.ID,
				TrackName:        track.Name,
				ArtistName:       artistName,
				AlbumImage:       track.AlbumImage,
				Genre:            track.PrimaryGenre,
				FromPlaylist:     "",
				FromPlaylistName: "",
				ToPlaylist:       toPlaylistID,
				ToPlaylistName:   toPlaylistName,
				Reason:           reason,
			})
		}

		// Check if track is in wrong managed playlists
		for _, playlistID := range track.InPlaylists {
			// Find the playlist
			var playlist *domain.Playlist
			for i := range analysis.Playlists {
				if analysis.Playlists[i].ID == playlistID {
					playlist = &analysis.Playlists[i]
					break
				}
			}

			if playlist == nil || !playlist.ManagedByApp || playlist.OwnerID != userID {
				continue
			}

			// Apply grouping to both playlist and track genres for comparison
			playlistEffectiveGenre := genre.ApplyGrouping(playlist.AssignedGenre, enabledGroups)
			playlistGenreNorm := genre.NormalizeGenre(playlistEffectiveGenre)
			trackEffectiveGenre := genre.ApplyGrouping(track.PrimaryGenre, enabledGroups)
			trackGenreNorm := genre.NormalizeGenre(trackEffectiveGenre)

			// Check if this playlist is the correct target playlist for this track
			isTargetPlaylist := exists && targetPlaylist != nil && playlist.ID == targetPlaylist.ID

			// Remove track if it's in the wrong playlist (genres don't match after grouping)
			// and it's not already in the correct target playlist
			if playlistGenreNorm != trackGenreNorm && !isTargetPlaylist {
				// Track is in wrong playlist
				artistName := ""
				if len(track.Artists) > 0 {
					artistName = track.Artists[0].Name
				}

				plan.TracksToRemove = append(plan.TracksToRemove, domain.TrackMove{
					TrackID:          track.ID,
					TrackName:        track.Name,
					ArtistName:       artistName,
					AlbumImage:       track.AlbumImage,
					Genre:            track.PrimaryGenre,
					FromPlaylist:     playlist.ID,
					FromPlaylistName: playlist.Name,
					ToPlaylist:       "",
					ToPlaylistName:   "",
					Reason:           fmt.Sprintf("Song genre (%s) doesn't match playlist (%s)", trackEffectiveGenre, playlistEffectiveGenre),
				})
			}
		}
	}

	// Add needed genres to playlists to create
	for genreName := range neededGenres {
		plan.PlaylistsToCreate = append(plan.PlaylistsToCreate, genreName)
	}

	// Generate genre statistics
	genreCounts := make(map[string]int)
	for _, track := range analysis.Tracks {
		if track.PrimaryGenre != "" {
			genreCounts[track.PrimaryGenre]++
		}
	}

	for genreName, count := range genreCounts {
		normalizedGenre := genre.NormalizeGenre(genreName)
		playlist, exists := genreToPlaylist[normalizedGenre]

		stat := domain.GenreStat{
			Genre:      genreName,
			TrackCount: count,
			IsNew:      !exists,
		}

		if exists {
			stat.PlaylistID = playlist.ID
		}

		plan.GenreStats = append(plan.GenreStats, stat)
	}

	log.Info().
		Int("tracksToAdd", len(plan.TracksToAdd)).
		Int("tracksToRemove", len(plan.TracksToRemove)).
		Int("playlistsToCreate", len(plan.PlaylistsToCreate)).
		Int("uncategorized", len(plan.UncategorizedTracks)).
		Msg("Sort plan generated")

	return plan, nil
}

// ValidateSortPlan checks if a sort plan is valid
func (s *SorterService) ValidateSortPlan(plan *domain.SortPlan) error {
	if plan == nil {
		return fmt.Errorf("sort plan is nil")
	}

	if plan.ID == "" {
		return fmt.Errorf("sort plan ID is empty")
	}

	// Add more validation as needed
	return nil
}
