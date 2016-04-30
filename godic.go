/*
Package godic provides a simple utility to auto-generate idiomatic dependency
injection containers that can benefit from Go's type checking system.
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

func NewGenerator() *Generator {
	return &Generator{
		opts: defaultOpts(),
	}
}

func (g *Generator) SetPackage(p string) *Generator {
	g.opts.Package = p
	return g
}

func (g *Generator) AddImports(l ...string) *Generator {
	g.opts.Imports = append(g.opts.Imports, l...)
	return g
}

func (g *Generator) SetName(n string) *Generator {
	g.opts.Name = n
	return g
}

func (g *Generator) SetType(t string) *Generator {
	g.opts.Type = t
	return g
}

func (g *Generator) AddDependency(n, t, f string) *Generator {
	g.deps = append(g.deps, dep{n, t, f})
	return g
}

func (g *Generator) Generate() {
	g.Fgenerate(os.Stdout)
}

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
