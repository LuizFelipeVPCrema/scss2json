package parser

import (
	"regexp"
	"strings"
)

var reMediaStart = regexp.MustCompile(`^@(media|supports|layer|container)\s+(.+)\s*\{`)

func parseMediaBlocks(lines []string) []ScssMediaBlock {
	var blocks []ScssMediaBlock
	var inBlock bool
	var currentType, currentCondition string
	var currentRaw []string
	var currentBody []string
	var currentProperties []string
	braceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if !inBlock {
			if matches := reMediaStart.FindStringSubmatch(trimmed); len(matches) > 0 {
				inBlock = true
				currentType = matches[1]
				currentCondition = matches[2]
				currentRaw = []string{line}
				currentBody = []string{}
				currentProperties = []string{}
				braceDepth = 1
				continue
			}
		} else {
			currentRaw = append(currentRaw, line)
			currentBody = append(currentBody, line)

			braceDepth += strings.Count(line, "{")
			braceDepth -= strings.Count(line, "}")

			if strings.Contains(line, ":") && !strings.HasPrefix(trimmed, "@") {
				currentProperties = append(currentProperties, strings.TrimSpace(strings.TrimSuffix(trimmed, ";")))
			}

			if braceDepth == 0 {
				blocks = append(blocks, ScssMediaBlock{
					Type:       currentType,
					Condition:  currentCondition,
					Properties: currentProperties,
					Rules:      parseRules(currentBody),
					Raw:        strings.Join(currentRaw, "\n"),
				})
				inBlock = false
				currentRaw = nil
				currentBody = nil
				currentProperties = nil
			}
		}
	}

	return blocks
}

func FilterOutMediaBlocks(lines []string) []string {
	var filtered []string
	inMedia := false
	braceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if !inMedia && reMediaStart.MatchString(trimmed) {
			inMedia = true
			braceDepth = 1
			continue
		}

		if inMedia {
			braceDepth += strings.Count(trimmed, "{")
			braceDepth -= strings.Count(trimmed, "}")
			if braceDepth == 0 {
				inMedia = false
			}
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}
