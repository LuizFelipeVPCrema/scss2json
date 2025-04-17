package parser_test

import (
	"bytes"
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
		if !containsLine(mixin.Body, "border-radius: $radius") {
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
		if !containsLine(fn.Body, "@return $number * 2") {
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
		if !containsLine(ph.Body, "@include border-radius(5px)") {
			t.Errorf("corpo do placeholder incorreto")
		}
	}

	// Verifica rules
	if len(result.Rules) == 0 {
		t.Errorf("nenhuma regra CSS encontrada")
	} else {
		// Encontra a regra "nav ul"
		navUlFound := false
		for _, r := range result.Rules {
			if r.Selector == "nav ul" {
				navUlFound = true
				if !containsLine(r.Properties, "margin: 0") {
					t.Errorf("propriedade 'margin: 0' não encontrada na regra 'nav ul'")
				}
				break
			}
		}
		if !navUlFound {
			t.Errorf("regra esperada 'nav ul' não encontrada")
		}

		// Encontra a regra "nav a" com nested rules
		navAFound := false
		for _, r := range result.Rules {
			if r.Selector == "nav a" {
				navAFound = true
				if len(r.Nested) == 0 {
					t.Errorf("regra aninhada não encontrada em 'nav a'")
				} else if _, ok := r.Nested["&:hover"]; !ok {
					t.Errorf("regra aninhada '&:hover' não encontrada em 'nav a'")
				}
				break
			}
		}
		if !navAFound {
			t.Errorf("regra esperada 'nav a' não encontrada")
		}
	}

	// Verifica loops
	if len(result.Loops) == 0 {
		t.Errorf("nenhum loop for encontrado")
	} else {
		loop := result.Loops[0]
		if loop.Expression != "$i from 1 through 3" {
			t.Errorf("expressão do loop incorreta: %s", loop.Expression)
		}
		if len(loop.Body) == 0 || !strings.Contains(loop.Body[0].Selector, ".column-") {
			t.Errorf("corpo do loop não contém a regra esperada")
		}
	}

	// Exporta resultado para inspeção manual (debug)
	rawJSON, err := parser.ToUnescapedJSON(result)
	if err != nil {
		t.Fatalf("Erro ao serializar resultado para JSON: %v", err)
	}

	// Formatando com indentação
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, rawJSON, "", "  ")
	if err != nil {
		t.Fatalf("Erro ao formatar JSON: %v", err)
	}

	err = os.WriteFile("testdata/saida.json", prettyJSON.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Erro ao escrever arquivo de saída: %v", err)
	}
}

func containsLine(lines []string, expected string) bool {
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimSuffix(line, ";"))
		expected = strings.TrimSpace(strings.TrimSuffix(expected, ";"))
		if line == expected {
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

func TestParseScssContent(t *testing.T) {
	scss := `
$primary-color: #3498db;
$padding: 10px;

@mixin border-radius($radius) {
  border-radius: $radius;
}

@function double($n) {
  @return $n * 2;
}

%box-style {
  padding: $padding;
}

nav {
  ul {
    margin: 0;
  }
}
`

	result, err := parser.ParseScssContent(scss)
	if err != nil {
		t.Fatalf("Erro ao parsear SCSS de string: %v", err)
	}

	// Verifica variáveis
	if len(result.Variables) != 2 {
		t.Errorf("esperado 2 variáveis, obtido %d", len(result.Variables))
	}

	// Verifica mixins
	if len(result.Mixins) != 1 || result.Mixins[0].Name != "border-radius" {
		t.Errorf("esperado 1 mixin 'border-radius', obtido %+v", result.Mixins)
	}

	// Verifica função
	if len(result.Functions) != 1 || result.Functions[0].Name != "double" {
		t.Errorf("esperado 1 função 'double', obtido %+v", result.Functions)
	}

	// Verifica placeholder
	if len(result.Placeholders) != 1 || result.Placeholders[0].Name != "box-style" {
		t.Errorf("esperado 1 placeholder 'box-style', obtido %+v", result.Placeholders)
	}

	// Verifica regra CSS nav ul
	found := false
	for _, rule := range result.Rules {
		if strings.TrimSpace(rule.Selector) == "nav ul" {
			found = true
			if len(rule.Properties) == 0 || rule.Properties[0] != "margin: 0" {
				t.Errorf("nav ul encontrado, mas propriedades incorretas: %+v", rule.Properties)
			}
			break
		}
	}
	if !found {
		t.Error("regra 'nav ul' não encontrada")
	}
}
