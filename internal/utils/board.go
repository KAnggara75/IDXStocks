package utils

import "strings"

// MapBoard mapping Indonesian board name to DB Enum board
func MapBoard(board string) string {
	switch strings.ToLower(strings.TrimSpace(board)) {
	case "utama":
		return "Main"
	case "pengembangan":
		return "Development"
	case "akselerasi":
		return "Acceleration"
	case "pemantauan khusus":
		return "Watchlist"
	case "ekonomi baru":
		return "Ekonomi Baru"
	default:
		// Default to "Main" if empty or unknown to match DB default
		if strings.TrimSpace(board) == "" {
			return "Main"
		}
		return "Main"
	}
}
