package domain

type Track struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Artists      []Artist `json:"artists"`
	AlbumName    string   `json:"albumName"`
	AlbumImage   string   `json:"albumImage"`
	Duration     int      `json:"duration"` // milliseconds
	PrimaryGenre string   `json:"primaryGenre"`
	InPlaylists  []string `json:"inPlaylists"` // Playlist IDs
}

type Artist struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}
