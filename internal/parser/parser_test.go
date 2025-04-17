package parser_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LuizFelipeVPCrema/scss2json/internal/parser"
)

func TestParseSampleScss(t *testing.T) {
	samplePath := filepath.Join("testdata", "sample.scss")
	result, err := parser.ParseScssFile(samplePath)
	if err != nil {
		t.Fatalf("Erro ao parsear sample.scss: %v", err)
	}

	// Verifica variáveis
	if len(result.Variables) != 3 {
		t.Errorf("esperado 3 variáveis, obtido %d", len(result.Variables))
	}

	// Verifica mixins
	if len(result.Mixins) != 1 {
		t.Errorf("esperado 1 mixin, obtido %d", len(result.Mixins))
	} else {
		mixin := result.Mixins[0]
		if mixin.Name != "border-radius" {
			t.Errorf("nome do mixin inválido: esperado 'border-radius', obtido '%s'", mixin.Name)
		}
		if len(mixin.Params) != 1 || mixin.Params[0] != "$radius" {
			t.Errorf("parâmetros do mixin inválidos: %+v", mixin.Params)
		}
		if !containsLine(mixin.Body, "border-radius: $radius;") {
			t.Errorf("corpo do mixin incompleto: %+v", mixin.Body)
		}
	}

	// Verifica funções
	if len(result.Functions) != 1 {
		t.Errorf("esperado 1 função, obtido %d", len(result.Functions))
	} else {
		fn := result.Functions[0]
		if fn.Name != "double" {
			t.Errorf("nome da função inválido: %s", fn.Name)
		}
		if !containsLine(fn.Body, "@return $number * 2;") {
			t.Errorf("corpo da função não encontrado corretamente")
		}
	}

	// Verifica placeholders
	if len(result.Placeholders) != 1 {
		t.Errorf("esperado 1 placeholder, obtido %d", len(result.Placeholders))
	} else {
		ph := result.Placeholders[0]
		if ph.Name != "button-style" {
			t.Errorf("nome do placeholder inválido: %s", ph.Name)
		}
		if !containsLine(ph.Body, "@include border-radius(5px);") {
			t.Errorf("corpo do placeholder incorreto")
		}
	}

	// Exporta resultado para inspeção manual (debug)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	_ = os.WriteFile("testdata/saida.json", jsonData, 0644)
}

func containsLine(lines []string, expected string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) == expected {
			return true
		}
	}
	return false
}

func TestParseEmptyFile(t *testing.T) {
	f, err := os.CreateTemp("", "empty.scss")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	result, err := parser.ParseScssFile(f.Name())
	if err != nil {
		t.Errorf("Erro inesperado ao parsear arquivo vazio: %v", err)
	}
	if result == nil {
		t.Error("Resultado nulo para arquivo vazio")
	}
}

func TestInvalidSyntaxIgnored(t *testing.T) {
	path := filepath.Join("testdata", "invalid.scss")
	content := `
        $valid: 10px;
        this is not SCSS;
        @function test() { @return 1; }
    `
	_ = os.WriteFile(path, []byte(content), 0644)
	defer os.Remove(path)

	result, err := parser.ParseScssFile(path)
	if err != nil {
		t.Fatalf("Erro inesperado ao parsear SCSS inválido: %v", err)
	}
	if len(result.Variables) != 1 {
		t.Errorf("Esperado 1 variável válida, obtido %d", len(result.Variables))
	}
	if len(result.Functions) != 1 {
		t.Errorf("Esperado 1 função válida, obtido %d", len(result.Functions))
	}
}

func TestParseLoopPreservation(t *testing.T) {
	path := filepath.Join("testdata", "loop.scss")
	content := `
        @for $i from 1 through 3 {
            .column-#{$i} {
                width: 100% / $i;
            }
        }
    `
	_ = os.WriteFile(path, []byte(content), 0644)
	defer os.Remove(path)

	result, err := parser.ParseScssFile(path)
	if err != nil {
		t.Fatalf("Erro ao parsear loop: %v", err)
	}

	raw := ""
	if len(result.Functions) > 0 {
		raw = result.Functions[0].Raw
	}

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	_ = os.WriteFile("testdata/loop_output.json", jsonBytes, 0644)

	// Esse teste garante apenas que o loop é ignorado e não quebra nada.
	if result == nil || raw == "" {
		t.Log("Loop foi ignorado (como esperado por enquanto).")
	}
}
