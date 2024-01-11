package tools

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

type EnvTag struct {
	Value    string
	Required bool
}

func ParseEnvTag(tag string) *EnvTag {
	vals := strings.Split(tag, ",")
	out := &EnvTag{
		Value: vals[0],
	}
	if len(vals) > 1 {
		for _, v := range vals[1:] {
			switch v {
			case "required":
				out.Required = true
			}
		}
	}
	return out
}

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
		envTag := ParseEnvTag(tag)
		val, ok := os.LookupEnv(envTag.Value)
		if envTag.Required && !ok {
			panic(fmt.Sprintf("Could not read env %s", tag))
		}
		if ok {
			valField.SetString(val)
		}
	}
	return
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
