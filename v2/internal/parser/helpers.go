package parser

import (
	"bytes"
	"encoding/json"
	"strings"
)

func splitParams(raw string) []string {
	params := strings.Split(raw, ",")
	for i := range params {
		params[i] = strings.TrimSpace(params[i])
	}
	return params
}

func extractProperties(lines []string) []string {
	var props []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, ":") && strings.HasSuffix(trimmed, ";") {
			prop := strings.TrimSuffix(trimmed, ";")
			props = append(props, prop)
		}
	}
	return props
}

func ToUnescapedJSON(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

func extractBody(line string) string {
	open := strings.Index(line, "{")
	close := strings.LastIndex(line, "}")
	if open == -1 || close == -1 || open >= close {
		return ""
	}
	return strings.TrimSpace(line[open+1 : close])
}

func extractMediaCondition(line string) string {
	line = strings.TrimSpace(strings.TrimPrefix(line, "@media"))
	if idx := strings.Index(line, "{"); idx != -1 {
		line = line[:idx]
	}
	return strings.TrimSpace(line)
}

func parseBlockToNodes(lines []string) []ASTNode {
	var result []ASTNode
	var stack []*ASTRule
	var current *ASTRule
	var currentRaw []string
	var inComment bool
	var commentLines []string
	var commentStartLine int

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if inComment {
			commentLines = append(commentLines, line)
			if strings.Contains(line, "*/") {
				inComment = false
				comment := &ASTComment{
					Content: strings.Join(commentLines, "\n"),
					raw:     strings.Join(commentLines, "\n"),
					Line:    commentStartLine,
				}
				if current != nil {
					current.Children = append(current.Children, comment)
				} else {
					result = append(result, comment)
				}
				commentLines = nil
			}
			continue
		}

		if strings.Contains(trimmed, "/*") && !strings.Contains(trimmed, "*/") {
			inComment = true
			commentStartLine = i + 1
			commentLines = []string{line}
			continue
		}

		if strings.HasPrefix(trimmed, "/*") && strings.Contains(trimmed, "*/") {
			comment := &ASTComment{
				Content: trimmed,
				raw:     trimmed,
				Line:    i + 1,
			}
			if current != nil {
				current.Children = append(current.Children, comment)
			} else {
				result = append(result, comment)
			}
			continue
		}

		currentRaw = append(currentRaw, line)

		if strings.HasSuffix(trimmed, "{") {
			selector := strings.TrimSpace(strings.TrimSuffix(trimmed, "{"))
			if current != nil && strings.Contains(selector, "&") {
				selector = strings.ReplaceAll(selector, "&", current.Selector)
			} else if current != nil {
				selector = current.Selector + " " + selector
			}

			newRule := &ASTRule{
				Selector:   selector,
				Properties: []string{},
				Children:   []ASTNode{},
				Line:       i + 1,
			}

			if current != nil {
				current.Children = append(current.Children, newRule)
				stack = append(stack, current)
			} else {
				result = append(result, newRule)
			}
			current = newRule
			currentRaw = []string{line}
			continue
		}

		if trimmed == "}" {
			if current != nil {
				current.raw = strings.Join(currentRaw, "\n")
			}
			if len(stack) > 0 {
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else {
				current = nil
			}
			currentRaw = nil
			continue
		}

		if current != nil && strings.Contains(trimmed, ":") && strings.HasSuffix(trimmed, ";") {
			current.Properties = append(current.Properties, strings.TrimSuffix(trimmed, ";"))
		}
	}

	return result
}
