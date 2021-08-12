package substate

const tpl = `// Code generated by gensubstate at {{.Source}}; DO NOT EDIT.

package {{.Package}}

{{if .Imports}}
import (
	{{- range $import := .Imports}}
	{{$import}}
	{{- end}}
)
{{end}}

// NewSubstateForTesting returns an implementation of Substate which can be used
// for testing.
func NewSubstateForTesting(_ *testing.TB, injectors...Injector) *substate {
	var s substate

	for _, injector := range injectors {
		injector.Inject(&s)
	}

	return &s
}

type Injector interface {
	Inject(*substate)
}

// InjectorFunc defines a convenience type making it easy to implement
// Injectors.
type InjectorFunc func(*substate)

// Inject implements the Injector interface.
func (fn InjectorFunc) Inject(s *substate) {
	fn(s)
}
{{range $index, $field := .Fields}}
// With{{$field.Method}} returns an Injector which sets the {{$field.Name}} on substate.
func With{{$field.Method}}({{$field.Name}} {{$field.Type}}) InjectorFunc {
	return func(s *substate) {
		s.{{$field.Name}} = {{$field.Name}}
	}
}
{{end}}
type substate struct {
{{- range $index, $field := .Fields}}
{{$field.Name}} {{$field.Type}}
{{- end}}
}
{{range $index, $field:= .Fields}}
// {{$field.Method}} implements the Substate interface.
func (s *substate) {{$field.Method}}{{$field.Params}} {{$field.Results}} {
	return s.{{$field.Name}}
}
{{- end}}
`