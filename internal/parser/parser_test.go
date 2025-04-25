package parser_test

import (
	"strings"
	"testing"

	"github.com/LuizFelipeVPCrema/scss2json/internal/parser"
)

func assertNodeCount(t *testing.T, ast *parser.AST, expected int) {
	if len(ast.Nodes) != expected {
		t.Fatalf("esperado %d nós, obtido %d", expected, len(ast.Nodes))
	}
}

func TestAST_Mixin(t *testing.T) {
	scss := `@mixin rounded($radius) { border-radius: $radius; }`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	mixin, ok := ast.Nodes[0].(*parser.ASTMixin)
	if !ok {
		t.Fatalf("esperado ASTMixin, obtido %T", ast.Nodes[0])
	}
	if mixin.Name != "rounded" {
		t.Errorf("nome do mixin incorreto: %s", mixin.Name)
	}
	if len(mixin.Params) != 1 || mixin.Params[0] != "$radius" {
		t.Errorf("parametros incorretos: %+v", mixin.Params)
	}
	if len(mixin.Body) == 0 {
		t.Errorf("esperado corpo do mixin")
	}
}

func TestAST_Function(t *testing.T) {
	scss := `@function double($n) { @return $n * 2; }`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	fn, ok := ast.Nodes[0].(*parser.ASTFunction)
	if !ok {
		t.Fatalf("esperado ASTFunction, obtido %T", ast.Nodes[0])
	}
	if fn.Name != "double" {
		t.Errorf("nome da função incorreto: %s", fn.Name)
	}
	if len(fn.Params) != 1 || fn.Params[0] != "$n" {
		t.Errorf("parametros incorretos: %+v", fn.Params)
	}
	if len(fn.Body) == 0 {
		t.Errorf("esperado corpo da função")
	}
}

func TestAST_Variable(t *testing.T) {
	scss := `$color: #ff0000;`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	varDecl, ok := ast.Nodes[0].(*parser.ASTVariable)
	if !ok {
		t.Fatalf("esperado ASTVariable, obtido %T", ast.Nodes[0])
	}
	if varDecl.Name != "$color" || varDecl.Value != "#ff0000" {
		t.Errorf("variavel incorreta: %+v", varDecl)
	}
}

func TestAST_Placeholder(t *testing.T) {
	scss := `%button-style {
		display: inline-block;
		padding: 10px;
	}`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	ph, ok := ast.Nodes[0].(*parser.ASTPlaceholder)
	if !ok {
		t.Fatalf("esperado ASTPlaceholder, obtido %T", ast.Nodes[0])
	}
	if ph.Name != "button-style" {
		t.Errorf("nome do placeholder incorreto: %s", ph.Name)
	}
	if len(ph.Body) == 0 {
		t.Errorf("esperado conteúdo no corpo do placeholder")
	}
}

func TestAST_Media(t *testing.T) {
	scss := `@media screen and (min-width: 900px) {
		.container {
			width: 100%;
		}
	}`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	media, ok := ast.Nodes[0].(*parser.ASTMediaBlock)
	if !ok {
		t.Fatalf("esperado ASTMediaBlock, obtido %T", ast.Nodes[0])
	}
	if media.Condition != "screen and (min-width: 900px)" {
		t.Errorf("condição do media block incorreta: %s", media.Condition)
	}
	if len(media.Body) == 0 {
		t.Errorf("esperado regras dentro do media block")
	}
}

func TestAST_Loop(t *testing.T) {
	scss := `@for $i from 1 through 3 {
		.item-#{$i} {
			width: 100px;
		}
	}`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	loop, ok := ast.Nodes[0].(*parser.ASTLoop)
	if !ok {
		t.Fatalf("esperado ASTLoop, obtido %T", ast.Nodes[0])
	}
	if loop.Expression != "$i from 1 through 3" {
		t.Errorf("expressão incorreta: %s", loop.Expression)
	}
	if len(loop.Body) == 0 {
		t.Errorf("esperado corpo no loop")
	}
}

