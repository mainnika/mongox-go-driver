package base

import (
	"reflect"
)

// GetProtection function finds protection field in the source document otherwise returns nil
func GetProtection(source interface{}) *Protection {

	v := reflect.ValueOf(source)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}

	el := v.Elem()
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)

		switch field.Interface().(type) {
		case *Protection:
			return field.Interface().(*Protection)
		case Protection:
			ptr := field.Addr()
			return ptr.Interface().(*Protection)
		default:
			continue
		}
	}

	return nil
}
