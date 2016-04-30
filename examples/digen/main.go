package main

import "github.com/carlosdavidepto/godic"

func main() {
	g := godic.NewGenerator()

	// we only need to set the fields which are different from the defaults
	g.
		AddImports("fmt", "os").
		SetName("cnt").
		SetType("DIContainer").

		// and obviously, we need to register the dependencies
		AddDependency("a", "*A", `{
			fmt.Fprintln(os.Stdout, "creating A...")
			return &A{}
		}`).
		AddDependency("b", "*B", `{
			a := cnt.A()
			fmt.Fprintln(os.Stdout, "creating B...")
			return &B{a}
		}`).
		AddDependency("c", "*C", `{
			a := cnt.A()
			b := cnt.B()
			fmt.Fprintln(os.Stdout, "creating C...")
			return &C{a, b}
		}`).
		AddDependency("d", "*D", `{
			b := cnt.B()
			c := cnt.C()
			fmt.Fprintln(os.Stdout, "creating D...")
			return &D{b, c}
		}`)

	// output to stdout
	g.Generate()
}
