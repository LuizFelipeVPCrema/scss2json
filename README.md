<p align="right">
  🇧🇷 Português | <a href="README.en.md">🇺🇸 English</a>
</p>

# SCSS2JSON 

**SCSS2JSON** é uma biblioteca e ferramenta de linha de comando escrita em Go que faz o _parse_ de arquivos `.scss` e transforma seus conteúdos em uma estrutura JSON bem definida. Ela é capaz de identificar variáveis, mixins, funções, placeholders, regras CSS e loops como `@for`, suportando hierarquia e aninhamentos SCSS complexos.

---

## Funcionalidades

- Extração estruturada de:
  - Variáveis (`$variavel`)
  - Mixins (`@mixin`)
  - Funções (`@function`)
  - Placeholders (`%placeholder`)
  - Regras aninhadas (`nav ul`, `a:hover`)
  - Loops (`@for`, com corpo e expressões)

- Suporte a:
  - Entrada via arquivo SCSS
  - Entrada via string SCSS (ideal para uso em APIs ou ferramentas visuais)
  - Exportação JSON indentada
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
fmt.Println(result.Variables)
```

### ➔ Parseando conteúdo SCSS como `string`

```go
content := `$color: red;\n@mixin test { color: $color; }`
result, err := scss2json.ParseContent(content)
if err != nil {
    panic(err)
}
fmt.Println(result.Mixins[0].Name) // test
```

---

## Estrutura do JSON

```json
{
  "variables": [...],
  "mixins": [...],
  "functions": [...],
  "placeholders": [...],
  "rules": [...],
  "loops": [...]
}
```

---

## Exemplos Avançados

### Uso com API Go

```go
import (
  "net/http"
  "io/ioutil"
  "github.com/LuizFelipeVPCrema/scss2json"
)

func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadAll(r.Body)
    result, _ := scss2json.ParseContent(string(body))
    json.NewEncoder(w).Encode(result)
}
```

### Exportando para disco

```go
scss2json.ExportToJson(result, "saida.json")
```

---


## API

| Função                       | Descrição                                |
|-----------------------------|------------------------------------------|
| `ParseFile(path)`          | Parse de arquivo SCSS                    |
| `ParseContent(content)`    | Parse de string SCSS                     |
| `ExportToJson(result, out)`| Exporta resultado como JSON              |

---

## Tecnologias

- Go 1.21+
- Regex para parsing
- `encoding/json` para exportação
- `flag` para CLI

---

## Contribuições

Contribuições, issues e melhorias são bem-vindas!

1. Forke o projeto
2. Crie sua branch: `git checkout -b minha-feature`
3. Commit: `git commit -m 'feat: minha nova feature'`
4. Push: `git push origin minha-feature`
5. Crie um Pull Request

---

## Licença

Este projeto está licenciado sob a Licença MIT.

MIT © [LuizFelipeVPCrema](https://github.com/LuizFelipeVPCrema)

---

