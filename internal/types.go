package parser

// Interface an interface that will need to be mocked.
type Interface struct {
	Name    string
	Methods []Method
}

// Method There are one or more Methods within every Interface.
type Method struct {
	Name    string
	Params  []Param
	Results []Result
}

// Param A Method has zero or more call parameters.
type Param struct {
	Name string
	Type string
}

// Result A Method has zero or more results that are returned.
type Result struct {
	Name string
	Type string
}
