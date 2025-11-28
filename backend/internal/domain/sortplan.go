package domain

import "time"

type SortPlan struct {
	ID                  string       `json:"id"`
	CreatedAt           time.Time    `json:"createdAt"`
	DryRun              bool         `json:"dryRun"`
	TotalLikedTracks    int          `json:"totalLikedTracks"`
	TracksToAdd         []TrackMove  `json:"tracksToAdd"`
	TracksToRemove      []TrackMove  `json:"tracksToRemove"`
	PlaylistsToCreate   []string     `json:"playlistsToCreate"` // Genre names
	UncategorizedTracks []Track      `json:"uncategorizedTracks"`
	GenreStats          []GenreStat  `json:"genreStats"`
}

type TrackMove struct {
	TrackID        string `json:"trackId"`
	TrackName      string `json:"trackName"`
	ArtistName     string `json:"artistName"`
	AlbumImage     string `json:"albumImage"`
	Genre          string `json:"genre"`
	FromPlaylist   string `json:"fromPlaylist"`   // Playlist ID (empty if from liked songs)
	FromPlaylistName string `json:"fromPlaylistName"`
	ToPlaylist     string `json:"toPlaylist"`     // Playlist ID or genre name if new
	ToPlaylistName string `json:"toPlaylistName"`
	Reason         string `json:"reason"`
}

type GenreStat struct {
	Genre      string `json:"genre"`
	TrackCount int    `json:"trackCount"`
	PlaylistID string `json:"playlistId"` // Empty if needs to be created
	IsNew      bool   `json:"isNew"`
}

type ExecutionResult struct {
	Success           bool              `json:"success"`
	PlaylistsCreated  int               `json:"playlistsCreated"`
	TracksAdded       int               `json:"tracksAdded"`
	TracksRemoved     int               `json:"tracksRemoved"`
	Errors            []ExecutionError  `json:"errors"`
}

type ExecutionError struct {
	Operation string `json:"operation"`
	TrackID   string `json:"trackId,omitempty"`
	Playlist  string `json:"playlist,omitempty"`
	Error     string `json:"error"`
}
