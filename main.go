package main

import (
	"fmt"
	"log"

	"github.com/nickwallen/farce/internal"
)

func main() {
	filename := "main.go"
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
	for _, iface := range ifaces {
		fmt.Printf("Found '%s' with method(s) %s", iface.Name, iface.Methods)
	}

	// TODO generate mock for the interfaces
}
