package utils

import (
	"errors"
	"fmt"
)

// Error codes
const (
	// Auth related error codes
	ErrUsernameExists    = "username_exists"
	ErrEmailExists       = "email_exists"
	ErrInvalidCredential = "invalid_credential"
	ErrTokenExpired      = "token_expired"
	ErrUserNotFound      = "user_not_found"
	// Add more error codes as needed
)

// AppError represents an application error with a code and message
type AppError struct {
	Code    string
	Message string
}

// Error implements the error interface for AppError
func (e *AppError) Error() string {
	return e.Message
}

// Is implements error comparison for AppError
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ""
}

// NewAppError creates a new AppError with code and message
func NewAppError(code string, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorf creates a new AppError with code and formatted message
func NewAppErrorf(code string, format string, args ...any) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// TranslatedError creates an error with translated message based on message ID
func TranslatedError(messageID string, templateData ...map[string]any) error {
	return errors.New(T(messageID, templateData...))
}

// NewTranslatedAppError creates an AppError with a code and translated message
func NewTranslatedAppError(code string, messageID string, templateData ...map[string]any) *AppError {
	return &AppError{
		Code:    code,
		Message: T(messageID, templateData...),
	}
}
