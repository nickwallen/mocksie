package main

import (
	"fmt"
	"github.com/nickwallen/farce/internal"
	"log"
)

func main() {
	filename := "internal/testdata/greeter.go"
	parser, err := parser.NewFileParser(filename)
	if err != nil {
		log.Fatalf("%v: %s", err, filename)
	}

	ifaces, err := parser.FindInterfaces()
	if err != nil {
		log.Fatalf("%v: %s", err, filename)
	}
	for _, iface := range ifaces {
		fmt.Printf("Found interface '%s' with method(s) %s", iface.Name, iface.Methods)
	}

	// TODO generate mock for the interfaces
}
