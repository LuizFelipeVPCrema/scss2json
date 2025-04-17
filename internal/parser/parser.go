package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var reCommentCategory = regexp.MustCompile(`^\\s*//\\s*(.+)`)

func ParseScssFile(path string) (*ScssJsonExport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &ScssJsonExport{}
	scanner := bufio.NewScanner(file)

	var currentCategory string
	var currentRaw []string
	var currentBody []string
	var captureMode string
	var captureName string
	var captureParams []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if matches := reCommentCategory.FindStringSubmatch(line); len(matches) > 0 {
			currentCategory = matches[1]
			continue
		}

		if captureMode != "" {
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				currentBody = append(currentBody, strings.TrimSpace(line))
				currentRaw = append(currentRaw, line)
				continue
			} else {
				switch captureMode {
				case "mixin":
					result.Mixins = append(result.Mixins, ScssMixin{
						Type:   "mixin",
						Name:   captureName,
						Params: captureParams,
						Body:   currentBody,
						Raw:    strings.Join(currentRaw, "\n"),
					})
				case "function":
					result.Functions = append(result.Functions, ScssFunction{
						Type:   "function",
						Name:   captureName,
						Params: captureParams,
						Body:   currentBody,
						Raw:    strings.Join(currentRaw, "\n"),
					})
				case "placeholder":
					result.Placeholders = append(result.Placeholders, ScssPlaceholder{
						Type: "placeholder",
						Name: captureName,
						Body: currentBody,
						Raw:  strings.Join(currentRaw, "\n"),
					})
				}
				captureMode = ""
				captureName = ""
				captureParams = nil
				currentRaw = nil
				currentBody = nil
			}
		}

		// Variables
		if ok, name, value, modifier := isVariableDeclaration(line); ok {
			result.Variables = append(result.Variables, ScssVariable{
				Type:      "variable",
				Name:      name,
				Value:     value,
				Unit:      extractUnit(value),
				Raw:       line,
				Modifiers: optionalModifier(modifier),
				Category:  currentCategory,
			})
			continue
		}

		// Mixin
		if ok, name, params := isMixinDeclaration(line); ok {
			captureMode = "mixin"
			captureName = name
			captureParams = params
			currentRaw = []string{line}
			currentBody = nil
			continue
		}

		// Function
		if ok, name, params := isFunctionDeclaration(line); ok {
			captureMode = "function"
			captureName = name
			captureParams = params
			currentRaw = []string{line}
			currentBody = nil
			continue
		}

		// Placeholder
		if ok, name := isPlaceholderDeclaration(line); ok {
			captureMode = "placeholder"
			captureName = name
			currentRaw = []string{line}
			currentBody = nil
			continue
		}
	}

	return result, scanner.Err()
}

func optionalModifier(input string) []string {
	if input == "" {
		return nil
	}
	return []string{strings.TrimSpace(input)}
}
