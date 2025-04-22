package parser

import (
	"strings"
)

func parseRules(lines []string) []ScssRule {
	var rules []ScssRule
	var stack []*RuleContext
	var braceDepth int
	seen := make(map[string]bool)

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "@media") || strings.HasPrefix(line, "@supports") || strings.HasPrefix(line, "@layer") || strings.HasPrefix(line, "@container") {
			continue
		}

		// Início de um seletor (bloco)
		if strings.HasSuffix(line, "{") {
			braceDepth++
			selector := strings.TrimSpace(strings.TrimSuffix(line, "{"))

			ctx := &RuleContext{
				selector:   selector,
				properties: []string{},
				nested:     make(map[string][]string),
			}

			// Concatenação de seletores aninhados
			if len(stack) > 0 {
				ctx.parent = stack[len(stack)-1]

				if strings.HasPrefix(selector, "&") {
					ctx.selector = selector
				} else if ctx.parent.selector != "" {
					ctx.selector = ctx.parent.selector + " " + selector
				}
			}

			stack = append(stack, ctx)
			continue
		}

		// Propriedades
		if len(stack) > 0 && strings.Contains(line, ":") && !strings.HasPrefix(line, "@") {
			line = strings.TrimSuffix(line, ";")
			stack[len(stack)-1].properties = append(stack[len(stack)-1].properties, line)
			continue
		}

		// Fim de bloco
		if strings.Contains(line, "}") {
			braceDepth--
			if braceDepth >= 0 && len(stack) > 0 {
				ctx := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if ctx.parent != nil && strings.HasPrefix(ctx.selector, "&") {
					nestedSelector := ctx.selector
					ctx.parent.nested[nestedSelector] = ctx.properties
					continue
				}

				key := strings.TrimSpace(ctx.selector)
				if key != "" && !seen[key] {
					rules = append(rules, ScssRule{
						Selector:   key,
						Properties: ctx.properties,
						Nested:     ctx.nested,
					})
					seen[key] = true
				}
			}
		}
	}

	return rules
}
