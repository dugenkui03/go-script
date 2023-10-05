package config

import (
	"goscript/ast"
	"goscript/function"
)


type Option func(config *Config)

type Config struct {
	expCache   map[string]ast.Expression
	funcByName map[string]function.Function
}

func (c *Config) FuncByName() map[string]function.Function {
	if c == nil || len(c.funcByName) == 0 {
		return map[string]function.Function{}
	}

	copyValue := make(map[string]function.Function)
	for k, v := range c.funcByName {
		copyValue[k] = v
	}

	return copyValue
}

func (c *Config) ExpCache() map[string]ast.Expression {
	if c == nil || len(c.expCache) == 0 {
		return map[string]ast.Expression{}
	}

	copyValue := make(map[string]ast.Expression)
	for k, v := range c.expCache {
		copyValue[k] = v
	}

	return copyValue
}
