package testdata

import "fmt"

type greeter interface {
	SayHello(name string) (string, error)
	SayGoodbye(name string) (string, error)
}

type warmGreeter struct {
}

func (w warmGreeter) SayHello(name string) (string, error) {
	if len(name) == 0 {
		return "", fmt.Errorf("no name given")
	}
	return fmt.Sprintf("Hello, %s! It is so nice to meet you.", name), nil
}

func (w warmGreeter) SayGoodbye(name string) (string, error) {
	if len(name) == 0 {
		return "", fmt.Errorf("no name given")
	}
	return fmt.Sprintf("Goodbye, %s! It was so nice to meet you.", name), nil
}