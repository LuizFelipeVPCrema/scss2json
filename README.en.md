<p align="right">
  <a href="README.md">🇧🇷 Português</a> | 🇺🇸 English
</p>


# SCSS2JSON

**A smart and robust Go parser for SCSS → JSON.**

This library converts SCSS files or strings into structured JSON objects, supporting variables, mixins, functions, placeholders, nested rules, and loops.

---

## Features

- ✅ SCSS Variables with units and optional modifiers
- ✅ Mixins and Functions (including inline)
- ✅ Placeholders (`%`)
- ✅ CSS Rules and nested selectors
- ✅ `@for` Loops with scoped rules
- ✅ CLI and programmatic use
- ✅ High-quality JSON output with nesting support

---

## Installation

```bash
go get github.com/LuizFelipeVPCrema/scss2json
```

---

## CLI Usage

```bash
go run main.go -i path/to/input.scss -o output.json
```

- `-i`: Input SCSS file path (**required**)
- `-o`: Output JSON file path (default: `output.json`)

---

## Usage as a Library

### Parse from a File

```go
import "github.com/LuizFelipeVPCrema/scss2json"

result, err := scss2json.ParseFile("styles.scss")
if err != nil {
    log.Fatal(err)
}
```

### Parse from String

```go
import "github.com/LuizFelipeVPCrema/scss2json"

scss := `
$color: red;
@mixin rounded { border-radius: 10px; }
`

result, err := scss2json.ParseString(scss)
if err != nil {
    log.Fatal(err)
}
```

### Export to JSON

```go
err := scss2json.ExportToJson(result, "output.json")
```

---

## Run Tests

```bash
go test ./...
```

---

## Output Sample

```json
{
  "variables": [
    {
      "type": "variable",
      "name": "color",
      "value": "red",
      "unit": "",
      "raw": "$color: red;"
    }
  ],
  "mixins": [
    {
      "type": "mixin",
      "name": "rounded",
      "params": [],
      "body": ["border-radius: 10px;"]
    }
  ]
}
```

---

## License

This project is licensed under the MIT License.

Copyright (c) 2025 LuizFelipeVPCrema



## Advanced Examples

- Nested CSS selectors (`a:hover`, `.card .title`)
- Rules inside loops (`.column-#{$i}`)
- Visual debugging via exported JSON
- Use it in servers or CLI tools for SCSS inspection


## Author

**Luiz Felipe Crema**  
[GitHub Profile](https://github.com/LuizFelipeVPCrema)


## Useful Links

- [MIT License](LICENSE)


---

