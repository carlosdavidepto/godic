package main

import "fmt"

type A struct{}

type B struct {
  a *A
}

type C struct {
  a *A
  b *B
}

type D struct {
  b *B
  c *C
}

//go:generate sh -c "go run ../digen/main.go > dicontainer.go && go fmt dicontainer.go"

func main() {
	dic := &DIContainer{}

	d := dic.D()

	fmt.Println(d)
}
