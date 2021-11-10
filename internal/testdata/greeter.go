package testdata

import (
	"fmt"
	"io"
)

type greeter interface {
	SayHello(name string, out io.Writer) (string, error)
	SayGoodbye(name string) (string, error)
}

type warmGreeter struct {
}

func (w warmGreeter) SayHello(name string) (string, error) {
	return fmt.Sprintf("Hello, %s", name), nil
}

func (w warmGreeter) SayGoodbye(name string) (string, error) {
	return fmt.Sprintf("Goodbye, %s", name), nil
}
