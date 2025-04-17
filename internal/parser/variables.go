package parser

import (
	"regexp"
	"strings"
)

var reVariable = regexp.MustCompile(`^\s*\$(\w[\w-]*)\s*:\s*([^;]+)(!\w+)?\s*;`)

func isVariableDeclaration(line string) (bool, string, string, string) {
	if matches := reVariable.FindStringSubmatch(line); len(matches) > 0 {
		name := matches[1]
		value := strings.TrimSpace(matches[2])
		modifier := strings.TrimSpace(matches[3])
		return true, name, value, modifier
	}
	return false, "", "", ""
}
