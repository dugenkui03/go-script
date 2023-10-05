package ast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestGetAllTokens(t *testing.T) {
	for _,exp := range validExpressions {
		tokens, err := getAllTokens(exp)
		assert.Nil(t, err)
		t.Logf("source code:%s, tokens: %+v", exp, tokens)
	}
}

