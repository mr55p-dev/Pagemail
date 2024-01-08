package db

import (
	"fmt"
	"reflect"
)

type StructData struct {
	Fields []string
	Values []any
	Refs   []any
}

func getStructData(record any) (*StructData, error) {
	recPtrVal := reflect.ValueOf(record)
	recPtrType := recPtrVal.Type()
	if recPtrType.Kind() != reflect.Pointer || recPtrType.Elem().Kind() != reflect.Struct {
		return nil, &DBError{ErrInvalidArg, "cannot create record with non-struct pointer arg"}
	}

	recValue := recPtrVal.Elem()
	recType := recValue.Elem().Type()

	var fields []string
	var values []any
	var refs []any

	for i := 0; i < recType.NumField(); i++ {
		fieldType := recType.Field(i)
		fieldValue := recValue.Field(i)
		colname := fieldType.Tag.Get("db")
		if colname == "" {
			continue
		}

		ptr := fieldValue.Addr().UnsafePointer()
		ref := reflect.NewAt(fieldValue.Type(), ptr).Interface()

		fields = append(fields, colname)
		values = append(values, recValue.Field(i).Interface())
		refs = append(refs, ref)
	}

	if len(fields) == 0 || len(values) == 0 {
		return nil, &DBError{ErrInvalidArg, "no tagged fields found on input"}
	}

	return &StructData{
		Fields: fields,
		Values: values,
		Refs:   refs,
	}, nil

}

func getStructField(obj any, fieldName string, value any) any {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr || objValue.Elem().Kind() != reflect.Struct {
		return &DBError{ErrInvalidArg, fmt.Sprintf("obj must be a pointer to a struct")}
	}

	field := objValue.Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return &DBError{ErrInvalidArg, fmt.Sprintf("field '%s' not found in the struct", fieldName)}
	}

	if !field.CanSet() {
		return &DBError{ErrInvalidArg, fmt.Sprintf("field '%s' cannot be set", fieldName)}
	}

	fieldType := field.Type()
	val := reflect.ValueOf(value)

	if !val.Type().AssignableTo(fieldType) {
		return &DBError{ErrInvalidArg, fmt.Sprintf("value type %v cannot be assigned to field type %v", val.Type(), fieldType)}
	}

	field.Set(val)

	return nil
}
func setStructField(obj any, fieldName string, value any) error {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr || objValue.Elem().Kind() != reflect.Struct {
		return &DBError{ErrInvalidArg, fmt.Sprintf("obj must be a pointer to a struct")}
	}

	field := objValue.Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return &DBError{ErrInvalidArg, fmt.Sprintf("field '%s' not found in the struct", fieldName)}
	}

	if !field.CanSet() {
		return &DBError{ErrInvalidArg, fmt.Sprintf("field '%s' cannot be set", fieldName)}
	}

	fieldType := field.Type()
	val := reflect.ValueOf(value)

	if !val.Type().AssignableTo(fieldType) {
		return &DBError{ErrInvalidArg, fmt.Sprintf("value type %v cannot be assigned to field type %v", val.Type(), fieldType)}
	}

	field.Set(val)

	return nil
}
