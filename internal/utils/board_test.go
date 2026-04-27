package utils

import "testing"

func TestMapBoard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Utama", "Main"},
		{"Pengembangan", "Development"},
		{"Akselerasi", "Acceleration"},
		{"Pemantauan Khusus", "Watchlist"},
		{"Ekonomi Baru", "Ekonomi Baru"},
		{"utama", "Main"},
		{"  Utama  ", "Main"},
		{"Unknown", "Main"},
	}

	for _, tt := range tests {
		got := MapBoard(tt.input)
		if got != tt.expected {
			t.Errorf("MapBoard(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
