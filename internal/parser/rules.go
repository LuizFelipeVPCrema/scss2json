package parser

import (
	"regexp"
	"strings"
)

var reRuleStart = regexp.MustCompile(`^([.#a-zA-Z0-9\s,&:\[\]\-_=><"']+?)\s*\{`)

func isRuleDeclaration(line string) (bool, string) {
	if matches := reRuleStart.FindStringSubmatch(line); len(matches) > 0 {
		return true, strings.TrimSpace(matches[1])
	}
	return false, ""
}

func parseRules(lines []string) []ScssRule {
	var rules []ScssRule
	var stack []*RuleContext
	braceDepth := 0
	seen := make(map[string]bool)

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		if ok, selector := isRuleDeclaration(line); ok {
			braceDepth++
			ctx := &RuleContext{
				selector:   selector,
				properties: []string{},
				nested:     make(map[string][]string),
			}

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

		if strings.Contains(line, "{") {
			braceDepth++
		}

		if strings.Contains(line, "}") {
			braceDepth--
			if braceDepth >= 0 && len(stack) > 0 {
				ctx := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if strings.Contains(ctx.selector, "#{") && strings.Contains(ctx.selector, "$i") {
					continue
				}

				if ctx.parent != nil && strings.HasPrefix(ctx.selector, "&") {
					nestedSelector := ctx.selector
					ctx.parent.nested[nestedSelector] = ctx.properties
					continue
				}

				key := strings.TrimSpace(ctx.selector)
				if !seen[key] && key != "" {
					rules = append(rules, ScssRule{
						Selector:   key,
						Properties: ctx.properties,
						Nested:     ctx.nested,
					})
					seen[key] = true
				}
			}
			continue
		}

		if len(stack) > 0 && !strings.Contains(line, "@") && strings.Contains(line, ":") {

			line = strings.TrimSuffix(line, ";")
			stack[len(stack)-1].properties = append(stack[len(stack)-1].properties, line)
		}
	}

	return rules
}
