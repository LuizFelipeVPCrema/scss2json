package parser

func ToScssNode(ast *AST) []*ScssNode {
	var result []*ScssNode

	for _, node := range ast.Nodes {
		switch n := node.(type) {

		case *ASTVariable:
			result = append(result, &ScssNode{
				Type:  "variable",
				Name:  n.Name,
				Value: n.Value,
				Raw:   n.raw,
			})

		case *ASTMixin:
			result = append(result, &ScssNode{
				Type:   "mixin",
				Name:   n.Name,
				Params: n.Params,
				Body:   n.Body,
				Raw:    n.raw,
			})

		case *ASTFunction:
			result = append(result, &ScssNode{
				Type:   "function",
				Name:   n.Name,
				Params: n.Params,
				Body:   n.Body,
				Raw:    n.raw,
			})

		case *ASTPlaceholder:
			result = append(result, &ScssNode{
				Type: "placeholder",
				Name: n.Name,
				Body: n.Body,
				Raw:  n.raw,
			})

		case *ASTRule:
			result = append(result, &ScssNode{
				Type:     "rule",
				Selector: n.Selector,
				Props:    n.Properties,
				Nested:   n.Nested,
				Children: ToScssNode(&AST{Nodes: n.Children}),
				Raw:      n.raw,
				Line:     n.Line,
			})

		case *ASTMediaBlock:
			result = append(result, &ScssNode{
				Type:      "media",
				Condition: n.Condition,
				Children:  ToScssNode(&AST{Nodes: n.Body}),
				Raw:       n.raw,
			})

		case *ASTLoop:
			result = append(result, &ScssNode{
				Type:      "loop",
				Condition: n.Expression,
				Children:  ToScssNode(&AST{Nodes: n.Body}),
				Raw:       n.raw,
			})

		case *ASTComment:
			result = append(result, &ScssNode{
				Type:  "comment",
				Value: n.Content,
				Line:  n.Line,
				Raw:   n.raw,
			})
		default:
			// future-proofing: silently ignore unknown node types
			// optionally: log.Printf("Unhandled ASTNode type: %T", n)
		}
	}
	return result
}
