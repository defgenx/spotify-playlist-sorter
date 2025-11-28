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

// ExecutorService executes sort plans
type ExecutorService struct {
	spotifyClient  *spotifyClient.Client
	libraryService *LibraryService
	broadcaster    *sse.Broadcaster
}

// NewExecutorService creates a new executor service
func NewExecutorService(client *spotifyClient.Client, libraryService *LibraryService, broadcaster *sse.Broadcaster) *ExecutorService {
	return &ExecutorService{
		spotifyClient:  client,
		libraryService: libraryService,
		broadcaster:    broadcaster,
	}
}

// ExecuteSortPlan executes a sort plan
func (s *ExecutorService) ExecuteSortPlan(ctx context.Context, client *spotify.Client, plan *domain.SortPlan, userID string) (*domain.ExecutionResult, error) {
	log.Info().Str("planID", plan.ID).Str("userID", userID).Msg("Executing sort plan")

	result := &domain.ExecutionResult{
		Success:          true,
		PlaylistsCreated: 0,
		PlaylistsDeleted: 0,
		TracksAdded:      0,
		TracksRemoved:    0,
		Errors:           []domain.ExecutionError{},
	}

	if plan.DryRun {
		s.broadcaster.SendInfo(userID, "Dry run mode - no changes will be made")
		return result, nil
	}

	// Step 1: Create new playlists
	if len(plan.PlaylistsToCreate) > 0 {
		s.broadcaster.SendProgress(userID, sse.PhaseCreatingPlaylists, 0, len(plan.PlaylistsToCreate), "Creating new playlists...")

		createdPlaylists, err := s.createPlaylists(ctx, client, plan.PlaylistsToCreate, userID)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, domain.ExecutionError{
				Operation: "create_playlists",
				Error:     err.Error(),
			})
			return result, err
		}

		result.PlaylistsCreated = len(createdPlaylists)

		// Update plan with created playlist IDs
		s.updatePlanWithCreatedPlaylists(plan, createdPlaylists)
	}

	// Step 2: Add tracks to playlists
	if len(plan.TracksToAdd) > 0 {
		s.broadcaster.SendProgress(userID, sse.PhaseAddingTracks, 0, len(plan.TracksToAdd), "Adding tracks to playlists...")

		added, errors := s.addTracksToPlaylists(ctx, client, plan.TracksToAdd, userID)
		result.TracksAdded = added
		result.Errors = append(result.Errors, errors...)
	}

	// Step 3: Handle uncategorized tracks
	if len(plan.UncategorizedTracks) > 0 {
		s.broadcaster.SendInfo(userID, fmt.Sprintf("Processing %d uncategorized tracks...", len(plan.UncategorizedTracks)))

		added, errors := s.handleUncategorizedTracks(ctx, client, plan.UncategorizedTracks, userID)
		result.TracksAdded += added
		result.Errors = append(result.Errors, errors...)
	}

	// Step 4: Remove tracks from wrong playlists
	if len(plan.TracksToRemove) > 0 {
		s.broadcaster.SendProgress(userID, sse.PhaseRemovingTracks, 0, len(plan.TracksToRemove), "Removing tracks from incorrect playlists...")

		removed, errors := s.removeTracksFromPlaylists(ctx, client, plan.TracksToRemove, userID)
		result.TracksRemoved = removed
		result.Errors = append(result.Errors, errors...)
	}

	// Step 5: Remove empty playlists
	s.broadcaster.SendInfo(userID, "Checking for empty playlists...")
	deleted, errors := s.removeEmptyPlaylists(ctx, client, userID)
	result.PlaylistsDeleted = deleted
	result.Errors = append(result.Errors, errors...)

	if len(result.Errors) > 0 {
		result.Success = false
	}

	s.broadcaster.SendComplete(userID, fmt.Sprintf("Sort complete! Created %d playlists, deleted %d empty playlists, added %d tracks, removed %d tracks",
		result.PlaylistsCreated, result.PlaylistsDeleted, result.TracksAdded, result.TracksRemoved))

	log.Info().
		Int("playlistsCreated", result.PlaylistsCreated).
		Int("playlistsDeleted", result.PlaylistsDeleted).
		Int("tracksAdded", result.TracksAdded).
		Int("tracksRemoved", result.TracksRemoved).
		Int("errors", len(result.Errors)).
		Msg("Sort plan execution complete")

	return result, nil
}

