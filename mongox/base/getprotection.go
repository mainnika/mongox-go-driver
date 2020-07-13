package base

import (
	"reflect"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
)

// GetProtection function finds protection field in the source document otherwise returns nil
func GetProtection(source interface{}) (key *protection.Key) {

	v := reflect.ValueOf(source)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}

	el := v.Elem()
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		if !field.CanInterface() {
			continue
		}

		switch field.Interface().(type) {
		case *protection.Key:
			key = field.Interface().(*protection.Key)
		case protection.Key:
			ptr := field.Addr()
			key = ptr.Interface().(*protection.Key)
		default:
			continue
		}

		return
	}

	return
}
