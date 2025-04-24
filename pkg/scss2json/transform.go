package scss2json

import "github.com/LuizFelipeVPCrema/scss2json/internal/parser"

// ToScssNode tranforma uma AST em nós exportáveis
func ToScssNode(ast *parser.AST) []*parser.ScssNode {
	return parser.ToScssNode(ast)
}
