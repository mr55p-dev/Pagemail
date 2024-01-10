package db

import "fmt"

type DBErrorType string

const (
	ErrInvalidArg DBErrorType = "Invalid argument"
)

type DBError struct {
	Type    DBErrorType
	Message string
}

func (e *DBError) Error() string {
	return fmt.Sprintf("DB error (%s): %s", string(e.Type), e.Message)
}
