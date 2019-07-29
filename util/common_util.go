package util

import (
	"reflect"
)

// Contains judge if `collection` contains `item`.
// Note that `collection` and `item` should be the the *same* type.
func Contains(item interface{}, collection interface{}) bool {
	targetValue := reflect.ValueOf(collection)
	switch reflect.TypeOf(collection).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == item {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(item)).IsValid() {
			return true
		}
	}
	return false
}
