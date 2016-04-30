/*
Package godic provides a simple utility to auto-generate idiomatic dependency
injection containers that can benefit from Go's type checking system.

Example usage:

		package main

		import "github.com/carlosdavidepto/godic"

		func main() {

			// get a new generator with default values
			g := godic.NewGenerator()
			g.AddDependency("configOption", "int", "{ return 10 }")
			g.Generate()
		}

This generation program would yield the following code:

		package main

		type Container struct{
			configOption int
		}

		func (c *Container) NewConfigOption() int { return 10 }

		func (c *Container) ConfigOption() int {
			if c.configOption == nil {
				c.configOption = c.NewConfigOption()
			}
			return c.configOption
		}

See the code in the examples directory for a more complete usage example.
*/
package godic

import (
	"io"
	"os"
	"text/template"
	"unicode"
)

type (
	opts struct {
		Package string
		Imports []string
		Name    string
		Type    string
	}

	dep struct {
		Name string
		Type string
		Func string
	}

	// Generator implements a static Go code generation mechanism for dependency
	// injection containers.
	Generator struct {
		opts opts
		deps []dep
	}
)

func defaultOpts() opts {
	return opts{
		Package: "main",
		Name:    "c",
		Type:    "Container",
	}
}

// NewGenerator returns a new instance of a Generator, preloaded with default
// parameters.
func NewGenerator() *Generator {
	return &Generator{
		opts: defaultOpts(),
	}
}

// SetPackage overrides the package that the generated code file will belong
// to. Default is main.
func (g *Generator) SetPackage(p string) *Generator {
	g.opts.Package = p
	return g
}

// AddImports registers the import statements to be added to the generated file
func (g *Generator) AddImports(l ...string) *Generator {
	g.opts.Imports = append(g.opts.Imports, l...)
	return g
}

// Set name overrides the receiver variable name for the container methods in
// the generated code. Default is "c".
func (g *Generator) SetName(n string) *Generator {
	g.opts.Name = n
	return g
}

// Set type overrides the type of the container struct in the generated code.
// Default is "Container".
func (g *Generator) SetType(t string) *Generator {
	g.opts.Type = t
	return g
}

// AddDependency registers a new dependency for the generated code.
//
// Each dependency that is registered will create:
//
// - A field in the container struct;
//
// - A "create" method, which builds new instances of a dependency;
//
// - A "lookup" method, which resolves (calling the create method, if needed)
// an instance of a dependency.
//
// AddDependency receives three parameters:
//
// - The name of the dependency, which will be used to determine the name of
// the struct field, lookup method and create method;
//
// - The type of the dependency (for all of the above);
//
// - The body of the create function (which can reference the container to
// fetch other dependencies);
func (g *Generator) AddDependency(n, t, f string) *Generator {
	g.deps = append(g.deps, dep{n, t, f})
	return g
}

// Generate writes the generated code to the standard output.
func (g *Generator) Generate() {
	g.Fgenerate(os.Stdout)
}

// Fgenerate writes the generated code to the io.Writer specified as parameter.
func (g *Generator) Fgenerate(w io.Writer) {
	tpkg := templateBuilder("package", `package {{ .Package }}{{"\n\n"}}`)

	timp := templateBuilder("imports", `import {{ if eq (len .Imports) 1 }}"{{ index .Imports 0 }}"{{else}}(
{{ range .Imports }}{{"\t"}}"{{ . }}"
{{ end }}){{ end }}{{"\n\n"}}`)

	ttyp := templateBuilder("type", `type {{ .Opts.Type }} struct{{"{"}}
{{- if gt (len .Deps) 0 -}}{{ "\n" }}
  {{- range .Deps -}}
    {{- "\t" }}{{ .Name }} {{ .Type -}}{{ "\n" }}
  {{- end -}}
{{- end -}}
}{{"\n\n"}}`)

	tdep := templateBuilder("deps", `{{- $cname := .Opts.Name -}}
{{- $ctype := .Opts.Type -}}
{{- range .Deps -}}
func ({{ $cname }} *{{ $ctype }}) New{{ .Name | ucfirst }}() {{ .Type }} {{ .Func }}{{ "\n\n" }}
func ({{ $cname }} *{{ $ctype }}) {{ .Name | ucfirst }}() {{ .Type }} {
{{ "\t" }}if {{ $cname }}.{{ .Name | lcfirst }} == nil {
{{ "\t\t" }}{{ $cname }}.{{ .Name | lcfirst }} = {{ $cname }}.New{{ .Name | ucfirst }}()
{{ "\t" }}}
{{ "\t" }}return {{ $cname }}.{{ .Name | lcfirst }}
}{{ "\n\n" }}
{{- end -}}`)

	err := tpkg.Execute(w, g.opts)

	if err != nil {
		panic(err)
	}

	if len(g.opts.Imports) > 0 {
		err = timp.Execute(w, g.opts)

		if err != nil {
			panic(err)
		}
	}

	gWithPublicFields := struct {
		Opts opts
		Deps []dep
	}{
		g.opts,
		g.deps,
	}

	err = ttyp.Execute(w, gWithPublicFields)

	if err != nil {
		panic(err)
	}

	if len(g.deps) > 0 {
		err = tdep.Execute(w, gWithPublicFields)

		if err != nil {
			panic(err)
		}
	}
}

func templateBuilder(name, tpl string) *template.Template {
	t, err := template.
		New(name).
		Funcs(template.FuncMap{
			"ucfirst": ucfirst,
			"lcfirst": lcfirst,
		}).
		Parse(tpl)

	if err != nil {
		panic(err)
	}

	return t
}

func ucfirst(s string) string {
	tmp := []rune(s)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

func lcfirst(s string) string {
	tmp := []rune(s)
	tmp[0] = unicode.ToLower(tmp[0])
	return string(tmp)
}
