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

	// 3. Robust Month Mapping (Abbreviation & Full Name)
	months := map[string]string{
		"Jan": "01", "Januari": "01", "January": "01",
		"Feb": "02", "Februari": "02", "February": "02",
		"Mar": "03", "Maret": "03", "March": "03",
		"Apr": "04", "April": "04",
		"Mei": "05", "May": "05",
		"Jun": "06", "Juni": "06", "June": "06",
		"Jul": "07", "Juli": "07", "July": "07",
		"Agu": "08", "Agt": "08", "Agustus": "08", "Aug": "08", "August": "08",
		"Sep": "09", "September": "09",
		"Okt": "10", "Oktober": "10", "Oct": "10", "October": "10",
		"Nov": "11", "November": "11",
		"Des": "12", "Desember": "12", "Dec": "12", "December": "12",
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
			if val, ok := months[p]; ok {
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
