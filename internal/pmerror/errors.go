package pmerror

import (
	"net/http"
)

type PMError struct {
	Message string
	Status  int
}

func (c *PMError) Error() string {
	return c.Message
}

var (
	ErrNoParam = &PMError{
		Status:  http.StatusBadRequest,
		Message: "Please include all marked fields",
	}
	ErrBadEmail = &PMError{
		Status:  http.StatusNotFound,
		Message: "Incorrect email address",
	}
	ErrBadPassword = &PMError{
		Status:  http.StatusUnauthorized,
		Message: "Incorrect password",
	}
	ErrDiffPasswords = &PMError{
		Status:  http.StatusBadRequest,
		Message: "Passwords do not match",
	}
	ErrDuplicateEmail = &PMError{
		Status:  http.StatusBadRequest,
		Message: "Looks like that email address is already taken. If you can't remember your password please reach out to help@pagemail.io for assistence.",
	}
	ErrNoAuth = &PMError{
		Status:  http.StatusForbidden,
		Message: "Looks like your account was created with a different provider.",
	}
	ErrUnspecified = &PMError{
		Status:  http.StatusInternalServerError,
		Message: "Something went wrong",
	}

	ErrBadPagination = &PMError{
		Status:  http.StatusBadRequest,
		Message: "Invalid page number",
	}
	ErrNoPage = &PMError{
		Status:  http.StatusNotFound,
		Message: "No page found",
	}

	ErrNotAllowed = &PMError{
		Status:  http.StatusForbidden,
		Message: "Permission denied",
	}

	ErrCreatingMail = &PMError{
		Status:  http.StatusInternalServerError,
		Message: "Failed to generate reset email",
	}
)

func NewInternalError(msg string) *PMError {
	return &PMError{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}
