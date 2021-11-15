package main

import (
	"io"
)

// mockGreeter ia a mock implementation of the greeter interface.
type mockGreeter struct {
	DoSayHello   func(in io.Reader, out io.Writer) error
	DoSayGoodbye func(in io.Reader, out io.Writer) error
}

// SayHello relies on DoSayHello for defining its behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(in io.Reader, out io.Writer) error {
	return m.DoSayHello(in, out)
}

// SayGoodbye relies on DoSayGoodbye for defining its behavior. If this is causing a panic,
// define DoSayGoodbye within your test case.
func (m *mockGreeter) SayGoodbye(in io.Reader, out io.Writer) error {
	return m.DoSayGoodbye(in, out)
}