// createPlaylists creates new playlists for genres
func (s *ExecutorService) createPlaylists(ctx context.Context, client *spotify.Client, genres []string, userID string) (map[string]string, error) {
	createdPlaylists := make(map[string]string) // genre -> playlistID

	for i, genreName := range genres {
		s.broadcaster.SendProgress(userID, sse.PhaseCreatingPlaylists, i+1, len(genres),
			fmt.Sprintf("Creating playlist for %s...", genreName))

		description := fmt.Sprintf("Automatically organized %s tracks", genreName)
		playlist, err := s.spotifyClient.CreatePlaylist(ctx, client, userID, genreName, description, false)
		if err != nil {
			log.Error().Err(err).Str("genre", genreName).Msg("Failed to create playlist")
			return createdPlaylists, fmt.Errorf("failed to create playlist for %s: %w", genreName, err)
		}

		createdPlaylists[genreName] = playlist.ID.String()
		log.Info().Str("genre", genreName).Str("playlistID", playlist.ID.String()).Msg("Created playlist")
	}

	return createdPlaylists, nil
}

// updatePlanWithCreatedPlaylists updates the plan with newly created playlist IDs
func (s *ExecutorService) updatePlanWithCreatedPlaylists(plan *domain.SortPlan, createdPlaylists map[string]string) {
	// Update TracksToAdd with correct playlist IDs
	// Use ToPlaylistName which contains the effectiveGenre (after grouping), not Genre which is the original
	for i := range plan.TracksToAdd {
		if plan.TracksToAdd[i].ToPlaylist == "" {
			// This track was going to a new playlist
			// Match by ToPlaylistName which is the effectiveGenre
			if playlistID, ok := createdPlaylists[plan.TracksToAdd[i].ToPlaylistName]; ok {
				plan.TracksToAdd[i].ToPlaylist = playlistID
			}
		}
	}

	// Update GenreStats
	// Note: GenreStats uses original genre names, but we need to match by effectiveGenre
	// For now, we'll match by genre name directly (this may need refinement if grouping affects stats)
	for i := range plan.GenreStats {
		if plan.GenreStats[i].IsNew {
			// Try matching by original genre first
			if playlistID, ok := createdPlaylists[plan.GenreStats[i].Genre]; ok {
				plan.GenreStats[i].PlaylistID = playlistID
			}
		}
	}
}

// addTracksToPlaylists adds tracks to their target playlists
func (s *ExecutorService) addTracksToPlaylists(ctx context.Context, client *spotify.Client, moves []domain.TrackMove, userID string) (int, []domain.ExecutionError) {
	// Group tracks by target playlist
	playlistTracks := make(map[string][]spotify.ID)
	for _, move := range moves {
		if move.ToPlaylist != "" {
			playlistTracks[move.ToPlaylist] = append(playlistTracks[move.ToPlaylist], spotify.ID(move.TrackID))
		}
	}

	totalAdded := 0
	var errors []domain.ExecutionError
	current := 0
	total := len(moves)

	for playlistID, trackIDs := range playlistTracks {
		s.broadcaster.SendProgress(userID, sse.PhaseAddingTracks, current, total,
			fmt.Sprintf("Adding %d tracks to playlist...", len(trackIDs)))

		err := s.spotifyClient.AddTracksToPlaylist(ctx, client, playlistID, trackIDs)
		if err != nil {
			log.Error().Err(err).Str("playlistID", playlistID).Msg("Failed to add tracks to playlist")
			errors = append(errors, domain.ExecutionError{
				Operation: "add_tracks",
				Playlist:  playlistID,
				Error:     err.Error(),
			})
		} else {
			totalAdded += len(trackIDs)
		}

		current += len(trackIDs)
	}

	return totalAdded, errors
}

// removeTracksFromPlaylists removes tracks from playlists
func (s *ExecutorService) removeTracksFromPlaylists(ctx context.Context, client *spotify.Client, moves []domain.TrackMove, userID string) (int, []domain.ExecutionError) {
	// Group tracks by source playlist
	playlistTracks := make(map[string][]spotify.ID)
	for _, move := range moves {
		if move.FromPlaylist != "" {
			playlistTracks[move.FromPlaylist] = append(playlistTracks[move.FromPlaylist], spotify.ID(move.TrackID))
		}
	}

	totalRemoved := 0
	var errors []domain.ExecutionError
	current := 0
	total := len(moves)

	for playlistID, trackIDs := range playlistTracks {
		s.broadcaster.SendProgress(userID, sse.PhaseRemovingTracks, current, total,
			fmt.Sprintf("Removing %d tracks from playlist...", len(trackIDs)))

		// Remove in batches to avoid API limits
		batchSize := 100
		for i := 0; i < len(trackIDs); i += batchSize {
			end := i + batchSize
			if end > len(trackIDs) {
				end = len(trackIDs)
			}
			batch := trackIDs[i:end]

			err := s.spotifyClient.RemoveTracksFromPlaylist(ctx, client, playlistID, batch)
			if err != nil {
				log.Error().Err(err).Str("playlistID", playlistID).Msg("Failed to remove tracks from playlist")
				errors = append(errors, domain.ExecutionError{
					Operation: "remove_tracks",
					Playlist:  playlistID,
					Error:     err.Error(),
				})
			} else {
				totalRemoved += len(batch)
			}
		}

		current += len(trackIDs)
	}

	return totalRemoved, errors
}

