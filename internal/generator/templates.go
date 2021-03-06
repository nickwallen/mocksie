package generator

const (
	// baseTemplate defines how the mock implementation is generated.
	baseTemplate = `
package {{ .Package }}

{{ template "imports" . }}

// mock{{ .Name | title }} ia a mock implementation of the {{ .Name }} interface.
type mock{{ .Name | title }} struct {
{{- range  .Methods }}
    Do{{ .Name }} func ({{ template "declare-params" . }}) {{ template "results" . }}
{{- end }}
}
{{ template "methods" . -}}
`
	// importsTemplate defines how the imports are generated.
	importsTemplate = `
{{- if gt (len .Imports) 0 -}}
import (
{{- range  .Imports }}
    "{{ .Path }}"
{{- end }}
)
{{- end -}}
`

	// methodsTemplate defines how the methods of the mock implementation are generated.
	methodsTemplate = `
{{- range .Methods }}
// {{ .Name }} relies on Do{{ .Name }} for defining its behavior. If this is causing a panic,
// define Do{{ .Name }} within your test case.
func (m *mock{{ $.Name | title }}) {{ .Name }}({{ template "declare-params" . }}) {{ template "results" . }} {
    {{ if gt (len .Results) 0 }}return {{ end }}m.Do{{ .Name }}({{ template "use-params" . }})
}
{{ end }}
`

	// declareParamsTemplate defines how the method parameters of the mock implementation are declared.
	declareParamsTemplate = `
{{- range $index, $param := .Params -}}
{{ if $index }}, {{ end }}{{ .Name }} {{ .Type }}
{{- end -}}
`

	// useParamsTemplate defines how the method parameters of the mock implementation are called.
	useParamsTemplate = `
{{- range $index, $param := .Params -}}
{{ if $index }}, {{ end }}{{ .Name }}
{{- end -}}
`

	// resultsTemplate defines how the method results of the mock implementation are generated.
	resultsTemplate = `
{{- if gt (len .Results) 1 -}} ( {{- end -}}
{{- range $index, $param := .Results -}}
{{- if $index -}}, {{ end -}}
{{- if gt (len .Name) 0 -}}
{{- .Name }} {{ .Type -}}
{{- else -}}
{{- .Type -}}
{{- end -}}
{{- end -}}
{{- if gt (len .Results) 1 -}} ) {{- end -}}
`
)
