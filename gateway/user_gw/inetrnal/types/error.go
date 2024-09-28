package types

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error represents an application-specific error with a message and code.
type Error struct {
	Message string `json:"message"`
	Code    uint16 `json:"code"`
	Success bool   `json:"success"`
}

// NewError creates a new custom error with the given code and message.
func NewError(code uint16, message string) *Error {
	return &Error{Code: code, Message: message, Success: false}
}

// NewInternalError creates a new internal server error with the given message.
func NewInternalError(message string) *Error {
	return &Error{Code: 13, Message: message, Success: false}
}

// NewNotFoundError creates a new not found error with the given message.
func NewNotFoundError(message string) *Error {
	return &Error{Code: 5, Message: message, Success: false}
}

// NewPermissionDeniedError creates a new permission denied error with the given message.
func NewPermissionDeniedError(message string) *Error {
	return &Error{Code: 7, Message: message, Success: false}
}

// NewBadRequestError creates a new bad request error with the given message.
func NewBadRequestError(message string) *Error {
	return &Error{Code: 10, Message: message, Success: false}
}

// ErrorToGRPCStatus converts the custom error to a gRPC status error based on the error code.
func (c *Error) ErrorToHttpStatus() int {

	switch c.Code {
	case 13:
		return fiber.StatusInternalServerError
	case 5:
		return fiber.StatusNotFound
	case 7:
		return fiber.StatusUnauthorized
	case 10:
		return fiber.StatusBadRequest
	default:
		return fiber.StatusBadRequest
	}
}

func (c *Error) ErrorToJsonMessage() map[string]interface{} {
	return map[string]interface{}{
		"message": c.Message,
		"success": false,
	}
}

func ExtractGRPCErrDetails(err error) *Error {
	details := &Error{}

	st, ok := status.FromError(err)
	if !ok {
		details.Code = uint16(codes.Unknown)
		details.Message = err.Error()
		return details
	}

	details.Code = uint16(st.Code())
	details.Message = st.Message()

	return details
}