func TestAST_Comment(t *testing.T) {
	scss := `/* Este é um comentário */`
	ast := parser.ParseAST(strings.NewReader(scss))
	assertNodeCount(t, ast, 1)

	comment, ok := ast.Nodes[0].(*parser.ASTComment)
	if !ok {
		t.Fatalf("esperado ASTComment, obtido %T", ast.Nodes[0])
	}
	if comment.Content == "" {
		t.Errorf("conteúdo do comentário está vazio")
	}
	if comment.Line != 1 {
		t.Errorf("linha incorreta do comentário, esperada 1, obtida %d", comment.Line)
	}
}

func TestAST_Rule(t *testing.T) {
	scss := `
	nav {
		ul {
			margin: 0;
			padding: 0;
			list-style: none;
		}
		li {
			display: inline-block;
		}
		a {
			text-decoration: none;
			color: $primary-color;
			&:hover {
				color: darken($primary-color, 10%);
			}
		}
	}`

	ast := parser.ParseAST(strings.NewReader(scss))
	if len(ast.Nodes) == 0 {
		t.Fatal("esperado ao menos 1 nó na AST")
	}

	rootRule := findRuleBySelectorRecursive(ast.Nodes, "nav")
	if rootRule == nil {
		t.Fatal("esperado regra 'nav' não encontrada")
	}

	if rootRule.Selector != "nav" {
		t.Errorf("esperado seletor 'nav', obtido '%s'", rootRule.Selector)
	}

	// Verifica "nav ul"
	ul := findRuleBySelectorRecursive(ast.Nodes, "nav ul")
	if ul == nil {
		t.Error("esperado seletor 'nav ul' não encontrado")
	} else {
		expectProps := []string{"margin: 0", "padding: 0", "list-style: none"}
		assertProperties(t, ul.Properties, expectProps)
	}

	// Verifica "nav li"
	li := findRuleBySelectorRecursive(ast.Nodes, "nav li")
	if li == nil {
		t.Error("esperado seletor 'nav li' não encontrado")
	} else {
		assertProperties(t, li.Properties, []string{"display: inline-block"})
	}

	// Verifica "nav a"
	a := findRuleBySelectorRecursive(ast.Nodes, "nav a")
	if a == nil {
		t.Error("esperado seletor 'nav a' não encontrado")
	} else {
		expectProps := []string{"text-decoration: none", "color: $primary-color"}
		assertProperties(t, a.Properties, expectProps)

		// Verifica "nav a:hover"
		hover := findRuleBySelectorRecursive(ast.Nodes, "nav a a:hover")
		if hover == nil {
			t.Error("esperado seletor 'nav a:hover' não encontrado")
		} else {
			assertProperties(t, hover.Properties, []string{"color: darken($primary-color, 10%)"})
		}
	}
}

// Helpers Rules
func findRuleBySelectorRecursive(nodes []parser.ASTNode, path string) *parser.ASTRule {
	parts := strings.Split(path, " ")
	return findRecursiveHelper(nodes, parts)
}

func findRecursiveHelper(nodes []parser.ASTNode, parts []string) *parser.ASTRule {
	if len(parts) == 0 {
		return nil
	}

	for _, node := range nodes {
		rule, ok := node.(*parser.ASTRule)
		if !ok || rule.Selector != parts[0] {
			continue
		}

		if len(parts) == 1 {
			return rule
		}

		return findRecursiveHelper(rule.Children, parts[1:])
	}

	return nil
}

func assertProperties(t *testing.T, got []string, want []string) {
	if len(got) != len(want) {
		t.Errorf("esperado %d propriedades, obtido %d", len(want), len(got))
	}
	for _, w := range want {
		found := false
		for _, g := range got {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("propriedade esperada não encontrada: %s", w)
		}
	}
}
