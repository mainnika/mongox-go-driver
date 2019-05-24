package base

import (
	"reflect"

	"github.com/mainnika/mongox-go-driver/mongox"
)

// Reset function creates new zero object for the target pointer
func Reset(target interface{}) {

	resettable, canReset := target.(mongox.Resetter)
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
