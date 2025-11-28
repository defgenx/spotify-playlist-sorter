package sse

import (
	"encoding/json"
	"sync"
)

// EventType represents different types of progress events
type EventType string

const (
	EventTypeProgress EventType = "progress"
	EventTypeError    EventType = "error"
	EventTypeComplete EventType = "complete"
	EventTypeInfo     EventType = "info"
)

// ProgressPhase represents different phases of the operation
type ProgressPhase string

const (
	PhaseFetchingLikedSongs  ProgressPhase = "fetching_liked_songs"
	PhaseFetchingPlaylists   ProgressPhase = "fetching_playlists"
	PhaseFetchingArtists     ProgressPhase = "fetching_artists"
	PhaseAnalyzing           ProgressPhase = "analyzing"
	PhaseGeneratingPlan      ProgressPhase = "generating_plan"
	PhaseCreatingPlaylists   ProgressPhase = "creating_playlists"
	PhaseAddingTracks        ProgressPhase = "adding_tracks"
	PhaseRemovingTracks      ProgressPhase = "removing_tracks"
	PhaseComplete            ProgressPhase = "complete"
)

// ProgressEvent represents a progress update event
type ProgressEvent struct {
	Type    EventType     `json:"type"`
	Phase   ProgressPhase `json:"phase,omitempty"`
	Current int           `json:"current,omitempty"`
	Total   int           `json:"total,omitempty"`
	Message string        `json:"message"`
}

// Client represents a connected SSE client
type Client struct {
	UserID  string
	Channel chan *ProgressEvent
}

// Broadcaster manages SSE connections and broadcasts events
type Broadcaster struct {
	clients map[string][]*Client // userID -> clients
	mu      sync.RWMutex
}

// NewBroadcaster creates a new SSE broadcaster
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[string][]*Client),
	}
}

// Subscribe adds a new client for a user
func (b *Broadcaster) Subscribe(userID string) *Client {
	b.mu.Lock()
	defer b.mu.Unlock()

	client := &Client{
		UserID:  userID,
		Channel: make(chan *ProgressEvent, 100),
	}

	b.clients[userID] = append(b.clients[userID], client)
	return client
}

// Unsubscribe removes a client
func (b *Broadcaster) Unsubscribe(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()

	clients := b.clients[client.UserID]
	for i, c := range clients {
		if c == client {
			// Remove from slice
			b.clients[client.UserID] = append(clients[:i], clients[i+1:]...)
			close(c.Channel)
			break
		}
	}

	// Clean up empty user entries
	if len(b.clients[client.UserID]) == 0 {
		delete(b.clients, client.UserID)
	}
}

// Broadcast sends an event to all clients for a specific user
func (b *Broadcaster) Broadcast(userID string, event *ProgressEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	clients := b.clients[userID]
	for _, client := range clients {
		select {
		case client.Channel <- event:
			// Successfully sent
		default:
			// Channel full, skip
		}
	}
}

// SendProgress sends a progress update
func (b *Broadcaster) SendProgress(userID string, phase ProgressPhase, current, total int, message string) {
	b.Broadcast(userID, &ProgressEvent{
		Type:    EventTypeProgress,
		Phase:   phase,
		Current: current,
		Total:   total,
		Message: message,
	})
}

// SendInfo sends an info message
func (b *Broadcaster) SendInfo(userID string, message string) {
	b.Broadcast(userID, &ProgressEvent{
		Type:    EventTypeInfo,
		Message: message,
	})
}

// SendError sends an error message
func (b *Broadcaster) SendError(userID string, message string) {
	b.Broadcast(userID, &ProgressEvent{
		Type:    EventTypeError,
		Message: message,
	})
}

// SendComplete sends a completion message
func (b *Broadcaster) SendComplete(userID string, message string) {
	b.Broadcast(userID, &ProgressEvent{
		Type:    EventTypeComplete,
		Phase:   PhaseComplete,
		Message: message,
	})
}

// FormatSSE formats an event as SSE message
func FormatSSE(event *ProgressEvent) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	// SSE format: data: {json}\n\n
	return []byte("data: " + string(data) + "\n\n"), nil
}

// GetClientCount returns the number of connected clients for a user
func (b *Broadcaster) GetClientCount(userID string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.clients[userID])
}
