package errorhandling

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mercadolibre/fury_go-core/pkg/log"
)

type customError interface {
	Handle() (ApiError, log.Level)
}

type RequestError struct {
	message string
}

func (re *RequestError) Error() string {
	return re.message
}

func (re *RequestError) Handle() (ApiError, log.Level) {
	return NewBadRequestAPIError(re.Error()), log.WarnLevel
}

func NewRequestError(message string) error {
	return &RequestError{message: message}
}

func NewRequestErrorf(message string, args ...interface{}) error {
	return &RequestError{message: fmt.Sprintf(message, args...)}
}

const (
	notFoundSQLErrorType = "NotFound"
	internalErrorType    = "InternalError"
)

type sqlError struct {
	message string
	method  string
	errType string
}

func (se *sqlError) Error() string {
	return fmt.Sprintf("Sql error in method [%s]: %s.", se.method, se.message)
}

func (se *sqlError) Handle() (ApiError, log.Level) {
	if se.errType == notFoundSQLErrorType {
		return NewNotFoundAPIError(se.message), log.WarnLevel
	}
	return NewInternalServerAPIError(se.message, se), log.ErrorLevel
}

func NewSQLError(err error, message string, method string) error {
	var errType string
	if errors.Is(err, sql.ErrNoRows) {
		errType = notFoundSQLErrorType
	} else {
		errType = internalErrorType
	}

	return &sqlError{message: fmt.Sprintf("%s. %s", message, err.Error()), method: method, errType: errType}
}
