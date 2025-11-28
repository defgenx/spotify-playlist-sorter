package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/session"
)

const (
	SessionCookieName = "spotify_session"
	UserIDContextKey  = "userID"
	SessionContextKey = "session"
)

// AuthMiddleware validates the session and adds user context
func AuthMiddleware(store *session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session cookie
		sessionID, err := c.Cookie(SessionCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No session found",
			})
			c.Abort()
			return
		}

		// Get session from store
		sess, err := store.Get(sessionID)
		if err != nil {
			if err == session.ErrSessionExpired {
				// Clear expired cookie
				c.SetCookie(SessionCookieName, "", -1, "/", "", false, true)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Session expired",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid session",
				})
			}
			c.Abort()
			return
		}

		// Refresh session expiration
		if err := store.Refresh(sessionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to refresh session",
			})
			c.Abort()
			return
		}

		// Add user ID and session to context
		c.Set(UserIDContextKey, sess.UserID)
		c.Set(SessionContextKey, sess)

		c.Next()
	}
}

// OptionalAuthMiddleware adds user context if session exists, but doesn't require it
func OptionalAuthMiddleware(store *session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(SessionCookieName)
		if err == nil {
			sess, err := store.Get(sessionID)
			if err == nil {
				c.Set(UserIDContextKey, sess.UserID)
				c.Set(SessionContextKey, sess)
			}
		}
		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetSession retrieves the session from the context
func GetSession(c *gin.Context) (*session.Session, bool) {
	sess, exists := c.Get(SessionContextKey)
	if !exists {
		return nil, false
	}
	return sess.(*session.Session), true
}
