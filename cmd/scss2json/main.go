package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/LuizFelipeVPCrema/scss2json/pkg/scss2json"
)

func main() {
	input := flag.String("i", "", "Arquivo SCSS de entrada")
	output := flag.String("o", "output.json", "Arquivo de saída JSON")
	flag.Parse()

	if *input == "" {
		log.Fatal("Você deve fornecer um arquivo de entrada com -i")
	}

	result, err := scss2json.ParseFile(*input)
	if err != nil {
		log.Fatalf("Erro ao parsear: %v", err)
	}

	err = scss2json.ExportToJson(result, *output)
	if err != nil {
		log.Fatalf("Erro ao exportar: %v", err)
	}

	fmt.Println("Exportado com sucesso para", *output)
}
