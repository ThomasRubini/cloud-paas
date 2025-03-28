package utils

import (
	"reflect"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/secretsprovider"
	"gorm.io/gorm"
)

type State struct {
	Db              *gorm.DB
	SecretsProvider secretsprovider.Helper
}

var state *State

func GetState() State {
	if state == nil {
		panic("state not set")
	}
	return *state
}

func SetState(s State) {
	if state != nil {
		panic("state already set")
	}
	state = &s
}

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
