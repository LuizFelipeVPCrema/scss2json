package parser

type ScssVariable struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Value     string   `json:"value"`
	Unit      string   `json:"unit"`
	Raw       string   `json:"raw"`
	Modifiers []string `json:"modifiers"`
	Category  string   `json:"category,omitempty"`
}

type ScssMixin struct {
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Params []string `json:"params"`
	Body   []string `json:"body"`
	Raw    string   `json:"raw"`
}

type ScssFunction struct {
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Params []string `json:"params"`
	Body   []string `json:"body"`
	Raw    string   `json:"raw"`
}

type ScssPlaceholder struct {
	Type string   `json:"type"`
	Name string   `json:"name"`
	Body []string `json:"body"`
	Raw  string   `json:"raw"`
}

type ScssRule struct {
	Selector   string              `json:"selector"`
	Properties []string            `json:"properties"`
	Nested     map[string][]string `json:"nested,omitempty"`
}

type RuleContext struct {
	selector   string
	properties []string
	parent     *RuleContext
	nested     map[string][]string
}

type ScssLoop struct {
	Type       string     `json:"type"`
	Expression string     `json:"expression"`
	Body       []ScssRule `json:"body"`
	Raw        string     `json:"raw"`
}

type ScssJsonExport struct {
	Variables    []ScssVariable    `json:"variables"`
	Mixins       []ScssMixin       `json:"mixins"`
	Functions    []ScssFunction    `json:"functions"`
	Placeholders []ScssPlaceholder `json:"placeholders"`
	Rules        []ScssRule        `json:"rules"`
	Loops        []ScssLoop        `json:"loops"`
	Comments     []ScssComment     `json:"comments"`
}

type ScssComment struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Line    int    `json:"line"`
}
