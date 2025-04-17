package scss2json

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/LuizFelipeVPCrema/scss2json/internal/parser"
)

func ExportToJson(result *parser.ScssJsonExport, path string) error {
	if result == nil {
		return errors.New("no data to export")
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
