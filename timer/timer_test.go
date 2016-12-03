package timer

import (
	"testing"
	"time"

	"github.com/dop251/goja"
)

func TestSetTimeout(t *testing.T) {
	registry := NewRegistry()
	vm := goja.New()
	registry.Enable(vm)

	start := time.Now()

	vm.Set("done", func(delay int) {
		et := time.Now().Sub(start)
		t.Logf("delay = %d, elapsed time = %v", delay, et)
	})

	_, err := vm.RunString(`
    setTimeout(function(){
        done(1000)
    }, 1000);


    setTimeout(function(){
        done(500)
    }, 500);
    `)

	if err != nil {
		t.Fatal(err)
	}

	registry.Wait()
}

func TestClearTimer(t *testing.T) {
	registry := NewRegistry()
	vm := goja.New()
	registry.Enable(vm)

	start := time.Now()

	vm.Set("done", func(delay int) {
		et := time.Now().Sub(start)
		t.Fatalf("cancell error: delay = %d elapsed time = %v", delay, et)
	})

	vm.Set("log", t.Log)

	_, err := vm.RunString(`
    var timer = setTimeout(function(){
        done(500)
    }, 500);

    clearTimeout(timer)
    `)

	if err != nil {
		t.Fatal(err)
	}

	registry.Wait()
}

func TestSetInterval(t *testing.T) {
	registry := NewRegistry()
	vm := goja.New()
	registry.Enable(vm)

	start := time.Now()

	vm.Set("done", func(delay int) {
		et := time.Now().Sub(start)
		t.Fatalf("cancell error: delay = %d elapsed time = %v", delay, et)
	})

	vm.Set("log", t.Log)

	_, err := vm.RunString(`
    var i = 0
    var timer = setInterval(function(){
        i++
        if(i >= 3) {
            clearTimeout(timer)
        }
        log(50)
    }, 50);
    `)

	if err != nil {
		t.Fatal(err)
	}

	registry.Wait()
}
