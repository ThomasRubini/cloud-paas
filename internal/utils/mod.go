package utils

import (
	"reflect"
)

func IsStatusCodeOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func CopyFields[A, B any](src *A, dst *B) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()
	copyMatchingFields(srcVal, dstVal)
}

func copyMatchingFields(srcVal, dstVal reflect.Value) {
	srcType := srcVal.Type()

	for i := 0; i < srcVal.NumField(); i++ {
		field := srcType.Field(i)
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(field.Name)

		// Handle embedded (anonymous) fields recursively
		if field.Anonymous {
			copyMatchingFields(srcField, dstVal)
			continue
		}

		// Copy matching fields
		if dstField.IsValid() && dstField.CanSet() && dstField.Type() == field.Type {
			dstField.Set(srcField)
		}
	}
}
