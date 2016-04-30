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

- Factory/Builder: to abstract the complexity of producing an instance;
- Service Locator: to easily locate produced instances and hold only one
  instance of each dependency;
- Lazy loading: to prevent instantiation of dependencies not needed at run time.

# Status

`godic` is under development and should be considered experimental. However,
due to its simplicity and its extremely narrow scope, it can be used in
production grade projects *as long as its version is locked*. There will be
no guarantees of API stability at this point. Since this is a code generation
tool, the impact of changes to the generated code should be minimal, though.

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
an executable that will produce the DIC code. Run this executable, preferably
(but not necessarily) with `go generate` to produce the code for the
DI container.

Read the [documentation](https://godoc.org/github.com/carlosdavidepto/godic)
for the API reference and check the `examples` directory for a sample of
how the generation code and the usage (i.e app/project) code tie together.
