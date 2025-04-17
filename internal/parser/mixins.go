package parser

import "regexp"

var reMixin = regexp.MustCompile(`^\s*@mixin\s+([a-zA-Z0-9_-]+)\((.*?)\)\s*\{?`)

func isMixinDeclaration(line string) (bool, string, []string) {
	if matches := reMixin.FindStringSubmatch(line); len(matches) > 0 {
		return true, matches[1], splitParams(matches[2])
	}
	return false, "", nil
}
