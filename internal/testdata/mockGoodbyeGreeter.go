package main

import (
	"io"
)

// mockGoodbyeGreeter ia a mock implementation of the goodbyeGreeter interface.
type mockGoodbyeGreeter struct {
	DoSayGoodbye func(in io.Reader, out io.Writer) error
}

// SayGoodbye relies on DoSayGoodbye for defining its behavior. If this is causing a panic,
// define DoSayGoodbye within your test case.
func (m *mockGoodbyeGreeter) SayGoodbye(in io.Reader, out io.Writer) error {
	return m.DoSayGoodbye(in, out)
}
