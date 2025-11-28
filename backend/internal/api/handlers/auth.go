package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/api/middleware"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/config"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/session"
	spotifyClient "github.com/adelvecchio/spotify-playlist-sorter/internal/spotify"
)

// stateStore holds OAuth states in memory
var (
	stateStore   = make(map[string]time.Time)
	stateStoreMu sync.Mutex
)

// tempTokenStore holds temporary tokens for completing login
var (
	tempTokenStore   = make(map[string]string) // tempToken -> sessionID
	tempTokenStoreMu sync.Mutex
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	spotifyClient *spotifyClient.Client
	sessionStore  *session.Store
	config        *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(spotifyClient *spotifyClient.Client, sessionStore *session.Store, config *config.Config) *AuthHandler {
	return &AuthHandler{
		spotifyClient: spotifyClient,
		sessionStore:  sessionStore,
		config:        config,
	}
}

// Login initiates the OAuth flow
func (h *AuthHandler) Login(c *gin.Context) {
	// Generate state for CSRF protection
	state := uuid.New().String()

	// Store state server-side
	stateStoreMu.Lock()
	stateStore[state] = time.Now().Add(10 * time.Minute)
	log.Info().Str("state", state).Int("storeSize", len(stateStore)).Msg("Stored OAuth state")
	stateStoreMu.Unlock()

	// Get auth URL
	authURL := h.spotifyClient.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"url": authURL,
	})
}

// Callback handles the OAuth callback
func (h *AuthHandler) Callback(c *gin.Context) {
	// Verify state from server-side store
	state := c.Query("state")

	stateStoreMu.Lock()
	// Debug: log all stored states
	storedStates := make([]string, 0, len(stateStore))
	for k := range stateStore {
		storedStates = append(storedStates, k)
	}
	log.Info().Str("receivedState", state).Strs("storedStates", storedStates).Msg("Checking OAuth state")
	expiry, exists := stateStore[state]
	if exists {
		delete(stateStore, state) // Use state only once
	}
	stateStoreMu.Unlock()

	if !exists || time.Now().After(expiry) {
		log.Warn().Str("state", state).Bool("exists", exists).Msg("Invalid OAuth state")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state parameter",
		})
		return
	}

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No code provided",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := h.spotifyClient.Exchange(ctx, code)
	if err != nil {
		log.Error().Err(err).Msg("Failed to exchange code for token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to authenticate with Spotify",
		})
		return
	}

	// Get user profile
	client := h.spotifyClient.NewSpotifyClient(ctx, token)
	user, err := h.spotifyClient.GetCurrentUser(ctx, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user profile")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile",
		})
		return
	}

	// Create session
	sessionID := uuid.New().String()
	_ = h.sessionStore.Create(sessionID, user.ID, token)

	// Create a temporary token for the frontend to complete login
	tempToken := uuid.New().String()
	tempTokenStoreMu.Lock()
	tempTokenStore[tempToken] = sessionID
	tempTokenStoreMu.Unlock()

	log.Info().Str("userID", user.ID).Msg("User logged in successfully")

	// Redirect to frontend with temp token (frontend will call /api/auth/complete to set cookie)
	c.Redirect(http.StatusFound, h.config.Server.FrontendURL+"/callback?token="+tempToken)
}

// GetMe returns the current user's profile
func (h *AuthHandler) GetMe(c *gin.Context) {
	_, exists := middleware.GetUserID(c)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	// Update session with refreshed token
	if token.AccessToken != sess.Token.AccessToken {
		sessionID, _ := c.Cookie(middleware.SessionCookieName)
		h.sessionStore.Update(sessionID, token)
	}

	client := h.spotifyClient.NewSpotifyClient(ctx, token)
	user, err := h.spotifyClient.GetCurrentUser(ctx, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user profile")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile",
		})
		return
	}

	imageURL := ""
	if len(user.Images) > 0 {
		imageURL = user.Images[0].URL
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"displayName": user.DisplayName,
		"email":       user.Email,
		"imageUrl":    imageURL,
		"product":     user.Product,
	})
}

// CompleteLogin exchanges a temp token for a session cookie (called by frontend on localhost)
func (h *AuthHandler) CompleteLogin(c *gin.Context) {
	tempToken := c.Query("token")
	if tempToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No token provided",
		})
		return
	}

	// Get and remove temp token
	tempTokenStoreMu.Lock()
	sessionID, exists := tempTokenStore[tempToken]
	if exists {
		delete(tempTokenStore, tempToken)
	}
	tempTokenStoreMu.Unlock()

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	// Set session cookie on localhost
	c.SetCookie(
		middleware.SessionCookieName,
		sessionID,
		86400, // 24 hours
		"/",
		"",
		false,
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// Logout logs out the current user
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie(middleware.SessionCookieName)
	if err == nil {
		h.sessionStore.Delete(sessionID)
	}

	c.SetCookie(middleware.SessionCookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
