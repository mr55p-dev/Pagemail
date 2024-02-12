package tools

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
)

func LoadFromEnv(obj any) {
	// Loads strings from environment variables specified in `env` tag
	// will no-op if the argument is not a pointer to a struct
	rawVal := reflect.ValueOf(obj)
	if !isStructPtr(rawVal) {
		return
	}
	cVal := rawVal.Elem()
	cType := cVal.Type()
	for i := 0; i < cType.NumField(); i++ {
		valField := cVal.Field(i)
		if valField.Kind() != reflect.String {
			continue
		}
		typeField := cType.Field(i)
		tag := typeField.Tag.Get("env")
		val, ok := os.LookupEnv(tag)
		if ok {
			valField.SetString(val)
		}
	}
	return
}

func ValidateRequiredFields(obj any) error {
	rawVal := reflect.ValueOf(obj)
	if !isStructPtr(rawVal) {
		return fmt.Errorf("non struct type passed")
	}
	cVal := rawVal.Elem()
	cType := cVal.Type()
	for i := 0; i < cType.NumField(); i++ {
		valField := cVal.Field(i)
		typeField := cType.Field(i)
		tag := typeField.Tag.Get("required")
		if tag == "true" && valField.IsZero() {
			return fmt.Errorf("Zero value provided for required field %s", typeField.Name)
		}
	}
	return nil
}

func isStructPtr(v reflect.Value) bool {
	return v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct
}

func LogValue(v any) []slog.Attr {
	rawVal := reflect.ValueOf(v)
	var vals []slog.Attr
	if !isStructPtr(rawVal) {
		return vals
	}

	cVal := rawVal.Elem()
	cType := cVal.Type()
	for i := 0; i < cType.NumField(); i++ {
		valField := cVal.Field(i)
		if valField.Kind() != reflect.String {
			continue
		}
		typeField := cType.Field(i)
		env := typeField.Tag.Get("log")
		if env != "" {
			vals = append(vals, slog.String(env, valField.String()))
		}
	}
	return vals
}
