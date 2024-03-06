package httpit

import (
	"fmt"
	"net/http"
	"reflect"
)

var (
	formKey  = "form"
	queryKey = "query"
)

type bindError struct{ msg string }

func (b *bindError) Error() string {
	return fmt.Sprintf("Error binding data: %s", b.msg)
}

func newBindError(format string, parts ...any) *bindError {
	return &bindError{fmt.Sprintf(format, parts...)}
}

func bind(in any, r *http.Request) error {
	refTypeOf := reflect.TypeOf(in)
	refKindOf := refTypeOf.Kind()
	if refKindOf != reflect.Pointer && refTypeOf.Elem().Kind() != reflect.Struct {
		return newBindError("Cannot bind to non-struct pointers or value types")
	}

	valueOf := reflect.ValueOf(in).Elem()
	typeOf := valueOf.Type()
	err := r.ParseForm()
	if err != nil {
		return newBindError("Failed to parse the request form body: %s", err.Error())
	}

	for i := 0; i < typeOf.NumField(); i++ {
		fieldVal := valueOf.Field(i)
		fieldType := typeOf.Field(i)

		formVal := fieldType.Tag.Get(formKey)
		queryVal := fieldType.Tag.Get(queryKey)
		if formVal == "" && queryVal == "" {
			continue
		}
		if !fieldVal.CanSet() {
			return newBindError("Attempted to bind parameter on un-settable field %s", fieldType.Name)
		}

		switch fieldVal.Kind() {
		case reflect.String:
			if formVal != "" {
				v := r.Form.Get(formVal)
				fieldVal.SetString(v)
			} else if queryVal != "" {
				v := r.URL.Query().Get(queryVal)
				fieldVal.SetString(v)
			}
		default:
			return newBindError("Cannot bind field of type %s (field %s in %s)", fieldType.Type, fieldType.Name, typeOf.Name())
		}
	}

	return nil
}
