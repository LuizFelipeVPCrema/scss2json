package parser

import (
	"strings"
)

// MultilineCommentTracker captura comentários multilinha como nós ASTComment.
type MultilineCommentTracker struct {
	inComment    bool
	startLine    int
	contentLines []string
	collected    []ASTNode
}

// NewMultilineCommentTracker cria um novo tracker.
func NewMultilineCommentTracker() *MultilineCommentTracker {
	return &MultilineCommentTracker{}
}

// Comments retorna os nós ASTComment gerados.
func (m *MultilineCommentTracker) Comments() []ASTNode {
	return m.collected
}

// ProcessLine analisa a linha e atualiza o estado interno da captura de comentários.
func (m *MultilineCommentTracker) ProcessLine(line string, lineNumber int) {
	switch {
	case m.inComment:
		m.contentLines = append(m.contentLines, line)
		if strings.Contains(line, "*/") {
			m.inComment = false
			m.collected = append(m.collected, &ASTComment{
				Content: strings.TrimSpace(strings.Join(m.contentLines, "\n")),
				Line:    m.startLine,
				raw:     strings.Join(m.contentLines, "\n"),
			})
			m.contentLines = nil
		}

	case strings.Contains(line, "/*") && strings.Contains(line, "*/"):
		start := strings.Index(line, "/*")
		end := strings.Index(line, "*/")
		if start >= 0 && end > start {
			comment := strings.TrimSpace(line[start : end+2])
			m.collected = append(m.collected, &ASTComment{
				Content: comment,
				Line:    lineNumber,
				raw:     comment,
			})
		}

	case strings.Contains(line, "/*"):
		m.inComment = true
		m.startLine = lineNumber
		m.contentLines = []string{line}
	}
}
