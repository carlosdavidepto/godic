# godic

Easily generate idiomatic, lazy Dependency Injection Containers for your Go
app/project.

# Why? Aren't there enough DI frameworks in the Go ecosystem?

`godic` is not a framework. It is a code generation tool, particularly inclined
to be used with `go generate`. It is not meant to do any magic shenanigans with
your code. It is meant to make writing dependency injection containers easier,
leaving the developer in control of how to wire the dependencies and their
dependents.

`godic` requires only the standard library, and the generated code makes no
use of reflection, since it does not try to figure out where each dependency
should be injected.

In terms of engineering, the generated code combines the only the minimum
number of patterns necessary for a useful DIC:

- Factory/Builder: to abstract the complexity of producing an instance
- Service Locator: to easily locate produced instances

# Status

`godic` is under development and should be considered experimental. However,
due to its simplicity and its extremely narrow scope, it is ok to use for
production use *as long as its version is locked*. There will be
no guarantees of API stability at this point.

# Features/Objectives

- Generated code is simple, clean and idiomatic;
- Minimalist: unnecessary bells and whistles will be avoided;
- Generated DI containers will benefit from static type checking, along with
  all the other goodies that the Go compiler offers;
- Declare your project dependencies and their configuration in a language you
  already know (Go);
- All dependencies made available to the rest of the code via an obvious, easy
  to refactor interface (i.e. a getter method on the DIC instance);
- Lazy instantiation: no dependency instantiation occurs until the dependency
  is needed for the first time.

# Usage

`godic` is designed to be used with `go generate`, but it's not mandatory.

To use `godic` to generate a DIC, simply create a package in your project with
an executable that will produce the DIC code:

    package main

    import "github.com/carlosdavidepto/godic"

    func main() {

      c := &godic.Generator{
        Package: "mypackage",
        Type:    "DIContainer",
        Deps: []godic.Dep{
          godic.Dep{"*ADependency", "aDependency", `{
            return &ADependency{}
          }`},
        },
      }

      c.Generate()
    }

The program above, when executed, will output the following Go code to `STDOUT`:

    package mypackage

    type DIContainer struct {
      aDependency *ADependency
    }

    func (c *DIContainer) NewADependency() *ADependency {
      return &ADependency{}
    }

    func (c *DIContainer) ADependency() *ADependency {
      if c.aDependency == nil {
        c.aDependency = c.NewADependency()
      }

      return c.aDependency
    }

Then proceed as normal in the rest of the code:

    package otherpackage

    import "mypackage"

    func main() {

      dic := &mypackage.DIC{}

      o := &NewObjectWithADependency{
        dependency: dic.ADependency()
      }

      p := &NewOtherObject{
        dependency: dic.ADependency() // here the same instance is reused
      }
    }

Check the `examples` directory for a more complete example of how it all ties
together.
