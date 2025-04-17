package parser

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

func splitParams(raw string) []string {
	params := strings.Split(raw, ",")
	for i := range params {
		params[i] = strings.TrimSpace(params[i])
	}
	return params
}

func extractUnit(value string) string {
	re := regexp.MustCompile(`(?i)[0-9\.]+([a-z%]+)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func ToUnescapedJSON(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}
