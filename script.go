package goscript

import (
	"goscript/vm"
)

func NewVm() *vm.VM {
	return vm.NewVM()
}

func Eval(exp string, env map[string]interface{}) (*vm.Value, error) {
	return vm.Eval(exp, env)
}
