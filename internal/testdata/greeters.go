package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type helloGreeter interface {
	SayHello(in io.Reader, out io.Writer) error
}

type warmHelloGreeter struct {
}

func (w warmHelloGreeter) SayHello(in io.Reader, out io.Writer) error {
	name, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "Hello and welcome, %s", name)
	return err
}

type goodbyeGreeter interface {
	SayGoodbye(in io.Reader, out io.Writer) error
}

type warmGoodbyeGreeter struct {
}

func (w warmGoodbyeGreeter) SayGoodbye(in io.Reader, out io.Writer) error {
	name, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "Goodbye, %s", name)
	return err
}

func main() {
	// Say hello
	var hello helloGreeter
	hello = &warmHelloGreeter{}
	err := hello.SayHello(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("failed to say hello: %v", err)
	}

	// Say goodbye
	var goodbye goodbyeGreeter
	goodbye = &warmGoodbyeGreeter{}
	err = goodbye.SayGoodbye(os.Stdin, os.Stdin)
	if err != nil {
		log.Fatalf("failed to say goodbye: %v", err)
	}
}
