package scss2json

import (
	"errors"

	"github.com/LuizFelipeVPCrema/scss2json/internal/parser"
)

// ParseFile carrega e analisa um arquivo SCSS para AST
func ParseFile(path string) (*parser.AST, error) {
	return parser.ParseScssFile(path)
}

// ParseContent analisa conteúdo SCSS diretamente (string) para AST
func ParseContent(content string) (*parser.AST, error) {
	return parser.ParseScssContent(content)
}

// ParseScss realiza o parse do SCSS (de arquivo ou string) e retorna a AST
func ParseScss(opts ParseOptions) (*parser.AST, error) {
	switch {
	case opts.Input.FilePath != "":
		return ParseFile(opts.Input.FilePath)
	case opts.Input.Content != "":
		return ParseContent(opts.Input.Content)
	default:
		return nil, errors.New("no input provided: set either FilePath or Content")
	}
}
