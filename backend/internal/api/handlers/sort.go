package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/api/middleware"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/domain"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/genre"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/service"
	spotifyClient "github.com/adelvecchio/spotify-playlist-sorter/internal/spotify"
)

// SortHandler handles sort endpoints
type SortHandler struct {
	spotifyClient   *spotifyClient.Client
	libraryService  *service.LibraryService
	sorterService   *service.SorterService
	executorService *service.ExecutorService
}

// NewSortHandler creates a new sort handler
func NewSortHandler(
	spotifyClient *spotifyClient.Client,
	libraryService *service.LibraryService,
	sorterService *service.SorterService,
	executorService *service.ExecutorService,
) *SortHandler {
	return &SortHandler{
		spotifyClient:   spotifyClient,
		libraryService:  libraryService,
		sorterService:   sorterService,
		executorService: executorService,
	}
}

// GeneratePlanRequest represents a request to generate a sort plan
type GeneratePlanRequest struct {
	DryRun            bool     `json:"dryRun"`
	EnabledGroups     []string `json:"enabledGroups"`     // Parent genres to group (e.g., ["Rock", "Pop"])
	DisabledPlaylists []string `json:"disabledPlaylists"` // Genre names to skip creating playlists for
}

// GeneratePlanResponse includes the sort plan and grouping suggestions
type GeneratePlanResponse struct {
	*domain.SortPlan
	GroupingSuggestions []genre.GroupSuggestion `json:"groupingSuggestions"`
}

// GeneratePlan generates a sort plan
func (h *SortHandler) GeneratePlan(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	sess, exists := middleware.GetSession(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No session found",
		})
		return
	}

	var req GeneratePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.DryRun = true // Default to dry run
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create Spotify client with token refresh
	tokenSource := h.spotifyClient.TokenSource(ctx, sess.Token)
	token, err := tokenSource.Token()
	if err != nil {
		log.Error().Err(err).Msg("Failed to refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to refresh token",
		})
		return
	}

	client := h.spotifyClient.NewSpotifyClient(ctx, token)

	// Analyze library
	log.Info().Str("userID", userID).Msg("Analyzing library for sort plan")
	analysis, err := h.libraryService.AnalyzeLibrary(ctx, client, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to analyze library")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze library: " + err.Error(),
		})
		return
	}

	// Convert enabled groups to map for easier lookup
	enabledGroupsMap := make(map[string]bool)
	for _, g := range req.EnabledGroups {
		enabledGroupsMap[g] = true
	}

	// Convert disabled playlists to map for easier lookup
	disabledPlaylistsMap := make(map[string]bool)
	for _, p := range req.DisabledPlaylists {
		disabledPlaylistsMap[p] = true
	}

	// Generate sort plan
	log.Info().Str("userID", userID).Bool("dryRun", req.DryRun).Int("enabledGroups", len(req.EnabledGroups)).Int("disabledPlaylists", len(req.DisabledPlaylists)).Msg("Generating sort plan")
	plan, err := h.sorterService.GenerateSortPlan(ctx, analysis, userID, req.DryRun, enabledGroupsMap)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate sort plan")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate sort plan: " + err.Error(),
		})
		return
	}

	// Filter out disabled playlists from playlistsToCreate
	if len(disabledPlaylistsMap) > 0 {
		filteredPlaylists := []string{}
		for _, p := range plan.PlaylistsToCreate {
			if !disabledPlaylistsMap[p] {
				filteredPlaylists = append(filteredPlaylists, p)
			}
		}
		plan.PlaylistsToCreate = filteredPlaylists

		// Also filter out tracks that would go to disabled playlists
		filteredTracksToAdd := []domain.TrackMove{}
		for _, t := range plan.TracksToAdd {
			if !disabledPlaylistsMap[t.ToPlaylistName] {
				filteredTracksToAdd = append(filteredTracksToAdd, t)
			}
		}
		plan.TracksToAdd = filteredTracksToAdd
	}

	log.Info().
		Str("userID", userID).
		Str("planID", plan.ID).
		Int("tracksToAdd", len(plan.TracksToAdd)).
		Int("tracksToRemove", len(plan.TracksToRemove)).
		Int("playlistsToCreate", len(plan.PlaylistsToCreate)).
		Msg("Sort plan generated")

	response := GeneratePlanResponse{
		SortPlan:            plan,
		GroupingSuggestions: analysis.GroupingSuggestions,
	}

	c.JSON(http.StatusOK, response)
}

// ExecutePlanRequest represents a request to execute a sort plan
type ExecutePlanRequest struct {
	DryRun            bool     `json:"dryRun"`
	EnabledGroups     []string `json:"enabledGroups"`     // Parent genres to group (e.g., ["Rock", "Pop"])
	DisabledPlaylists []string `json:"disabledPlaylists"` // Genre names to skip creating playlists for
}

// ExecutePlan executes a sort plan
func (h *SortHandler) ExecutePlan(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	sess, exists := middleware.GetSession(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No session found",
		})
		return
	}

	var req ExecutePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.DryRun = false // Default to actual execution
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Create Spotify client with token refresh
	tokenSource := h.spotifyClient.TokenSource(ctx, sess.Token)
	token, err := tokenSource.Token()
	if err != nil {
		log.Error().Err(err).Msg("Failed to refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to refresh token",
		})
		return
	}

	client := h.spotifyClient.NewSpotifyClient(ctx, token)

	// Analyze library
	log.Info().Str("userID", userID).Msg("Analyzing library for execution")
	analysis, err := h.libraryService.AnalyzeLibrary(ctx, client, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to analyze library")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze library: " + err.Error(),
		})
		return
	}

	// Convert enabled groups to map for easier lookup
	enabledGroupsMap := make(map[string]bool)
	for _, g := range req.EnabledGroups {
		enabledGroupsMap[g] = true
	}

	// Generate sort plan
	log.Info().Str("userID", userID).Bool("dryRun", req.DryRun).Msg("Generating sort plan for execution")
	plan, err := h.sorterService.GenerateSortPlan(ctx, analysis, userID, req.DryRun, enabledGroupsMap)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate sort plan")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate sort plan: " + err.Error(),
		})
		return
	}

	// Execute plan
	log.Info().Str("userID", userID).Str("planID", plan.ID).Msg("Executing sort plan")
	result, err := h.executorService.ExecuteSortPlan(ctx, client, plan, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute sort plan")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute sort plan: " + err.Error(),
		})
		return
	}

	log.Info().
		Str("userID", userID).
		Str("planID", plan.ID).
		Bool("success", result.Success).
		Int("playlistsCreated", result.PlaylistsCreated).
		Int("tracksAdded", result.TracksAdded).
		Int("tracksRemoved", result.TracksRemoved).
		Msg("Sort plan executed")

	c.JSON(http.StatusOK, result)
}
