package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/api/middleware"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/sse"
)

// EventsHandler handles SSE endpoints
type EventsHandler struct {
	broadcaster *sse.Broadcaster
}

// NewEventsHandler creates a new events handler
func NewEventsHandler(broadcaster *sse.Broadcaster) *EventsHandler {
	return &EventsHandler{
		broadcaster: broadcaster,
	}
}

// StreamEvents streams progress events via SSE
func (h *EventsHandler) StreamEvents(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	log.Info().Str("userID", userID).Msg("SSE client connected")

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	// Subscribe to events
	client := h.broadcaster.Subscribe(userID)
	defer func() {
		h.broadcaster.Unsubscribe(client)
		log.Info().Str("userID", userID).Msg("SSE client disconnected")
	}()

	// Send initial connection message
	initialEvent := &sse.ProgressEvent{
		Type:    sse.EventTypeInfo,
		Message: "Connected to event stream",
	}
	if data, err := sse.FormatSSE(initialEvent); err == nil {
		c.Writer.Write(data)
		c.Writer.Flush()
	}

	// Keep track of client connection
	clientGone := c.Writer.CloseNotify()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-clientGone:
			// Client disconnected
			return

		case event, ok := <-client.Channel:
			if !ok {
				// Channel closed
				return
			}

			// Format and send event
			data, err := sse.FormatSSE(event)
			if err != nil {
				log.Error().Err(err).Msg("Failed to format SSE event")
				continue
			}

			if _, err := c.Writer.Write(data); err != nil {
				if err == io.EOF {
					log.Debug().Str("userID", userID).Msg("Client disconnected (EOF)")
					return
				}
				log.Error().Err(err).Msg("Failed to write SSE event")
				return
			}

			c.Writer.Flush()

		case <-ticker.C:
			// Send keep-alive comment
			if _, err := c.Writer.Write([]byte(": keep-alive\n\n")); err != nil {
				log.Debug().Str("userID", userID).Msg("Failed to send keep-alive, client likely disconnected")
				return
			}
			c.Writer.Flush()
		}
	}
}

// TestEvent sends a test event (for debugging)
func (h *EventsHandler) TestEvent(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	h.broadcaster.SendInfo(userID, "Test event from server")

	c.JSON(http.StatusOK, gin.H{
		"message": "Test event sent",
	})
}
