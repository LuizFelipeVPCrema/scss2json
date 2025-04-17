package scss2json

import "github.com/LuizFelipeVPCrema/scss2json/internal/parser"

func ParseFile(path string) (*parser.ScssJsonExport, error) {
	return parser.ParseScssFile(path)
}
