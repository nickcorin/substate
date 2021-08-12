package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/nickcorin/substate"
)

var (
	outFile = flag.String("outFile", "substate_gen.go", "Output file")
)

func main() {
	flag.Parse()

	if !strings.HasSuffix(*outFile, ".go") {
		log.Fatalf("outFile must be a .go file")
	}

	srcFile := os.Getenv("GOFILE")

	if err := substate.Generate(srcFile, *outFile); err != nil {
		log.Fatalf("generate: %s", err)
	}
}
