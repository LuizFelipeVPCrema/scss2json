package parser

// ScssNode representa um nó genérica da AST(Árvore de Sintaxe Abstrata) de SCSS
type ScssNode struct {
	Type      string              `json:"type"`            // ex: "varible", "mixin", "function", "rule", etc
	Name      string              `json:"name,omitempty"`  // nome da variável, mixin, seletor etc
	Value     string              `json:"value,omitempty"` // valor de variável ou propriedade
	Unit      string              `json:"unit,omitempty"`
	Params    []string            `json:"params,omitempty"`     // parâmetros para mixins e funções
	Selector  string              `json:"selector,omitempty"`   // seletor CSS (para rules)
	Condition string              `json:"condition,omitempty"`  // condições do bloco @media, @supperts, etc
	Body      []string            `json:"body,omitempty"`       // corpo textual (usando em mixin, funções, placeholders)
	Props     []string            `json:"properties,omitempty"` // propriedades CSS (para rules)
	Children  []*ScssNode         `json:"children,omitempty"`   // filhos aninhados (ex: blocos dentro do for, @media etc)
	Nested    map[string][]string `json:"nested,omitempty"`     // regras aninhadas via "&"
	Line      int                 `json:"line,omitempty"`       // número da linha no SCSS original
	Raw       string              `json:"raw,omitempty"`        // conteúdo bruto
	Category  string              `json:"category,omitempty"`
	Modifiers []string            `json:"modifiers,omitempty"` // ex: modificadores extras (!default)
}
