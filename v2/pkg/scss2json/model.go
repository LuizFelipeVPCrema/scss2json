package scss2json

// InputSource representa a origem dos dados SCSS
type InputSource struct {
	FilePath string
	Content  string
}

// ParseOptions define as opções de entrada para análise SCSS
type ParseOptions struct {
	Input InputSource
}
