package vm

import (
	"errors"
	"goscript/function"
)

func (vm *VM) RemoveFunc(name string) {
	if vm == nil || len(vm.expressionCache) == 0 {
		return
	}

	delete(vm.expressionCache, name)
}

func (vm *VM) RegisterFunc0(name string, allowFold bool, f func() (interface{}, error)) error {
	return vm.registerFuncN(name, allowFold, 0, f)
}

func (vm *VM) RegisterFunc1(name string, allowFold bool, f func(arg Value) (interface{}, error)) error {
	return vm.registerFuncN(name, allowFold, 1, f)
}

func (vm *VM) RegisterFunc2(name string, allowFold bool, f func(arg1, arg2 Value) (interface{}, error)) error {
	return vm.registerFuncN(name, allowFold, 2, f)
}

func (vm *VM) RegisterFunc3(name string, allowFold bool, f func(arg1, arg2, arg3 Value) (interface{}, error)) error {
	return vm.registerFuncN(name, allowFold, 3, f)
}

func (vm *VM) registerFuncN(name string, allowFold bool, i int, f interface{}) error {
	if f == nil {
		return errors.New("function should not be nil")
	}

	if vm == nil {
		return errors.New("vm is nil ptr")
	}

	if vm.funcByName == nil{
		vm.funcByName = make(map[string]function.Function)
	}

	if _, ok := vm.funcByName[name]; ok {
		return errors.New("already register function named " + name)
	}

	vm.funcByName[name] = function.NewFunction(name, f, i, allowFold)

	return nil
}
