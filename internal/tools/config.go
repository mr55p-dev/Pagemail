package tools

import (
	"log/slog"
	"reflect"
)

func LogValue(v any) []slog.Attr {
	// TODO: use code generation for this instead
	rawVal := reflect.ValueOf(v)
	var vals []slog.Attr

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
