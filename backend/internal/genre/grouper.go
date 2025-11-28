package genre

import (
	"sort"
	"strings"
)

// GenreGroup represents a parent genre and its sub-genres
type GenreGroup struct {
	Parent   string   `json:"parent"`
	Children []string `json:"children"`
	Count    int      `json:"count"`
}

// GroupSuggestion represents a suggestion to merge genres
type GroupSuggestion struct {
	ParentGenre    string   `json:"parentGenre"`
	ChildGenres    []string `json:"childGenres"`
	TotalTracks    int      `json:"totalTracks"`
	PlaylistsToMerge int    `json:"playlistsToMerge"`
}

// Predefined genre families - maps sub-genres to parent genres
var genreFamilies = map[string]string{
	// Rock family
	"indie rock":        "Rock",
	"alternative rock":  "Rock",
	"classic rock":      "Rock",
	"hard rock":         "Rock",
	"soft rock":         "Rock",
	"progressive rock":  "Rock",
	"psychedelic rock":  "Rock",
	"garage rock":       "Rock",
	"punk rock":         "Rock",
	"post-punk":         "Rock",
	"art rock":          "Rock",
	"folk rock":         "Rock",
	"blues rock":        "Rock",
	"southern rock":     "Rock",
	"glam rock":         "Rock",
	"stoner rock":       "Rock",
	"grunge":            "Rock",
	"britpop":           "Rock",
	"rock":              "Rock",

	// Pop family
	"indie pop":         "Pop",
	"synth-pop":         "Pop",
	"synthpop":          "Pop",
	"electropop":        "Pop",
	"dream pop":         "Pop",
	"chamber pop":       "Pop",
	"art pop":           "Pop",
	"dance pop":         "Pop",
	"power pop":         "Pop",
	"baroque pop":       "Pop",
	"k-pop":             "Pop",
	"j-pop":             "Pop",
	"c-pop":             "Pop",
	"pop":               "Pop",

	// Electronic family
	"house":             "Electronic",
	"deep house":        "Electronic",
	"tech house":        "Electronic",
	"progressive house": "Electronic",
	"techno":            "Electronic",
	"trance":            "Electronic",
	"drum and bass":     "Electronic",
	"dubstep":           "Electronic",
	"edm":               "Electronic",
	"ambient":           "Electronic",
	"idm":               "Electronic",
	"downtempo":         "Electronic",
	"chillwave":         "Electronic",
	"electronica":       "Electronic",
	"electronic":        "Electronic",
	"synthwave":         "Electronic",
	"retrowave":         "Electronic",
	"vaporwave":         "Electronic",

	// Hip-Hop family
	"hip hop":           "Hip-Hop",
	"hip-hop":           "Hip-Hop",
	"rap":               "Hip-Hop",
	"trap":              "Hip-Hop",
	"southern hip hop":  "Hip-Hop",
	"east coast hip hop": "Hip-Hop",
	"west coast hip hop": "Hip-Hop",
	"underground hip hop": "Hip-Hop",
	"conscious hip hop": "Hip-Hop",
	"boom bap":          "Hip-Hop",
	"gangsta rap":       "Hip-Hop",
	"drill":             "Hip-Hop",

	// R&B family
	"r&b":               "R&B/Soul",
	"rnb":               "R&B/Soul",
	"soul":              "R&B/Soul",
	"neo soul":          "R&B/Soul",
	"contemporary r&b":  "R&B/Soul",
	"funk":              "R&B/Soul",
	"motown":            "R&B/Soul",

	// Metal family
	"heavy metal":       "Metal",
	"death metal":       "Metal",
	"black metal":       "Metal",
	"thrash metal":      "Metal",
	"progressive metal": "Metal",
	"doom metal":        "Metal",
	"power metal":       "Metal",
	"metalcore":         "Metal",
	"nu metal":          "Metal",
	"symphonic metal":   "Metal",
	"metal":             "Metal",

	// Jazz family
	"jazz":              "Jazz",
	"smooth jazz":       "Jazz",
	"acid jazz":         "Jazz",
	"jazz fusion":       "Jazz",
	"bebop":             "Jazz",
	"swing":             "Jazz",
	"big band":          "Jazz",
	"free jazz":         "Jazz",

	// Classical family
	"classical":         "Classical",
	"baroque":           "Classical",
	"romantic":          "Classical",
	"opera":             "Classical",
	"orchestral":        "Classical",
	"chamber music":     "Classical",
	"contemporary classical": "Classical",

	// Country family
	"country":           "Country",
	"country rock":      "Country",
	"alt-country":       "Country",
	"americana":         "Country",
	"bluegrass":         "Country",
	"country pop":       "Country",
	"outlaw country":    "Country",

	// Folk family
	"folk":              "Folk",
	"indie folk":        "Folk",
	"contemporary folk": "Folk",
	"acoustic":          "Folk",
	"singer-songwriter": "Folk",

	// Reggae family
	"reggae":            "Reggae",
	"dancehall":         "Reggae",
	"dub":               "Reggae",
	"ska":               "Reggae",

	// Latin family
	"latin":             "Latin",
	"latin pop":         "Latin",
	"salsa":             "Latin",
	"bachata":           "Latin",
	"cumbia":            "Latin",
	"bossa nova":        "Latin",
	"samba":             "Latin",
	"reggaeton":         "Latin",

	// Blues family
	"blues":             "Blues",
	"delta blues":       "Blues",
	"chicago blues":     "Blues",
	"electric blues":    "Blues",
}

