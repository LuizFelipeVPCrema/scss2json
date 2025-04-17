package parser

import (
	"regexp"
	"strings"
)

var reFunctionInline = regexp.MustCompile(`^\s*@function\s+([a-zA-Z0-9_-]+)\((.*?)\)\s*\{\s*(.*?)\s*\}`)

var reFunction = regexp.MustCompile(`^\s*@function\s+([a-zA-Z0-9_-]+)\((.*?)\)\s*\{?`)

func isFunctionDeclaration(line string) (bool, string, []string, []string, string) {
	// Função inline com conteúdo direto
	if matches := reFunctionInline.FindStringSubmatch(line); len(matches) > 0 {
		name := matches[1]
		params := splitParams(matches[2])
		body := []string{strings.TrimSpace(matches[3])}
		raw := line
		return true, name, params, body, raw
	}

	// Função que começa um bloco
	if matches := reFunction.FindStringSubmatch(line); len(matches) > 0 {
		return true, matches[1], splitParams(matches[2]), nil, ""
	}

	return false, "", nil, nil, ""
}
