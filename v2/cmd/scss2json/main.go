package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/LuizFelipeVPCrema/scss2json/v2/pkg/scss2json"
)

func main() {
	input := flag.String("i", "", "Input SCSS file")
	output := flag.String("o", "output.json", "JSON output file")
	flag.Parse()

	if *input == "" {
		log.Fatal("You must provide an input file with -i")
	}

	result, err := scss2json.ParseFile(*input)
	if err != nil {
		log.Fatalf("Error parsing: %v", err)
	}

	err = scss2json.ExportToJson(result, *output)
	if err != nil {
		log.Fatalf("Error when exporting: %v", err)
	}

	fmt.Println("Successfully exported to", *output)
}
