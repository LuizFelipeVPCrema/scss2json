package parser

import "regexp"

var reFunction = regexp.MustCompile(`^\s*@function\s+([a-zA-Z0-9_-]+)\((.*?)\)`)

func isFunctionDeclaration(line string) (bool, string, []string) {
	if matches := reFunction.FindStringSubmatch(line); len(matches) > 0 {
		return true, matches[1], splitParams(matches[2])
	}
	return false, "", nil
}
