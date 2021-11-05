package testdata

// There are multiple interfaces defined in this file.

type thisOne interface {
	DoThisThing() (string, error)
}

type thatOne interface {
	DoThatThing() (string, error)
}

type anotherOne interface {
	DoAnotherThing() (string, error)
}