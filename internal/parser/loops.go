package parser

import (
	"regexp"
	"strings"
)

var reLoop = regexp.MustCompile(`^@for\s+(\$\w+)\s+from\s+([0-9]+)\s+through\s+([0-9]+)\s*\{`)

func isLoopDeclaration(line string) (bool, string) {
	if matches := reLoop.FindStringSubmatch(line); len(matches) > 0 {
		expr := matches[1] + " from " + matches[2] + " through " + matches[3]
		return true, expr
	}
	return false, ""
}

func parseLoops(lines []string) []ScssLoop {
	var loops []ScssLoop
	var loop ScssLoop
	var raw []string
	collecting := false
	braceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !collecting {
			if ok, expr := isLoopDeclaration(trimmed); ok {
				collecting = true
				loop = ScssLoop{
					Type:       "for",
					Expression: expr,
					Body:       []ScssRule{},
				}
			}
		}

		if collecting {
			raw = append(raw, line)
			braceDepth += strings.Count(line, "{")
			braceDepth -= strings.Count(line, "}")

			if braceDepth == 0 {
				loop.Raw = strings.Join(raw, "\n")
				loop.Body = parseRules(raw[1:])
				loops = append(loops, loop)
				collecting = false
				raw = nil
			}
		}
	}

	return loops
}
