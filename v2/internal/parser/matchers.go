package parser

import "strings"

func isVariableDeclaration(line string) (bool, string, string, string) {
	if matches := reVariable.FindStringSubmatch(line); len(matches) > 0 {
		name := matches[1]
		value := strings.TrimSpace(matches[2])
		modifier := strings.TrimSpace(matches[3])
		return true, name, value, modifier
	}
	return false, "", "", ""
}

func isPlaceholderDeclaration(line string) (bool, string) {
	if matches := rePlaceholder.FindStringSubmatch(line); len(matches) > 0 {
		return true, matches[1]
	}
	return false, ""
}

func isMixinDeclaration(line string) (bool, string, []string) {
	if matches := reMixin.FindStringSubmatch(line); len(matches) > 0 {
		return true, matches[1], splitParams(matches[2])
	}
	return false, "", nil
}

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
