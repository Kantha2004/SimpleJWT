package services

// Custom error types for service layer
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{Message: message}
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	return e.Message
}

func NewInternalError(message string, err error) *InternalError {
	return &InternalError{Message: message, Err: err}
}
