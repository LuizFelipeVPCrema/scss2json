package parser

import "regexp"

var reVariable = regexp.MustCompile(`^\s*\$(\w[\w-]*)\s*:\s*([^;]+)(!\w+)?\s*;`)

var rePlaceholder = regexp.MustCompile(`^\s*%([a-zA-Z0-9_-]+)\s*\{?`)

var reMixin = regexp.MustCompile(`^\s*@mixin\s+([a-zA-Z0-9_-]+)\((.*?)\)\s*\{?`)

var reFunctionInline = regexp.MustCompile(`^\s*@function\s+([a-zA-Z0-9_-]+)\((.*?)\)\s*\{\s*(.*?)\s*\}`)

var reFunction = regexp.MustCompile(`^\s*@function\s+([a-zA-Z0-9_-]+)\((.*?)\)\s*\{?`)

var reLoop = regexp.MustCompile(`^@for\s+(\$\w+)\s+from\s+([0-9]+)\s+through\s+([0-9]+)\s*\{`)
