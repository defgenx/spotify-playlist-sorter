package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/api/middleware"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/service"
	spotifyClient "github.com/adelvecchio/spotify-playlist-sorter/internal/spotify"
)

// LibraryHandler handles library endpoints
type LibraryHandler struct {
	spotifyClient  *spotifyClient.Client
	libraryService *service.LibraryService
}

// NewLibraryHandler creates a new library handler
func NewLibraryHandler(spotifyClient *spotifyClient.Client, libraryService *service.LibraryService) *LibraryHandler {
	return &LibraryHandler{
		spotifyClient:  spotifyClient,
		libraryService: libraryService,
	}
}

// GetAnalysis analyzes the user's library
func (h *LibraryHandler) GetAnalysis(c *gin.Context) {
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
	log.Info().Str("userID", userID).Msg("Starting library analysis")
	analysis, err := h.libraryService.AnalyzeLibrary(ctx, client, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to analyze library")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze library: " + err.Error(),
		})
		return
	}

	log.Info().
		Str("userID", userID).
		Int("totalTracks", analysis.TotalLikedSongs).
		Int("playlists", len(analysis.Playlists)).
		Msg("Library analysis complete")

	c.JSON(http.StatusOK, analysis)
}
