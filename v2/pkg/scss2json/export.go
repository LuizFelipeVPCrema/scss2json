package scss2json

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/LuizFelipeVPCrema/scss2json/v2/internal/parser"
)

// ExportToJson salva o resultado do parser em um arquivo JSON no caminho especificado
func ExportToJson(ast *parser.AST, path string) error {
	if ast == nil || len(ast.Nodes) == 0 {
		return errors.New("no data to export")
	}

	// Converte AST para forma exportável (ScssNode)
	exportable := parser.ToScssNode(ast)

	// Garante que o diretório de destino exista
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// Criar o arquivo de destino
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Codifica e salva o JSON com indentação
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportable); err != nil {
		return err
	}

	return nil
}
