package parser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// var reCommentCategory = regexp.MustCompile(`^\\s*//\\s*(.+)`)
var reLoop = regexp.MustCompile(`^@for\s+(\$\w+)\s+from\s+([0-9]+)\s+through\s+([0-9]+)\s*\{`)

func ParseScssFile(path string) (*AST, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseAST(strings.NewReader(string(contentBytes))), nil
}

func ParseScssContent(content string) (*AST, error) {
	return ParseAST(strings.NewReader(content)), nil
}

func ParseAST(reader io.Reader) *AST {
	scanner := bufio.NewScanner(reader)
	ast := &AST{}
	commentTracker := NewMultilineCommentTracker()

	var (
		inBlock        bool
		blockType      string
		blockName      string
		blockParams    []string
		blockCondition string
		blockLines     []string
		blockRaw       []string
		lineNumber     int
		braceDepth     int
	)

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		trimmed := strings.TrimSpace(line)

		commentTracker.ProcessLine(line, lineNumber)

		// Ignora linhas vazias e comentários inline
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		// Dentro de um bloco (mixin, function, rule, etc)
		if inBlock {
			blockRaw = append(blockRaw, line)
			blockLines = append(blockLines, trimmed)

			braceDepth += strings.Count(trimmed, "{")
			braceDepth -= strings.Count(trimmed, "}")

			if braceDepth == 0 {
				inBlock = false
				raw := strings.Join(blockRaw, "\n")
				bodyLines := blockLines[1 : len(blockLines)-1]

				switch blockType {
				case "mixin":
					ast.Nodes = append(ast.Nodes, &ASTMixin{
						Name:   blockName,
						Params: blockParams,
						Body:   bodyLines,
						raw:    raw,
					})
				case "function":
					ast.Nodes = append(ast.Nodes, &ASTFunction{
						Name:   blockName,
						Params: blockParams,
						Body:   bodyLines,
						raw:    raw,
					})
				case "placeholder":
					ast.Nodes = append(ast.Nodes, &ASTPlaceholder{
						Name: blockName,
						Body: bodyLines,
						raw:  raw,
					})
				case "media":
					ast.Nodes = append(ast.Nodes, &ASTMediaBlock{
						Condition: blockCondition,
						Body:      parseBlockToNodes(bodyLines),
						raw:       raw,
					})
				case "loop":
					ast.Nodes = append(ast.Nodes, &ASTLoop{
						Expression: blockCondition,
						Body:       parseBlockToNodes(bodyLines),
						raw:        raw,
					})
				case "rule":
					parsed := parseBlockToNodes(bodyLines)
					rule := &ASTRule{
						Selector: blockName,
						raw:      raw,
						Line:     lineNumber - len(blockLines) + 1,
					}
					for _, child := range parsed {
						switch n := child.(type) {
						case *ASTRule:
							rule.Children = append(rule.Children, n)
						case *ASTComment:
							rule.Children = append(rule.Children, n)
						}
					}
					rule.Properties = extractProperties(bodyLines)
					ast.Nodes = append(ast.Nodes, rule)
				}

				// Reset bloco
				inBlock = false
				blockType, blockName, blockCondition = "", "", ""
				blockParams, blockLines, blockRaw = nil, nil, nil
			}
			continue
		}

		// Declaração de variável
		if ok, name, value, _ := isVariableDeclaration(line); ok {
			ast.Nodes = append(ast.Nodes, &ASTVariable{
				Name:  "$" + name,
				Value: value,
				raw:   line,
			})
			continue
		}

		// Mixin
		if ok, name, params := isMixinDeclaration(line); ok {
			if strings.Contains(line, "{") && strings.Contains(line, "}") {
				ast.Nodes = append(ast.Nodes, &ASTMixin{
					Name:   name,
					Params: params,
					Body:   []string{extractBody(line)},
					raw:    line,
				})
			} else {
				inBlock = true
				blockType = "mixin"
				blockName = name
				blockParams = params
				blockLines = []string{trimmed}
				blockRaw = []string{line}
				braceDepth = 1
			}
			continue
		}

		// Função
		if ok, name, params, inlineBody, raw := isFunctionDeclaration(line); ok {
			if inlineBody != nil {
				ast.Nodes = append(ast.Nodes, &ASTFunction{
					Name:   name,
					Params: params,
					Body:   inlineBody,
					raw:    raw,
				})
			} else {
				inBlock = true
				blockType = "function"
				blockName = name
				blockParams = params
				blockLines = []string{trimmed}
				blockRaw = []string{line}
				braceDepth = 1
			}
			continue
		}

		// Placeholder
		if ok, name := isPlaceholderDeclaration(line); ok {
			inBlock = true
			blockType = "placeholder"
			blockName = name
			blockLines = []string{trimmed}
			blockRaw = []string{line}
			braceDepth = 1
			continue
		}

		// Media query
		if strings.HasPrefix(trimmed, "@media") {
			condition := extractMediaCondition(trimmed)
			inBlock = true
			blockType = "media"
			blockCondition = condition
			blockLines = []string{trimmed}
			blockRaw = []string{line}
			braceDepth = 1
			continue
		}

		// Loop
		if matches := reLoop.FindStringSubmatch(trimmed); len(matches) > 0 {
			expr := matches[1] + " from " + matches[2] + " through " + matches[3]
			inBlock = true
			blockType = "loop"
			blockCondition = expr
			blockLines = []string{trimmed}
			blockRaw = []string{line}
			braceDepth = 1
			continue
		}

		// Regra CSS normal
		if strings.HasSuffix(trimmed, "{") {
			selector := strings.TrimSuffix(trimmed, "{")
			inBlock = true
			blockType = "rule"
			blockName = strings.TrimSpace(selector)
			blockLines = []string{trimmed}
			blockRaw = []string{line}
			braceDepth = 1
			continue
		}
	}

	// Adiciona comentários finais
	ast.Nodes = append(ast.Nodes, commentTracker.Comments()...)
	return ast
}

