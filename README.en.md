<p align="right">
  <a href="README.md">🇧🇷 Português</a> | 🇺🇸 English
</p>

<p align="center">
  <img src="https://github.com/user-attachments/assets/620d68df-62d8-4e3a-8421-88a771e6d50b" alt="SCSS2JSON Logo" width="400">
</p>

# SCSS2JSON

**A smart and robust Go parser for SCSS → JSON.**

This library parses SCSS files or raw content and outputs a full **AST (Abstract Syntax Tree)** in JSON format. It supports variables, mixins, functions, placeholders, loops, media blocks, comments, and deeply nested CSS rules.

---


## Features

- ✅ Full AST-based SCSS parsing
- ✅ Supports:
  - Variables (`$color`)
  - Mixins (`@mixin`)
  - Functions (`@function`)
  - Placeholders (`%placeholder`)
  - Loops (`@for`)
  - Media blocks (`@media`)
  - Nested CSS rules (like `a:hover`, `nav ul li`)
  - Multi-line comments
- ✅ File or string input
- ✅ Pretty JSON output
- ✅ CLI and Go library integration

---

## Installation

```bash
go install github.com/LuizFelipeVPCrema/scss2json@latest
```

Or clone and build locally:

```bash
git clone https://github.com/LuizFelipeVPCrema/scss2json.git
cd scss2json
go build -o scss2json ./cmd/scss2json
```

---

## CLI Usage

```bash
scss2json -i path/to/input.scss -o output.json
```

- `-i`: Input SCSS file (**required**)
- `-o`: Output JSON path (default: `output.json`)

---

## Using as a Library

### Parse a SCSS file

```go
result, err := scss2json.ParseFile("styles.scss")
if err != nil {
    panic(err)
}
fmt.Println(result.Nodes) // AST nodes
```

### Parse a SCSS string

```go
content := `$color: red;\n@mixin test { color: $color; }`
result, err := scss2json.ParseContent(content)
if err != nil {
    panic(err)
}
fmt.Println(result.Nodes)
```

### Export AST to JSON

```go
scss2json.ExportToJson(result, "output.json")
```

---

## Output Example (AST JSON)

```json
[
  {
    "type": "variable",
    "name": "$color",
    "value": "red",
    "raw": "$color: red;"
  },
  {
    "type": "mixin",
    "name": "rounded",
    "params": ["$radius"],
    "body": ["border-radius: $radius;"],
    "raw": "@mixin rounded($radius) { ... }"
  },
  {
    "type": "rule",
    "selector": ".button",
    "properties": ["color: $color"],
    "children": []
  }
]
```

---

## Advanced Use Cases

- Extract selectors with `&` logic (`a:hover`)
- Parse SCSS for linting or visualization
- Embed into backend systems or build tools
- Generate style metadata for UI libraries

---

## Run Tests

```bash
go test ./...
```

---

## API Summary

| Function                      | Description                              |
|------------------------------|------------------------------------------|
| `ParseFile(path)`            | Parses a SCSS file and returns AST       |
| `ParseContent(content)`      | Parses raw SCSS string content           |
| `ExportToJson(ast, path)`    | Saves AST result as JSON                 |

---

## License

This project is licensed under the MIT License.  
MIT © [LuizFelipeVPCrema](https://github.com/LuizFelipeVPCrema)

---

## Author

**Luiz Felipe Crema**  
[GitHub Profile](https://github.com/LuizFelipeVPCrema)
