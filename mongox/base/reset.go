package base

import (
	"fmt"
	"reflect"
)

// Reset function creates new zero object for the target pointer
func Reset(target interface{}) (created bool) {
	type resetter interface {
		Reset()
	}

	resettable, canReset := target.(resetter)
	if canReset {
		resettable.Reset()
		return false
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Errorf("reset target should be a pointer"))
	}

	t := v.Elem().Type()
	zero := reflect.Zero(t)

	v.Elem().Set(zero)

	return true
}
