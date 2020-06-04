package utils

import (
	"unsafe"
)

// IsNil function evaluates the interface value to nil
func IsNil(i interface{}) bool {

	type iface struct {
		_   unsafe.Pointer
		ptr unsafe.Pointer
	}

	unpacked := (*iface)(unsafe.Pointer(&i))
	if unpacked.ptr == nil {
		return true
	}

	return *(*unsafe.Pointer)(unpacked.ptr) == nil
}