// func ParseAST(reader io.Reader) *AST {
// 	scanner := bufio.NewScanner(reader)
// 	ast := &AST{}
// 	commentTracker := NewMultilineCommentTracker()

// 	var (
// 		inBlock        bool
// 		blockType      string
// 		blockName      string
// 		blockParams    []string
// 		blockCondition string
// 		blockLines     []string
// 		blockRaw       []string
// 		lineNumber     int
// 		braceDepth     int
// 	)

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		lineNumber++
// 		trimmed := strings.TrimSpace(line)

// 		commentTracker.ProcessLine(line, lineNumber)

// 		// Ignora linhas vazias e comentários por enquanto
// 		// TODO
// 		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
// 			continue
// 		}

// 		// Se estiver dentro de um bloco (mixin, function, media, etc)
// 		if inBlock {
// 			blockRaw = append(blockRaw, line)
// 			blockLines = append(blockLines, trimmed)

// 			braceDepth += strings.Count(trimmed, "{")
// 			braceDepth -= strings.Count(trimmed, "}")

// 			if braceDepth == 0 {
// 				inBlock = false

// 				switch blockType {
// 				case "mixin":
// 					ast.Nodes = append(ast.Nodes, &ASTMixin{
// 						Name:   blockName,
// 						Params: blockParams,
// 						Body:   blockLines[1 : len(blockLines)-1],
// 						raw:    strings.Join(blockRaw, "\n"),
// 					})
// 				case "function":
// 					ast.Nodes = append(ast.Nodes, &ASTFunction{
// 						Name:   blockName,
// 						Params: blockParams,
// 						Body:   blockLines[1 : len(blockLines)-1],
// 						raw:    strings.Join(blockRaw, "\n"),
// 					})
// 				case "media":
// 					ast.Nodes = append(ast.Nodes, &ASTMediaBlock{
// 						Condition: blockCondition,
// 						Body:      parseBlockToNodes(blockLines[1 : len(blockLines)-1]),
// 						raw:       strings.Join(blockRaw, "\n"),
// 					})
// 				case "loop":
// 					ast.Nodes = append(ast.Nodes, &ASTLoop{
// 						Expression: blockCondition,
// 						Body:       parseBlockToNodes(blockLines[1 : len(blockLines)-1]),
// 						raw:        strings.Join(blockRaw, "\n"),
// 					})
// 				case "placeholder":
// 					ast.Nodes = append(ast.Nodes, &ASTPlaceholder{
// 						Name: blockName,
// 						Body: blockLines[1 : len(blockLines)-1],
// 						raw:  strings.Join(blockRaw, "\n"),
// 					})
// 				case "rule":
// 					ast.Nodes = append(ast.Nodes, &ASTRule{
// 						Selector:   blockName,
// 						Properties: extractProperties(blockLines[1 : len(blockLines)-1]),
// 						raw:        strings.Join(blockRaw, "\n"),
// 						Line:       lineNumber - len(blockLines) + 1,
// 					})
// 				}

