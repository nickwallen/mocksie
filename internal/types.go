package mocksie

// Interface is an interface that will need to be mocked.
type Interface struct {
	Name    string
	Package string
	Methods []Method
}

// Method is a method that is part of an Interface. There are one or more methods
// within an Interface.
type Method struct {
	Name    string
	Params  []Param
	Results []Result
}

// Param is a parameter to a Method call. A Method has zero or more call parameters.
type Param struct {
	Name string
	Type string
}

// Result is the result that is returned by a Method call. A Method has zero or
// more results that can be either named or unnamed.
type Result struct {
	Name string
	Type string
}
