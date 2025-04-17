// rules.go com suporte a nested e exclusão de regras duplicadas vindas de loops
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

type ruleContext struct {
	selector   string
	properties []string
	parent     *ruleContext
	nested     map[string][]string
}

func parseRules(lines []string) []ScssRule {
	var rules []ScssRule
	var stack []*ruleContext
	braceDepth := 0
	seen := make(map[string]bool)

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		if ok, selector := isRuleDeclaration(line); ok {
			braceDepth++
			ctx := &ruleContext{
				selector:   selector,
				properties: []string{},
				nested:     make(map[string][]string),
			}
			if len(stack) > 0 {
				ctx.parent = stack[len(stack)-1]
				if strings.HasPrefix(selector, "&") {
					ctx.selector = strings.Replace(selector, "&", ctx.parent.selector, 1)
				} else {
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

				if strings.Contains(ctx.selector, "#{") {
					continue
				}

				if ctx.parent != nil && strings.HasPrefix(ctx.selector, ctx.parent.selector+":") {
					ctx.parent.nested[ctx.selector[len(ctx.parent.selector):]] = ctx.properties
					continue
				}

				key := strings.TrimSpace(ctx.selector)
				if !seen[key] {
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

		if len(stack) > 0 {
			stack[len(stack)-1].properties = append(stack[len(stack)-1].properties, line)
		}
	}

	return rules
}
