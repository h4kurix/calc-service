package errors

import "errors"

// Calculation Errors
var (
	ErrInvalidExpression     = NewError("недопустимое выражение")
	ErrUnbalancedParentheses = NewError("несбалансированные скобки")
	ErrDivisionByZero        = NewError("деление на ноль")
	ErrUnknownOperator       = NewError("неизвестный оператор")
)

// HTTP Handler Errors
var (
	ErrInvalidSymbol       = NewError("expression is not valid")
	ErrInternalServerError = NewError("internal server error")
	ErrMethodNotAllowed    = NewError("method not allowed")
	ErrInvalidRequestBody  = NewError("invalid request body")
)

// NewError creates a new error with the given message
func NewError(message string) error {
	return errors.New(message)
}

// HTTPError defines an HTTP error with status code
type HTTPError struct {
	Message    string
	StatusCode int
}

func (e HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(message string, statusCode int) HTTPError {
	return HTTPError{Message: message, StatusCode: statusCode}
}

// HTTP Status Codes
const (
	StatusMethodNotAllowed    = 405
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
)
