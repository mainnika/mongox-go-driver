package base

import (
	"reflect"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
)

// GetProtection function finds protection field in the source document otherwise returns nil
func GetProtection(source interface{}) *protection.Key {

	v := reflect.ValueOf(source)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}

	el := v.Elem()
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)

		switch field.Interface().(type) {
		case *protection.Key:
			return field.Interface().(*protection.Key)
		case protection.Key:
			ptr := field.Addr()
			return ptr.Interface().(*protection.Key)
		default:
			continue
		}
	}

	return nil
}
