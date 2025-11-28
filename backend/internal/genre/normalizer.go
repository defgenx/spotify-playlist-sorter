package genre

import (
	"regexp"
	"strings"
)

var (
	// Special characters to remove during normalization
	specialCharsRegex = regexp.MustCompile(`[^\w\s-]`)
	// Multiple spaces/dashes to single space
	multiSpaceRegex = regexp.MustCompile(`[\s-]+`)
)

// NormalizeGenre normalizes a genre string by lowercasing and removing special characters
func NormalizeGenre(genre string) string {
	if genre == "" {
		return ""
	}

	// Lowercase
	normalized := strings.ToLower(genre)

	// Remove special characters
	normalized = specialCharsRegex.ReplaceAllString(normalized, " ")

	// Replace multiple spaces/dashes with single space
	normalized = multiSpaceRegex.ReplaceAllString(normalized, " ")

	// Trim spaces
	normalized = strings.TrimSpace(normalized)

	return normalized
}

// MatchPlaylistToGenre performs fuzzy matching between playlist name and available genres
// Returns the best matching genre and a confidence score (0-1)
func MatchPlaylistToGenre(playlistName string, genres []string) (string, float64) {
	if len(genres) == 0 {
		return "", 0.0
	}

	normalizedPlaylist := NormalizeGenre(playlistName)
	if normalizedPlaylist == "" {
		return "", 0.0
	}

	bestMatch := ""
	bestScore := 0.0

	for _, genre := range genres {
		normalizedGenre := NormalizeGenre(genre)
		if normalizedGenre == "" {
			continue
		}

		score := calculateMatchScore(normalizedPlaylist, normalizedGenre)
		if score > bestScore {
			bestScore = score
			bestMatch = genre
		}
	}

	return bestMatch, bestScore
}

// calculateMatchScore calculates a similarity score between two strings
func calculateMatchScore(s1, s2 string) float64 {
	// Exact match
	if s1 == s2 {
		return 1.0
	}

	// Contains match (genre is substring of playlist name)
	if strings.Contains(s1, s2) {
		return 0.9
	}

	// Contains match (playlist name is substring of genre)
	if strings.Contains(s2, s1) {
		return 0.85
	}

	// Word-based matching
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Count common words
	commonWords := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 {
				commonWords++
				break
			}
		}
	}

	if commonWords == 0 {
		return 0.0
	}

	// Calculate Jaccard similarity
	totalUniqueWords := len(uniqueWords(append(words1, words2...)))
	score := float64(commonWords) / float64(totalUniqueWords)

	return score
}

// uniqueWords returns a slice of unique words
func uniqueWords(words []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}

	return result
}

// ExtractPrimaryGenre extracts the primary genre from a list of genres
// Uses a simple heuristic: prefer specific genres over generic ones
func ExtractPrimaryGenre(genres []string) string {
	if len(genres) == 0 {
		return ""
	}

	// Generic genres to deprioritize
	genericGenres := map[string]bool{
		"pop":           true,
		"rock":          true,
		"electronic":    true,
		"indie":         true,
		"alternative":   true,
	}

	// First pass: look for specific genres
	for _, genre := range genres {
		normalized := NormalizeGenre(genre)
		if !genericGenres[normalized] {
			return genre
		}
	}

	// If all are generic or no specific found, return first
	return genres[0]
}
