package main

import (
	"fmt"

	"github.com/dop251/goja"
)

func main() {
	vm := goja.New()
	v, err := vm.RunString(`2+2`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T, %#v\n", v, v)
}
