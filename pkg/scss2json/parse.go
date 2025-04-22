package scss2json

import (
	"errors"

	"github.com/LuizFelipeVPCrema/scss2json/internal/parser"
)

func ParseFile(path string) (*parser.ScssJsonExport, error) {
	return parser.ParseScssFile(path)
}

func parseContent(content string) (*parser.ScssJsonExport, error) {
	return parser.ParseScssContent(content)
}

type InputSource struct {
	FilePath string
	Content  string
}

type ParseOptions struct {
	Input InputSource
}

func ParseScssToJson(opts ParseOptions) (*parser.ScssJsonExport, error) {
	if opts.Input.FilePath != "" {
		return parser.ParseScssFile(opts.Input.FilePath)
	} else if opts.Input.Content != "" {
		return parser.ParseScssContent(opts.Input.Content)
	}
	return nil, errors.New("no input provided: set either FilePath or Content")
}
