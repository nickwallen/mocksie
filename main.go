package main

import (
	"log"

	"github.com/nickwallen/mocksie/internal/generator"
	"github.com/nickwallen/mocksie/internal/parser"
)

func main() {
	filename := "internal/testdata/greeter.go"

	// Find all interfaces
	parser, err := parser.NewFileParser(filename)
	if err != nil {
		log.Fatalf("%v: %s", err, filename)
	}
	ifaces, err := parser.FindInterfaces()
	if err != nil {
		log.Fatalf("%v: %s", err, filename)
	}
	if len(ifaces) == 0 {
		log.Printf("No interfaces found in %s", filename)
		return
	}

	// Generate a mock for each interface
	gen := generator.NewGenerator()
	for _, iface := range ifaces {

		err := gen.GenerateMock(iface)
		if err != nil {
			log.Fatalf("Failed to generate mock for %s: %v", iface.Name, err)
		}
	}
}
