package parser

import "regexp"

var rePlaceholder = regexp.MustCompile(`^\s*%([a-zA-Z0-9_-]+)\s*$`)

func isPlaceholderDeclaration(line string) (bool, string) {
	if matches := rePlaceholder.FindStringSubmatch(line); len(matches) > 0 {
		return true, matches[1]
	}
	return false, ""
}
