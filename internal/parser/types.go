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

// Interface de nó base
type ASTNode interface {
	NodeType() string
	Raw() string
}

// Estrutura de árvore principal
type AST struct {
	Nodes []ASTNode
}

// =========================
// ===== NÓS INDIVIDUAIS ====
// =========================
type ASTVariable struct {
	Name  string
	Value string
	raw   string
}

func (v *ASTVariable) NodeType() string { return "variable" }
func (v *ASTVariable) Raw() string      { return v.raw }

type ASTMixin struct {
	Name   string
	Params []string
	Body   []string
	raw    string
}

func (m *ASTMixin) NodeType() string { return "mixin" }
func (m *ASTMixin) Raw() string      { return m.raw }

type ASTFunction struct {
	Name   string
	Params []string
	Body   []string
	raw    string
}

func (f *ASTFunction) NodeType() string { return "function" }
func (f *ASTFunction) Raw() string      { return f.raw }

type ASTPlaceholder struct {
	Name string
	Body []string
	raw  string
}

func (p *ASTPlaceholder) NodeType() string { return "placeholder" }
func (p *ASTPlaceholder) Raw() string      { return p.raw }

type ASTLoop struct {
	Expression string
	Body       []ASTNode
	raw        string
}

func (l *ASTLoop) NodeType() string { return "loop" }
func (l *ASTLoop) Raw() string      { return l.raw }

type ASTMediaBlock struct {
	Condition string
	Body      []ASTNode
	raw       string
}

func (m *ASTMediaBlock) NodeType() string { return "media" }
func (m *ASTMediaBlock) Raw() string      { return m.raw }

type ASTRule struct {
	Selector   string
	Properties []string
	Nested     map[string][]string
	Children   []ASTNode
	raw        string
	Line       int
}

func (r *ASTRule) NodeType() string { return "rule" }
func (r *ASTRule) Raw() string      { return r.raw }

type ASTComment struct {
	Content string
	Line    int
	raw     string
}

func (c *ASTComment) NodeType() string { return "comment" }
func (c *ASTComment) Raw() string      { return c.raw }
