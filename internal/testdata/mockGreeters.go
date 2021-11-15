package main

import (
	"io"
)

// mockHelloGreeter ia a mock implementation of the helloGreeter interface.
type mockHelloGreeter struct {
	DoSayHello func(in io.Reader, out io.Writer) error
}

// SayHello relies on DoSayHello for defining its behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockHelloGreeter) SayHello(in io.Reader, out io.Writer) error {
	return m.DoSayHello(in, out)
}

// mockGoodbyeGreeter ia a mock implementation of the goodbyeGreeter interface.
type mockGoodbyeGreeter struct {
	DoSayGoodbye func(in io.Reader, out io.Writer) error
}

// SayGoodbye relies on DoSayGoodbye for defining its behavior. If this is causing a panic,
// define DoSayGoodbye within your test case.
func (m *mockGoodbyeGreeter) SayGoodbye(in io.Reader, out io.Writer) error {
	return m.DoSayGoodbye(in, out)
}
