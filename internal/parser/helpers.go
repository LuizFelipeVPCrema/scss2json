package parser

import (
	"regexp"
	"strings"
)

func splitParams(raw string) []string {
	params := strings.Split(raw, ",")
	for i := range params {
		params[i] = strings.TrimSpace(params[i])
	}
	return params
}

func extractUnit(value string) string {
	re := regexp.MustCompile(`(?i)[0-9\.]+([a-z%]+)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
