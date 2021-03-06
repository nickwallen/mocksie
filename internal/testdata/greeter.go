package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type greeter interface {
	SayHello(in io.Reader, out io.Writer) error
	SayGoodbye(in io.Reader, out io.Writer) error
}

type warmGreeter struct {
}

func (w warmGreeter) SayHello(in io.Reader, out io.Writer) error {
	name, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "Hello and welcome, %s", name)
	return err
}

func (w warmGreeter) SayGoodbye(in io.Reader, out io.Writer) error {
	name, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "Goodbye, %s", name)
	return err
}

func main() {
	var greeter greeter

	// Say hello
	greeter = &warmGreeter{}
	err := greeter.SayHello(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("failed to say hello: %v", err)
	}

	// Say goodbye
	err = greeter.SayGoodbye(os.Stdin, os.Stdin)
	if err != nil {
		log.Fatalf("failed to say goodbye: %v", err)
	}
}
