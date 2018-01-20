package maya

import "strings"

func sanitizeLineFeedMultiLine(lines []string) []string {
	for i, line := range lines {
		lines[i] = sanitizeLineFeedSingleLine(line)
	}
	return lines
}

func sanitizeLineFeedSingleLine(line string) string {
	return strings.Replace(line, "\r", "", -1)
}
