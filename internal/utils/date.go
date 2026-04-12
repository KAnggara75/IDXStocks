package utils

import (
	"fmt"
	"strings"
)

// NormalizeDate converts various date formats (DD-MM-YYYY, DD Mon YYYY) to YYYY-MM-DD
func NormalizeDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}

	// 1. If it's already YYYY-MM-DD, return it
	if len(dateStr) == 10 && dateStr[4] == '-' && dateStr[7] == '-' {
		return dateStr
	}

	// 2. Robust Month Mapping (Abbreviation & Full Name) - All Lowercase
	months := map[string]string{
		"jan": "01", "januari": "01", "january": "01",
		"feb": "02", "februari": "02", "february": "02",
		"mar": "03", "maret": "03", "march": "03",
		"apr": "04", "april": "04",
		"mei": "05", "may": "05",
		"jun": "06", "juni": "06", "june": "06",
		"jul": "07", "juli": "07", "july": "07",
		"agu": "08", "agt": "08", "agustus": "08", "aug": "08", "august": "08",
		"sep": "09", "september": "09",
		"okt": "10", "oktober": "10", "oct": "10", "october": "10",
		"nov": "11", "november": "11",
		"des": "12", "desember": "12", "dec": "12", "december": "12",
	}

	// 3. Handle various separators
	sep := "-"
	if strings.Contains(dateStr, " ") {
		sep = " "
	} else if strings.Contains(dateStr, "/") {
		sep = "/"
	}

	parts := strings.Split(dateStr, sep)
	if len(parts) == 3 {
		var d, m, y string
		mIdx := -1

		// 3a. Identify Month (either by name or by position in common formats)
		for i, p := range parts {
			pLow := strings.ToLower(p)
			if val, ok := months[pLow]; ok {
				m = val
				mIdx = i
				break
			}
		}

		// 3b. If no month name found, assume numeric month and identify parts by common patterns
		if mIdx == -1 {
			// Case: YYYY-MM-DD (already handled at top, but just in case)
			if len(parts[0]) == 4 {
				y = parts[0]
				m = parts[1]
				d = parts[2]
			} else if len(parts[2]) == 4 {
				y = parts[2]
				m = parts[1]
				d = parts[0]
			}
		} else {
			// Month name was found, identify Year and Day from others
			var others []string
			for i, p := range parts {
				if i != mIdx {
					others = append(others, p)
				}
			}
			if len(others) == 2 {
				if len(others[0]) == 4 {
					y = others[0]
					d = others[1]
				} else if len(others[1]) == 4 {
					y = others[1]
					d = others[0]
				}
			}
		}

		// Final assembly if all parts identified
		if y != "" && m != "" && d != "" {
			if len(m) == 1 {
				m = "0" + m
			}
			if len(d) == 1 {
				d = "0" + d
			}
			return fmt.Sprintf("%s-%s-%s", y, m, d)
		}
	}

	return dateStr
}
