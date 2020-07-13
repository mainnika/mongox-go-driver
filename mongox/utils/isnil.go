package utils

import (
	"unsafe"
)

// IsNil function evaluates the interface value to nil
func IsNil(i interface{}) (isNil bool) {

	type iface struct {
		_   unsafe.Pointer
		ptr unsafe.Pointer
	}

	unpacked := (*iface)(unsafe.Pointer(&i))
	if unpacked.ptr == nil {
		isNil = true
		return
	}

	isNil = *(*unsafe.Pointer)(unpacked.ptr) == nil
	return
}
