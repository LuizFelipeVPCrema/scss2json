<p align="right">
  🇧🇷 Português | <a href="README.en.md">🇺🇸 English</a>
</p>

<p align="center">
  <img src="https://github.com/user-attachments/assets/620d68df-62d8-4e3a-8421-88a771e6d50b" alt="SCSS2JSON Logo" width="400">
</p>


# SCSS2JSON

**SCSS2JSON** é uma biblioteca e ferramenta de linha de comando escrita em Go que realiza o _parse_ de arquivos `.scss` e transforma seu conteúdo em uma estrutura JSON baseada em **AST (Abstract Syntax Tree)**. Suporta SCSS moderno com variáveis, mixins, funções, placeholders, loops, regras aninhadas e comentários multilinha.

---

## Funcionalidades

- Extração estruturada de nós SCSS como:
  - Variáveis (`$cor: red`)
  - Mixins (`@mixin`)
  - Funções (`@function`)
  - Placeholders (`%nome`)
  - Loops `@for`
  - Regras aninhadas (`nav ul`, `a:hover`)
  - Blocos `@media`
  - Comentários (inclusive multilinha)

- Suporte a:
  - Entrada via arquivo `.scss`
  - Entrada via string SCSS (ideal para APIs)
  - Exportação JSON formatada
  - CLI simples para conversão direta

---

## Instalação

```bash
go install github.com/LuizFelipeVPCrema/scss2json@latest
```

Ou clone e compile localmente:

```bash
git clone https://github.com/LuizFelipeVPCrema/scss2json.git
cd scss2json
go build -o scss2json ./cmd/scss2json
```

---

## Como Usar (CLI)

```bash
scss2json -i input.scss -o output.json
```

- `-i`: Caminho do arquivo SCSS de entrada (**obrigatório**)
- `-o`: Caminho do arquivo JSON de saída (opcional, padrão: `output.json`)

---

## Uso como Biblioteca Go

### ➔ Parseando arquivo `.scss`

```go
result, err := scss2json.ParseFile("styles.scss")
if err != nil {
    panic(err)
}
fmt.Println(result.Nodes) // AST nodes
```

### ➔ Parseando conteúdo como string

```go
content := `$color: red;\n@mixin test { color: $color; }`
result, err := scss2json.ParseContent(content)
if err != nil {
    panic(err)
}
fmt.Println(result.Nodes)
```

---

## Estrutura do JSON (exemplo)

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
    "name": "button",
    "params": ["$radius"],
    "body": ["border-radius: $radius;"],
    "raw": "@mixin button($radius) { ... }"
  },
  ...
]
```

---

## Exemplos Avançados

### Usando em uma API HTTP com Go

```go
func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    result, _ := scss2json.ParseContent(string(body))
    json.NewEncoder(w).Encode(result)
}
```

### Exportando para disco

```go
scss2json.ExportToJson(result, "saida.json")
```

---

## API Pública

| Função                         | Descrição                               |
|-------------------------------|------------------------------------------|
| `ParseFile(path)`             | Faz o parse de um arquivo SCSS           |
| `ParseContent(content)`       | Faz o parse de conteúdo como string      |
| `ExportToJson(ast, outPath)`  | Exporta o resultado da AST em JSON       |

---

## Tecnologias

- Go 1.21+
- Expressões regulares (`regexp`)
- `encoding/json` para serialização
- `flag` para argumentos CLI

---

## Contribuições

Contribuições são bem-vindas!

1. Forke o projeto
2. Crie uma branch: `git checkout -b minha-feature`
3. Commit: `git commit -m 'feat: nova feature'`
4. Push: `git push origin minha-feature`
5. Crie um Pull Request

---

## Licença

Este projeto está licenciado sob a Licença MIT.

MIT © [LuizFelipeVPCrema](https://github.com/LuizFelipeVPCrema)


