package ast

import "goscript/function"

type OptimizeConfig struct{

}

func Optimize(exp Expression, udf map[string]function.Function) (Expression, error) {

	return exp, nil
}