// 				// Reset
// 				inBlock = false
// 				blockType, blockName, blockCondition = "", "", ""
// 				blockParams, blockLines, blockRaw = nil, nil, nil
// 			}
// 			continue
// 		}

// 		// === Blocos possíveis ===

// 		// Variáveis
// 		if ok, name, value, _ := isVariableDeclaration(line); ok {
// 			ast.Nodes = append(ast.Nodes, &ASTVariable{
// 				Name:  "$" + name,
// 				Value: value,
// 				raw:   line,
// 			})
// 			continue
// 		}

// 		// Mixins (inline por enquanto)
// 		if ok, name, params := isMixinDeclaration(line); ok {
// 			if strings.Contains(line, "{") && strings.Contains(line, "}") {
// 				body := extractBody(line)
// 				ast.Nodes = append(ast.Nodes, &ASTMixin{
// 					Name:   name,
// 					Params: params,
// 					Body:   []string{body},
// 					raw:    line,
// 				})
// 			} else {
// 				// início de bloco
// 				inBlock = true
// 				blockType = "mixin"
// 				blockName = name
// 				blockParams = params
// 				blockLines = []string{trimmed}
// 				blockRaw = []string{line}
// 				braceDepth = 1
// 			}
// 			continue
// 		}

// 		// Funções (mesma ideia)
// 		if ok, name, params, inlineBody, raw := isFunctionDeclaration(line); ok {
// 			if inlineBody != nil {
// 				ast.Nodes = append(ast.Nodes, &ASTFunction{
// 					Name:   name,
// 					Params: params,
// 					Body:   inlineBody,
// 					raw:    raw,
// 				})
// 			} else {
// 				// início de bloco
// 				inBlock = true
// 				blockType = "function"
// 				blockName = name
// 				blockParams = params
// 				blockLines = []string{trimmed}
// 				blockRaw = []string{line}
// 				braceDepth = 1
// 			}
// 			continue
// 		}

// 		// Placeholder simples
// 		if ok, name := isPlaceholderDeclaration(line); ok {
// 			inBlock = true
// 			blockType = "placeholder"
// 			blockName = name
// 			blockLines = []string{trimmed}
// 			blockRaw = []string{line}
// 			braceDepth = 1
// 			continue
// 		}

// 		// Media
// 		if strings.HasPrefix(trimmed, "@media") {
// 			condition := extractMediaCondition(trimmed)
// 			inBlock = true
// 			blockType = "media"
// 			blockCondition = condition
// 			blockLines = []string{trimmed}
// 			blockRaw = []string{line}
// 			braceDepth = 1
// 			continue
// 		}

// 		// Loop
// 		if matches := reLoop.FindStringSubmatch(trimmed); len(matches) > 0 {
// 			expr := matches[1] + " from " + matches[2] + " through " + matches[3]
// 			inBlock = true
// 			blockType = "loop"
// 			blockCondition = expr
// 			blockLines = []string{trimmed}
// 			blockRaw = []string{line}
// 			braceDepth = 1
// 			continue
// 		}

// 		// Regras CSS (inline)
// 		if strings.HasSuffix(trimmed, "{") {
// 			selector := strings.TrimSuffix(trimmed, "{")
// 			inBlock = true
// 			blockType = "rule"
// 			blockName = strings.TrimSpace(selector)
// 			blockLines = []string{trimmed}
// 			blockRaw = []string{line}
// 			braceDepth = 1
// 			continue
// 		}
// 	}

// 	// Comentários
// 	ast.Nodes = append(ast.Nodes, commentTracker.Comments()...)
// 	return ast
// }
