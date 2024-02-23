package configLoader

import "fmt"

type KeyNotPresent struct {
	msg string
}

func (e *KeyNotPresent) Error() string {
	return e.msg
}

type InvalidValue struct {
	msg string
}

func (e *InvalidValue) Error() string {
	return e.msg
}

func errKeyNotPresent(key string) *KeyNotPresent {
	return &KeyNotPresent{
		msg: fmt.Sprintf("Key %s not found\n", key),
	}
}

func errInvalidValue(key string) *InvalidValue {
	return &InvalidValue{
		msg: fmt.Sprintf("Attempted to coerce invalid value for %s", key),
	}
}
