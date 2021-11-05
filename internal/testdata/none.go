package testdata

// There are no interfaces defined in this file.

type coldGreeter struct {
}

func (c coldGreeter) SayHello(_ string) (string, error) {
	return "Hi.", nil
}

func (c coldGreeter) SayGoodbye(_ string) (string, error) {
	return "Bye.", nil
}