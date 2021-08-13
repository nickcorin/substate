package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/nickcorin/substate"
)

var (
	typeName = flag.String("type", "Substate", "The name of the interface")
	outFile  = flag.String("outFile", "substate_gen.go", "Output file to write")
)

func main() {
	flag.Parse()

	if *typeName == "" {
		log.Fatalf("type cannot be empty")
	}

	if !strings.HasSuffix(*outFile, ".go") {
		log.Fatalf("outFile must be a .go file")
	}

	srcFile := os.Getenv("GOFILE")

	if err := substate.Generate(srcFile, *outFile, *typeName); err != nil {
		log.Fatalf("generate: %s", err)
	}
}
