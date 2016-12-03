package main

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/k2wanko/goja-example/timer"
)

func main() {
	registry := timer.NewRegistry()

	vm := goja.New()
	registry.Enable(vm)

	vm.Set("debug", func(c goja.FunctionCall) goja.Value {
		fmt.Printf("%v\n", c.Arguments[0])
		return goja.Null()
	})

	_, err := vm.RunString(`
        setTimeout(function() {
            debug('timeout')
        }, 1000)
    `)

	if err != nil {
		panic(err)
	}

	registry.Wait()
}
