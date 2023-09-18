package errorhandling

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustom_Error(t *testing.T) {
	tests := []struct {
		name string
		re   *RequestError
		want string
	}{
		{
			name: "TestCustom_Error",
			re:   &RequestError{message: "message"},
			want: "message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.re.Error()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func TestCustom_Handle(t *testing.T) {
	tests := []struct {
		name string
		re   *RequestError
		want ApiError
	}{
		{
			name: "TestCustom_Handle",
			re:   &RequestError{message: "message"},
			want: NewBadRequestAPIError("message"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.re.Handle()
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewRequestError(t *testing.T) {
	tests := []struct {
		name string
		args string
		want error
	}{
		{
			name: "TestNewRequestError",
			args: "message",
			want: &RequestError{message: "message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRequestError(tt.args)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewRequestErrorf(t *testing.T) {
	tests := []struct {
		name string
		args string
		want error
	}{
		{
			name: "TestNewRequestErrorf",
			args: "message",
			want: &RequestError{message: "message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRequestErrorf(tt.args)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestSqlError_Error(t *testing.T) {
	tests := []struct {
		name string
		se   *sqlError
		want string
	}{
		{
			name: "TestSqlError_Error",
			se:   &sqlError{message: "message", method: "method", errType: "errType"},
			want: "Sql error in method [method]: message.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.se.Error()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func TestSqlError_Handle(t *testing.T) {
	tests := []struct {
		name string
		se   *sqlError
		want ApiError
	}{
		{
			name: "TestSqlError_Handle",
			se:   &sqlError{message: "message", method: "method", errType: "errType"},
			want: NewInternalServerAPIError("message", &sqlError{message: "message", method: "method", errType: "errType"}),
		},
		{
			name: "TestSqlError_Handle_notFoundSQLErrorType",
			se:   &sqlError{message: "message", method: "method", errType: "NotFound"},
			want: NewNotFoundAPIError("message"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.se.Handle()
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewSqlError(t *testing.T) {
	tests := []struct {
		name string
		args error
		want error
	}{
		{
			name: "TestNewSqlError",
			args: NewRequestError("cause"),
			want: &sqlError{message: "message. cause", method: "method", errType: "InternalError"},
		},
		{
			name: "TestNewSqlError_notFoundSQLErrorType",
			args: sql.ErrNoRows,
			want: &sqlError{message: "message. sql: no rows in result set", method: "method", errType: "NotFound"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSQLError(tt.args, "message", "method")
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}
