package main

import "github.com/carlosdavidepto/godic"

type Dep godic.Dep

func main() {
	c := &godic.Generator{

		// generate our DIC belonging to package "main"
		Package: "main",

		// since some of the functions in the DIC use the fmt package,
		// the generated file will need to import it
		// this would normally be used to import the packages that make
		// the dependencies unless (as is the case with this example)
		// they are part of the same package
		Imports: []string{"fmt"},

		// our DIC will be created as a "type DIContainer struct {...}"
		Type: "DIContainer",

		// each dependency is defined as a struct with three fields:
		// type, name, and the body of the builder function
		// additionally, each dependency creates three elements in generated
		// container code:
		// - a private field in the container struct
		// - a public container method that creates new instances of that
		//   dependency
		// - a public container method that acts like a singleton builder
		//   for that dependency
		Deps: []godic.Dep{

			// example: this dependency definition will yield in the generated code
			// the following elements:
			//
			//    type DIContainer struct {
			//        ...
			//        a *A
			//        ...
			//    }
			//
			//    func (c *DIContainer) NewA() *A {
			//        fmt.Println("creating A...")
			//        return &A{}
			//    }
			//
			//    func (c *DIContainer) A() *A {
			//        if c.a == nil {
			//            c.a = c.NewA()
			//        }
			//
			//        return c.a
			//    }
			godic.Dep{"*A", "a", `{
        fmt.Println("creating A...")
        return &A{}
      }`},

			godic.Dep{"*B", "b", `{
        a := c.A()
        fmt.Println("creating B...")
        return &B{a}
      }`},

			godic.Dep{"*C", "c", `{
        a := c.A()
        b := c.B()
        fmt.Println("creating C...")
        return &C{a,b}
      }`},

			godic.Dep{"*D", "d", `{
        b := c.B()
        depc := c.C()
        fmt.Println("creating D...")
        return &D{b, depc}
      }`},
		},
	}

	c.Generate()
}
