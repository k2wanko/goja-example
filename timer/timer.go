package timer

import (
	"time"

	"sync"

	"github.com/dop251/goja"
)

type (
	Registry struct {
		sync.Mutex
		wg     *sync.WaitGroup
		vm     *goja.Runtime
		timers map[*timer]*timer
	}

	timer struct {
		timer    *time.Timer
		duration time.Duration
		interval bool
		call     goja.Callable
	}
)

func NewRegistry() *Registry {
	return &Registry{
		wg:     new(sync.WaitGroup),
		timers: map[*timer]*timer{},
	}
}

func (r *Registry) newTimer(call goja.Callable, delay int64, interval bool) *timer {
	t := &timer{
		call:     call,
		duration: time.Duration(delay) * time.Millisecond,
		interval: interval,
	}
	r.Lock()
	r.timers[t] = t
	r.Unlock()

	r.wg.Add(1)
	t.timer = time.AfterFunc(t.duration, func() {
		r.call(t)
	})

	return t
}

func (r *Registry) call(t *timer) {
	_, err := t.call(nil)
	if err != nil {
		panic(err)
	}

	if t.interval {
		t.timer.Reset(t.duration)
	} else {
		r.clearTimer(t)
	}
}

func (r *Registry) setTimeout(c goja.FunctionCall) goja.Value {
	call, ok := goja.AssertFunction(c.Argument(0))

	if !ok {
		panic("argument 1 is not function")
	}

	delay := c.Argument(1).ToInteger()

	return r.vm.ToValue(r.newTimer(call, delay, false))
}

func (r *Registry) setInterval(c goja.FunctionCall) goja.Value {
	call, ok := goja.AssertFunction(c.Argument(0))

	if !ok {
		panic("argument 1 is not function")
	}

	delay := c.Argument(1).ToInteger()

	return r.vm.ToValue(r.newTimer(call, delay, true))
}

func (r *Registry) clearTimer(t *timer) {
	timer, ok := r.timers[t]
	if !ok {
		return
	}

	timer.timer.Stop()
	delete(r.timers, timer)
	r.wg.Done()
}

func (r *Registry) clearTimeout(c goja.FunctionCall) goja.Value {
	t, ok := c.Argument(0).Export().(*timer)
	if !ok {
		return goja.Undefined()
	}

	r.clearTimer(t)

	return goja.Undefined()
}

func (r *Registry) Enable(vm *goja.Runtime) {
	r.vm = vm
	vm.Set("setTimeout", r.setTimeout)
	vm.Set("setInterval", r.setInterval)
	vm.Set("clearTimeout", r.clearTimeout)
	vm.Set("clearInterval", r.clearTimeout)
}

func (r *Registry) Wait() {
	r.Lock()
	defer r.Unlock()
	if len(r.timers) <= 0 {
		return
	}
	r.wg.Wait()
}
