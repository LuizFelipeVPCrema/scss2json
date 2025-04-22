package parser

import (
	"strings"
)

// MultilineCommentTracker encapsula o estado de captura de comentários multilinha.
type MultilineCommentTracker struct {
	inComment bool
	startLine int
	content   []string
	comments  []ScssComment
}

// NewMultilineCommentTracker cria um novo tracker.
func NewMultilineCommentTracker() *MultilineCommentTracker {
	return &MultilineCommentTracker{}
}

// Comments retorna todos os comentários coletados.
func (m *MultilineCommentTracker) Comments() []ScssComment {
	return m.comments
}

// ProcessLine analisa uma linha e atualiza o estado da captura.
func (m *MultilineCommentTracker) ProcessLine(line string, lineNumber int) {
	if m.inComment {
		m.content = append(m.content, line)
		if strings.Contains(line, "*/") {
			m.inComment = false
			m.comments = append(m.comments, ScssComment{
				Type:    "comment",
				Content: strings.TrimSpace(strings.Join(m.content, "\n")),
				Line:    m.startLine,
			})
			m.content = nil
		}
		return
	}

	if strings.Contains(line, "/*") {
		if strings.Contains(line, "*/") {
			// Comentário em uma linha
			start := strings.Index(line, "/*")
			end := strings.Index(line, "*/")
			if start >= 0 && end > start {
				comment := strings.TrimSpace(line[start : end+2])
				m.comments = append(m.comments, ScssComment{
					Type:    "comment",
					Content: comment,
					Line:    lineNumber,
				})
			}
		} else {
			// Início de um comentário multilinha
			m.inComment = true
			m.startLine = lineNumber
			m.content = []string{line}
		}
	}
}
