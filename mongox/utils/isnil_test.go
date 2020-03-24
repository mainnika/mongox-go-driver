package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {

	testvalues := []struct {
		i     interface{}
		isnil bool
	}{
		{nil, true},
		{(*string)(nil), true},
		{([]string)(nil), true},
		{(map[string]string)(nil), true},
		{(func() bool)(nil), true},
		{(chan func() bool)(nil), true},
		{"", true},
		{0, true},
		{append(([]string)(nil), ""), false},
		{[]string{}, false},
		{1, false},
		{"1", false},
	}

	for _, tt := range testvalues {
		assert.Equal(t, tt.isnil, IsNil(tt.i))
	}
}
