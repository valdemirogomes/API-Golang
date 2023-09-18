package errorhandling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApiError(t *testing.T) {
	type args struct {
		message string
		error   string
		status  int
		cause   CauseList
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewApiError",
			args: args{
				message: "message",
				error:   "error",
				status:  200,
				cause:   CauseList{},
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "error",
				ErrorStatus:  200,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIError(tt.args.message, tt.args.error, tt.args.status, tt.args.cause)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewNotFoundApiError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewNotFoundApiError",
			args: args{
				message: "message",
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "not_found",
				ErrorStatus:  404,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNotFoundAPIError(tt.args.message)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewTooManyRequestsError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewTooManyRequestsError",
			args: args{
				message: "message",
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "too_many_requests",
				ErrorStatus:  429,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTooManyRequestsError(tt.args.message)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewBadRequestApiError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewBadRequestApiError",
			args: args{
				message: "message",
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "bad_request",
				ErrorStatus:  400,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBadRequestAPIError(tt.args.message)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewMethodNotAllowedApiError(t *testing.T) {
	got := NewMethodNotAllowedAPIError()
	assert.Contains(t, got.Error(), "Method not allowed", "error message mismatch")
}

func TestNewInternalServerApiError(t *testing.T) {
	type args struct {
		message string
		cause   error
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewInternalServerApiError",
			args: args{
				message: "message",
				cause:   nil,
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "internal_server_error",
				ErrorStatus:  500,
				ErrorCause:   CauseList{},
			},
		},
		{
			name: "TestNewInternalServerApiError",
			args: args{
				message: "message",
				cause:   NewRequestError("cause"),
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "internal_server_error",
				ErrorStatus:  500,
				ErrorCause:   CauseList{NewRequestError("cause")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInternalServerAPIError(tt.args.message, tt.args.cause)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewForbiddenApiError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewForbiddenApiError",
			args: args{
				message: "message",
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "forbidden",
				ErrorStatus:  403,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewForbiddenAPIError(tt.args.message)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewUnauthorizedApiError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewUnauthorizedApiError",
			args: args{
				message: "message",
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "unauthorized_scopes",
				ErrorStatus:  401,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUnauthorizedAPIError(tt.args.message)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewConflictApiError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewConflictApiError",
			args: args{
				message: "message",
			},
			want: apiErr{
				ErrorMessage: "Can't update message due to a conflict error",
				ErrorCode:    "conflict_error",
				ErrorStatus:  409,
				ErrorCause:   CauseList{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConflictAPIError(tt.args.message)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestNewValidationApiError(t *testing.T) {
	type args struct {
		message string
		error   string
		cause   CauseList
	}
	tests := []struct {
		name string
		args args
		want ApiError
	}{
		{
			name: "TestNewValidationApiError",
			args: args{
				message: "message",
				error:   "error",
				cause:   CauseList{},
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "error",
				ErrorStatus:  400,
				ErrorCause:   CauseList{},
			},
		},
		{
			name: "TestNewValidationApiError",
			args: args{
				message: "message",
				error:   "error",
				cause:   CauseList{1, 2, 3},
			},
			want: apiErr{
				ErrorMessage: "message",
				ErrorCode:    "error",
				ErrorStatus:  400,
				ErrorCause:   CauseList{1, 2, 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewValidationAPIError(tt.args.message, tt.args.error, tt.args.cause)
			assert.Equal(t, tt.want.Error(), got.Error(), "error message mismatch")
		})
	}
}

func TestCauseList_ToString(t *testing.T) {
	tests := []struct {
		name string
		c    CauseList
		want string
	}{
		{
			name: "TestCauseList_ToString",
			c:    CauseList{},
			want: "[]",
		},
		{
			name: "TestCauseList_ToString",
			c:    CauseList{1, 2, 3},
			want: "[1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.ToString()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func Test_apiErr_Code(t *testing.T) {
	tests := []struct {
		name string
		e    apiErr
		want string
	}{
		{
			name: "Test_apiErr_Code",
			e:    apiErr{ErrorCode: "error"},
			want: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Code()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func Test_apiErr_Error(t *testing.T) {
	tests := []struct {
		name string
		e    apiErr
		want string
	}{
		{
			name: "Test_apiErr_Error",
			e:    apiErr{ErrorMessage: "message", ErrorCode: "error", ErrorStatus: 200, ErrorCause: CauseList{}},
			want: "Message: message;Error Code: error;Status: 200;Cause: []",
		},
		{
			name: "Test_apiErr_Error",
			e:    apiErr{ErrorMessage: "message", ErrorCode: "error", ErrorStatus: 200, ErrorCause: CauseList{1, 2, 3}},
			want: "Message: message;Error Code: error;Status: 200;Cause: [1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Error()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func Test_apiErr_Status(t *testing.T) {
	tests := []struct {
		name string
		e    apiErr
		want int
	}{
		{
			name: "Test_apiErr_Status",
			e:    apiErr{ErrorStatus: 200},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Status()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func Test_apiErr_Cause(t *testing.T) {
	tests := []struct {
		name string
		e    apiErr
		want CauseList
	}{
		{
			name: "Test_apiErr_Cause",
			e:    apiErr{ErrorCause: CauseList{1, 2, 3}},
			want: CauseList{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Cause()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}

func Test_apiErr_Message(t *testing.T) {
	tests := []struct {
		name string
		e    apiErr
		want string
	}{
		{
			name: "Test_apiErr_Message",
			e:    apiErr{ErrorMessage: "message"},
			want: "message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Message()
			assert.Equal(t, tt.want, got, "error message mismatch")
		})
	}
}
