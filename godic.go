/*
Package godic provides a simple utility to auto-generate idiomatic dependency
injection containers that can benefit from Go's type checking system.
*/
package godic

import (
	"os"
	"text/template"
	"unicode"
)

// Generator defines what a dependency injection container's code is going to
// look like. The properties of an instance of Generator directly map to
// source code entities in the produced code.
type Generator struct {

	// Package specifies which package the container is going
	// to belong to.
	Package string

	// Imports is the list of imports the generated code will need.
	Imports []string

	// Type is the Go type of the generated container.
	Type string

	// Deps is the list of dependency definitions that will make up the
	// lookup/create methods and the properties of the container.
	Deps []Dep
}

// Dep is the definition for a dependency.
type Dep struct {
	// Type is the Go type of the object to be created when the lookup/create
	// function is invoked.
	Type string

	// Name is used in constructing the name of the struct field, lookup and
	// create function. As an example, if the Name is "dependency", the struct
	// field name will be "dependency", the create function will be named
	// "NewDependency" and the lookup function will be named just like a getter,
	// i.e. "Dependency" (note that lookup/create are public, but the struct
	// field is private).
	Name string

	// Func is the text for the Go code that will be the body of the builder
	// function. Providing only the body of the function only avoids repetition;
	// the lookup function always has the same structure, and the signatures
	// for both the lookup and create functions are similar and easy to generate
	// automatically.
	Func string
}

// Generate renders the code template with the Generator properties to create
// the final DIC code which will be output to STDOUT.
func (c *Generator) Generate() {
	t, err := template.
		New("generate").
		Funcs(template.FuncMap{
			"ucfirst": ucfirst,
			"lcfirst": lcfirst,
		}).
		Parse(tplstr)

	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, c)

	if err != nil {
		panic(err)
	}
}

var tplstr = `package {{ .Package }}

{{ range .Imports }}import "{{.}}"
{{ end }}

type {{ .Type }} struct {
{{ range .Deps }} {{ .Name | lcfirst }} {{ .Type }}
{{ end }}{{ "}" }}

{{ $ctype := .Type }}
{{ range .Deps }}func (c *{{ $ctype }}) New{{ .Name | ucfirst }}() {{ .Type }} {{ .Func }}

func (c *{{ $ctype }}) {{ .Name | ucfirst }}() {{ .Type }} {
  if c.{{ .Name | lcfirst }} == nil {
    c.{{ .Name | lcfirst }} = c.New{{ .Name | ucfirst }}()
  }

  return c.{{ .Name | lcfirst }}
}

{{ end }}
`

// ucfirst is an internal utility function to make the first character of a
// string upper case (i.e. an exported Go name)
func ucfirst(s string) string {
	tmp := []rune(s)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

// lcfirst is an internal utility function to make the first character of a
// string lower case (i.e. a non-exported Go name)
func lcfirst(s string) string {
	tmp := []rune(s)
	tmp[0] = unicode.ToLower(tmp[0])
	return string(tmp)
}