// parentKeywords maps keywords to their parent genre for smart detection
var parentKeywords = map[string]string{
	"rock":       "Rock",
	"pop":        "Pop",
	"electronic": "Electronic",
	"electro":    "Electronic",
	"house":      "Electronic",
	"techno":     "Electronic",
	"trance":     "Electronic",
	"hip hop":    "Hip-Hop",
	"hip-hop":    "Hip-Hop",
	"rap":        "Hip-Hop",
	"r&b":        "R&B/Soul",
	"soul":       "R&B/Soul",
	"metal":      "Metal",
	"jazz":       "Jazz",
	"classical":  "Classical",
	"country":    "Country",
	"folk":       "Folk",
	"reggae":     "Reggae",
	"latin":      "Latin",
	"blues":      "Blues",
	"punk":       "Rock",
	"indie":      "Rock", // Default indie to rock, but indie pop/electronic will match first
}

// GetParentGenre returns the parent genre for a given genre, or the genre itself if no parent
func GetParentGenre(genre string) string {
	normalized := strings.ToLower(strings.TrimSpace(genre))

	// First check explicit mappings
	if parent, ok := genreFamilies[normalized]; ok {
		return parent
	}

	// Then try smart detection by checking if genre contains a keyword
	// Check longer keywords first to avoid "pop" matching before "indie pop"
	keywordsByLength := []string{
		"hip hop", "hip-hop", "electronic", "classical", "country", "reggae",
		"electro", "techno", "trance", "house", "metal", "blues", "latin",
		"rock", "jazz", "folk", "soul", "r&b", "punk", "rap", "pop", "indie",
	}

	for _, keyword := range keywordsByLength {
		if strings.Contains(normalized, keyword) {
			return parentKeywords[keyword]
		}
	}

	return genre // Return original if no mapping found
}

// GroupGenres groups a list of genres by their parent categories
func GroupGenres(genreDistribution map[string]int) map[string]*GenreGroup {
	groups := make(map[string]*GenreGroup)

	for genre, count := range genreDistribution {
		parent := GetParentGenre(genre)

		if group, exists := groups[parent]; exists {
			group.Children = append(group.Children, genre)
			group.Count += count
		} else {
			groups[parent] = &GenreGroup{
				Parent:   parent,
				Children: []string{genre},
				Count:    count,
			}
		}
	}

	// Sort children by count (most tracks first)
	for _, group := range groups {
		sort.Slice(group.Children, func(i, j int) bool {
			return genreDistribution[group.Children[i]] > genreDistribution[group.Children[j]]
		})
	}

	return groups
}

// SuggestGroupings analyzes genre distribution and suggests which genres to group
func SuggestGroupings(genreDistribution map[string]int, minTracksThreshold int) []GroupSuggestion {
	groups := GroupGenres(genreDistribution)
	var suggestions []GroupSuggestion

	for parent, group := range groups {
		// Only suggest grouping if there are multiple child genres
		if len(group.Children) > 1 {
			// Check if any individual child genre has fewer tracks than threshold
			smallGenres := []string{}
			for _, child := range group.Children {
				if genreDistribution[child] < minTracksThreshold {
					smallGenres = append(smallGenres, child)
				}
			}

			// Suggest grouping if we have small genres that could be merged
			if len(smallGenres) > 0 || len(group.Children) > 3 {
				suggestions = append(suggestions, GroupSuggestion{
					ParentGenre:      parent,
					ChildGenres:      group.Children,
					TotalTracks:      group.Count,
					PlaylistsToMerge: len(group.Children),
				})
			}
		}
	}

	// Sort by number of playlists to merge (most impact first)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].PlaylistsToMerge > suggestions[j].PlaylistsToMerge
	})

	return suggestions
}

// ApplyGrouping maps a genre to its parent genre if grouping is enabled for that family
func ApplyGrouping(genre string, enabledGroups map[string]bool) string {
	parent := GetParentGenre(genre)
	if enabledGroups[parent] {
		return parent
	}
	return genre
}

// GetAllParentGenres returns all available parent genre categories
func GetAllParentGenres() []string {
	parentSet := make(map[string]bool)
	for _, parent := range genreFamilies {
		parentSet[parent] = true
	}

	parents := make([]string, 0, len(parentSet))
	for parent := range parentSet {
		parents = append(parents, parent)
	}
	sort.Strings(parents)
	return parents
}