// handleUncategorizedTracks creates/updates an "Uncategorized" playlist
func (s *ExecutorService) handleUncategorizedTracks(ctx context.Context, client *spotify.Client, tracks []domain.Track, userID string) (int, []domain.ExecutionError) {
	if len(tracks) == 0 {
		return 0, nil
	}

	var errors []domain.ExecutionError

	// Check if "Uncategorized" playlist already exists
	playlists, err := s.spotifyClient.FetchAllPlaylists(ctx, client, userID)
	if err != nil {
		errors = append(errors, domain.ExecutionError{
			Operation: "fetch_playlists",
			Error:     err.Error(),
		})
		return 0, errors
	}

	var uncategorizedPlaylist *domain.Playlist
	for i := range playlists {
		if playlists[i].ManagedByApp && playlists[i].OwnerID == userID {
			normalized := genre.NormalizeGenre(playlists[i].AssignedGenre)
			if normalized == "uncategorized" || playlists[i].Name == "Uncategorized" {
				uncategorizedPlaylist = &playlists[i]
				break
			}
		}
	}

	// Create if doesn't exist
	if uncategorizedPlaylist == nil {
		s.broadcaster.SendInfo(userID, "Creating Uncategorized playlist...")
		playlist, err := s.spotifyClient.CreatePlaylist(ctx, client, userID, "Uncategorized", "Songs without a clear genre", false)
		if err != nil {
			errors = append(errors, domain.ExecutionError{
				Operation: "create_uncategorized_playlist",
				Error:     err.Error(),
			})
			return 0, errors
		}
		uncategorizedPlaylist = &domain.Playlist{
			ID:   playlist.ID.String(),
			Name: playlist.Name,
		}
	}

	// Add tracks
	trackIDs := make([]spotify.ID, len(tracks))
	for i, track := range tracks {
		trackIDs[i] = spotify.ID(track.ID)
	}

	s.broadcaster.SendInfo(userID, fmt.Sprintf("Adding %d tracks to Uncategorized playlist...", len(trackIDs)))
	err = s.spotifyClient.AddTracksToPlaylist(ctx, client, uncategorizedPlaylist.ID, trackIDs)
	if err != nil {
		errors = append(errors, domain.ExecutionError{
			Operation: "add_uncategorized_tracks",
			Playlist:  uncategorizedPlaylist.ID,
			Error:     err.Error(),
		})
		return 0, errors
	}

	return len(trackIDs), errors
}

// removeEmptyPlaylists finds and deletes empty managed playlists
func (s *ExecutorService) removeEmptyPlaylists(ctx context.Context, client *spotify.Client, userID string) (int, []domain.ExecutionError) {
	var errors []domain.ExecutionError

	// Fetch all playlists
	playlists, err := s.spotifyClient.FetchAllPlaylists(ctx, client, userID)
	if err != nil {
		errors = append(errors, domain.ExecutionError{
			Operation: "fetch_playlists_for_cleanup",
			Error:     err.Error(),
		})
		return 0, errors
	}

	deletedCount := 0
	for _, playlist := range playlists {
		// Only process managed playlists owned by the user
		if !playlist.ManagedByApp || playlist.OwnerID != userID {
			continue
		}

		// Skip "Uncategorized" playlist - we don't want to delete it even if empty
		normalized := genre.NormalizeGenre(playlist.AssignedGenre)
		if normalized == "uncategorized" || playlist.Name == "Uncategorized" {
			continue
		}

		// Check if playlist is empty
		// Use TrackCount from the playlist object (from Spotify API)
		if playlist.TrackCount == 0 {
			s.broadcaster.SendInfo(userID, fmt.Sprintf("Deleting empty playlist: %s", playlist.Name))
			
			err := s.spotifyClient.DeletePlaylist(ctx, client, playlist.ID)
			if err != nil {
				log.Error().Err(err).Str("playlistID", playlist.ID).Str("playlistName", playlist.Name).Msg("Failed to delete empty playlist")
				errors = append(errors, domain.ExecutionError{
					Operation: "delete_empty_playlist",
					Playlist:  playlist.ID,
					Error:     err.Error(),
				})
			} else {
				deletedCount++
				log.Info().Str("playlistID", playlist.ID).Str("playlistName", playlist.Name).Msg("Deleted empty playlist")
			}
		}
	}

	return deletedCount, errors
}
