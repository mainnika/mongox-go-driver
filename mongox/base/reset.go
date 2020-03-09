package base

import (
	"reflect"
)

// Reset function creates new zero object for the target pointer
func Reset(target interface{}) {

	type resetter interface {
		Reset()
	}

	resettable, canReset := target.(resetter)
	if canReset {
		resettable.Reset()
		return
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		panic("reset target should be a pointer")
	}

	t := v.Elem().Type()
	zero := reflect.Zero(t)

	v.Elem().Set(zero)
}
