package utils

import "reflect"

func IsStatusCodeOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func CopyFields[A, B any](src *A, dst *B) {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)
		if dstField.IsValid() {
			dstField.Set(srcValue.Field(i))
		}
	}
}
