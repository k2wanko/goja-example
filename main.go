package main

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func main() {
	registry := new(require.Registry)

	vm := goja.New()
	rm := registry.Enable(vm)

	var v goja.Value
	if sum, err := rm.Require("sum.js"); err != nil {
		panic(err)
	} else {
		if sum, ok := goja.AssertFunction(sum); ok {
			v, err = sum(goja.Undefined(), vm.ToValue(6), vm.ToValue(7))
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Printf("%T, %#v\n", v, v)
}
