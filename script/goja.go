package script

import (
	"errors"
	"isp-script-service/conf"
	"isp-script-service/router"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/integration-system/isp-lib/v2/config"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return initVm()
	},
}

func initVm() *goja.Runtime {
	vm := goja.New()
	vm.Set("invoke", router.Invoke)

	return vm
}

type Goja struct{}

func (*Goja) Compile(script string) (Script, error) {
	prog, err := goja.Compile("script", script, false)
	if err != nil {
		return nil, err
	}

	return &GojaProgram{prog: prog}, nil
}

var errUnknownEngine = errors.New("unknown engine")

func (*Goja) Execute(program Script, arg interface{}) (interface{}, error) {
	value, ok := program.Src().(*goja.Program)

	if !ok {
		return nil, errUnknownEngine
	}

	vm := pool.Get().(*goja.Runtime)
	cfg := config.GetRemote().(*conf.RemoteConfig)
	awaitTime := time.Duration(cfg.ScriptExecutionTimeoutMs) * time.Millisecond
	t := time.AfterFunc(awaitTime, func() {
		vm.Interrupt("execution timeout")
	})
	defer func() {
		t.Stop()
		vm.ClearInterrupt()
		pool.Put(vm)
	}()

	vm.Set("arg", arg)
	res, err := vm.RunProgram(value)
	if err != nil {
		return nil, err
	}

	return res.Export(), nil
}

type GojaProgram struct {
	prog *goja.Program
}

func (a *GojaProgram) Src() interface{} {
	return a.prog
}
