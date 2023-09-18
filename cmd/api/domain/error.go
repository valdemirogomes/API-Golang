package domain

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mercadolibre/fury_go-core/pkg/web"
)

type AppError struct {
	Message string
	Err     error
}

func NewAppError(message string, originalError error) *AppError {
	if originalError == nil {
		originalError = errors.New(message)
	}

	return &AppError{
		Message: message,
		Err:     originalError,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

type NotFoundError struct {
	*AppError // AppError
}

type InternalError struct {
	*AppError // AppError
}

type BadRequest struct {
	*AppError // AppError
}

func NewNotFoundError(message string, originalError error) *NotFoundError {
	return &NotFoundError{
		AppError: NewAppError(message, originalError),
	}
}

func NewInternalError(message string, originalError error) *InternalError {
	return &InternalError{
		AppError: NewAppError(message, originalError),
	}
}

func NewBadRequest(message string, originalError error) *BadRequest {
	return &BadRequest{
		AppError: NewAppError(message, originalError),
	}
}

func ConvertToWebErr(err error) error {
	switch e := err.(type) {
	case *AppError:
		return web.NewError(http.StatusInternalServerError, e.Error())
	case *NotFoundError:
		return web.NewError(http.StatusNotFound, e.Error())
	case *InternalError:
		return web.NewError(http.StatusInternalServerError, e.Error())
	case *BadRequest:
		return web.NewError(http.StatusBadRequest, e.Error())
	default:
		return web.NewError(http.StatusInternalServerError, err.Error())

	}
}

// Funções para obter a mensagem de erro
func (e *AppError) getMessage() string {
	if e.Err != nil {
		return fmt.Sprintf("%s original_cause: %s", e.Message, e.Err.Error())
	}
	return e.Message
}
