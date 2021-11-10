package main

import (
	"bufio"
	"fmt"
	"io"
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
	greeter = &warmGreeter{}
	greeter.SayHello(os.Stdin, os.Stdout)
	greeter.SayGoodbye(os.Stdin, os.Stdin)
}
