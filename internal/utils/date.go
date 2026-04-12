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

	// 2. Handle DD-MM-YYYY (often from Pasardana)
	if strings.Contains(dateStr, "-") {
		parts := strings.Split(dateStr, "-")
		if len(parts) == 3 && len(parts[2]) == 4 {
			day := parts[0]
			if len(day) == 1 {
				day = "0" + day
			}
			month := parts[1]
			if len(month) == 1 {
				month = "0" + month
			}
			return fmt.Sprintf("%s-%s-%s", parts[2], month, day)
		}
	}

	// 3. Robust Month Mapping (Abbreviation & Full Name) - All Lowercase
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

	// 4. Handle "-" or " " separated dates
	sep := "-"
	if strings.Contains(dateStr, " ") {
		sep = " "
	}

	parts := strings.Split(dateStr, sep)
	if len(parts) == 3 {
		// Detect positions
		var d, m, y string

		// Try to identify month
		mIdx := -1
		for i, p := range parts {
			if val, ok := months[strings.ToLower(p)]; ok {
				m = val
				mIdx = i
				break
			}
		}

		if mIdx != -1 {
			// If month is found, identify day and year from remaining parts
			var others []string
			for i, p := range parts {
				if i != mIdx {
					others = append(others, p)
				}
			}

			if len(others) == 2 {
				// Usually year is 4 digits
				if len(others[1]) == 4 {
					y = others[1]
					d = others[0]
				} else if len(others[0]) == 4 {
					y = others[0]
					d = others[1]
				}

				if y != "" && d != "" {
					if len(d) == 1 {
						d = "0" + d
					}
					return fmt.Sprintf("%s-%s-%s", y, m, d)
				}
			}
		}
	}

	return dateStr
}
