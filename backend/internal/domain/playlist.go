package domain

const ManagedTag = "[Managed by SpotifyPlaylistSorter]"

type Playlist struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	OwnerID       string   `json:"ownerId"`
	TrackCount    int      `json:"trackCount"`
	ImageURL      string   `json:"imageUrl"`
	ManagedByApp  bool     `json:"managedByApp"`  // Has our tag in description
	AssignedGenre string   `json:"assignedGenre"` // Genre this playlist represents
	TrackIDs      []string `json:"trackIds"`
}

func (p *Playlist) IsManagedByApp() bool {
	return p.ManagedByApp
}
