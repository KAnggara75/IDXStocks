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

	// 3. Handle DD Mon YYYY (Indonesian/English e.g. "17 Des 2009" or "17 Dec 2009")
	months := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"Mei": "05", "May": "05", "Jun": "06", "Jul": "07",
		"Agu": "08", "Agt": "08", "Aug": "08", "Sep": "09",
		"Okt": "10", "Oct": "10", "Nov": "11", "Des": "12", "Dec": "12",
	}

	parts := strings.Split(dateStr, " ")
	if len(parts) == 3 {
		day := parts[0]
		if len(day) == 1 {
			day = "0" + day
		}
		month, ok := months[parts[1]]
		if !ok {
			return dateStr
		}
		year := parts[2]
		return fmt.Sprintf("%s-%s-%s", year, month, day)
	}

	return dateStr
}
