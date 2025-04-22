package err_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "github.com/eclipse-xfsc/microservice-core-go/pkg/err"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	e := errors.New("something went wrong")
	assert.Implements(t, (*error)(nil), e)

	// create error with from a Kind
	e = errors.New(errors.Forbidden)
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Forbidden, e))
	assert.Contains(t, e.Error(), "permission denied")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Empty(t, ec.Message)
		assert.Equal(t, errors.Forbidden, ec.Kind)
		assert.Nil(t, ec.Err)
	}

	// create error with a string message only
	e = errors.New("something went wrong")
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Unknown, e))
	assert.Contains(t, e.Error(), "something went wrong")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, "something went wrong", ec.Message)
		assert.Equal(t, errors.Unknown, ec.Kind)
		assert.Nil(t, ec.Err)
	}

	// create error with Kind and Message
	e = errors.New(errors.Internal, "something went wrong")
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Internal, e))
	assert.Contains(t, e.Error(), "something went wrong: internal error")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, "something went wrong", ec.Message)
		assert.Equal(t, errors.Internal, ec.Kind)
		assert.Nil(t, ec.Err)
	}

	// create error from a previous error
	e = errors.New(fmt.Errorf("oops it did it again"))
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Unknown, e))
	assert.Contains(t, e.Error(), "oops it did it again")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, errors.Unknown, ec.Kind)
		assert.NotNil(t, ec.Err)
		assert.Equal(t, ec.Err.Error(), "oops it did it again")
	}

	// create error from a previous structured error
	e = errors.New(errors.New(errors.Unauthorized, "no way out"))
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Unauthorized, e))
	assert.Contains(t, e.Error(), "no way out: not authenticated")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, errors.Unauthorized, ec.Kind)
		assert.Equal(t, ec.Message, "no way out")
		assert.NotNil(t, ec.Err)
		assert.Contains(t, ec.Err.Error(), "no way out: not authenticated")
	}

	// create error from a previous structured error
	e = errors.New(errors.BadRequest, "bad request", errors.New(errors.Unauthorized, "no way out"))
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.BadRequest, e))
	assert.Contains(t, e.Error(), "bad request: no way out: not authenticated")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, errors.BadRequest, ec.Kind)
		assert.Equal(t, "bad request", ec.Message)
		assert.NotNil(t, ec.Err)
		assert.Contains(t, ec.Err.Error(), "no way out: not authenticated")
	}

	// create error from a previous structured error and a message
	e = errors.New("bad request", errors.New(errors.Unauthorized, "no way out"))
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Unauthorized, e))
	assert.Contains(t, e.Error(), "bad request: not authenticated: no way out: not authenticated")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, errors.Unauthorized, ec.Kind)
		assert.Equal(t, "bad request", ec.Message)
		assert.NotNil(t, ec.Err)
		assert.Contains(t, ec.Err.Error(), "no way out: not authenticated")
	}

	// create error from a previous structured error and a message
	e = errors.New(errors.BadRequest, errors.New(errors.Unauthorized, "no way out"))
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.BadRequest, e))
	assert.Contains(t, e.Error(), "bad request: no way out: not authenticated")
	if ec, ok := e.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, errors.BadRequest, ec.Kind)
		assert.Equal(t, "no way out", ec.Message)
		assert.NotNil(t, ec.Err)
		assert.Contains(t, ec.Err.Error(), "no way out: not authenticated")
	}

	// create three nested errors
	e1 := fmt.Errorf("cannot insert record")
	e2 := errors.New("account already exists", e1)
	e3 := errors.New("failed to create account", e2)
	assert.Contains(t, e3.Error(), "failed to create account: account already exists: cannot insert record")
	if ec, ok := e3.(*errors.Error); ok {
		assert.NotEmpty(t, ec.ID)
		assert.Equal(t, errors.Unknown, ec.Kind)
		assert.Equal(t, "failed to create account", ec.Message)
		assert.NotNil(t, ec.Err)
		assert.Contains(t, ec.Err.Error(), "account already exists")
	}
}

func TestIs(t *testing.T) {
	e := errors.New(errors.Timeout)
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Timeout, e))

	// create error from a previous structured error and a message
	e = errors.New(errors.Timeout, errors.New(errors.Unauthorized))
	assert.IsType(t, &errors.Error{}, e)
	assert.True(t, errors.Is(errors.Timeout, e))
}

func TestError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		kind errors.Kind
		code int
	}{
		{
			name: "unknown error",
			kind: errors.Unknown,
			code: http.StatusInternalServerError,
		},
		{
			name: "bad request",
			kind: errors.BadRequest,
			code: http.StatusBadRequest,
		},
		{
			name: "unauthorized",
			kind: errors.Unauthorized,
			code: http.StatusUnauthorized,
		},
		{
			name: "forbidden",
			kind: errors.Forbidden,
			code: http.StatusForbidden,
		},
		{
			name: "exists",
			kind: errors.Exist,
			code: http.StatusConflict,
		},
		{
			name: "not found",
			kind: errors.NotFound,
			code: http.StatusNotFound,
		},
		{
			name: "timeout",
			kind: errors.Timeout,
			code: http.StatusRequestTimeout,
		},
		{
			name: "internal",
			kind: errors.Internal,
			code: http.StatusInternalServerError,
		},
		{
			name: "service unavailable",
			kind: errors.ServiceUnavailable,
			code: http.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := errors.New(test.kind)
			assert.IsType(t, &errors.Error{}, e)
			if ec, ok := e.(*errors.Error); ok {
				assert.Equal(t, test.code, ec.StatusCode())
			}
		})
	}
}

func TestError_MarshalJSON(t *testing.T) {
	e := errors.New(errors.NotFound, "item does not exist")
	ec, ok := e.(*errors.Error)
	assert.True(t, ok)
	id := ec.ID

	data, err := ec.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	ee := &errors.Error{}
	err = ee.UnmarshalJSON(data)
	assert.NoError(t, err)

	assert.Equal(t, id, ee.ID)
	assert.Equal(t, ec.Kind, ee.Kind)
	assert.Equal(t, ec.Message, ee.Message)
	assert.Equal(t, ec.Error(), ee.Error())
	assert.Equal(t, errors.NotFound, ee.Kind)
}

func TestError_Temporary(t *testing.T) {
	e := errors.New(errors.Forbidden)
	ec, ok := e.(*errors.Error)
	assert.True(t, ok)
	assert.False(t, ec.Temporary())

	e = errors.New(errors.Internal)
	ec, ok = e.(*errors.Error)
	assert.True(t, ok)
	assert.True(t, ec.Temporary())
}

func TestGetKind(t *testing.T) {
	tests := []struct {
		name string
		code int
		kind errors.Kind
	}{
		{
			name: "undefined HTTP status code",
			code: 9999,
			kind: errors.Unknown,
		},
		{
			name: "bad request",
			code: http.StatusBadRequest,
			kind: errors.BadRequest,
		},
		{
			name: "unauthorized",
			code: http.StatusUnauthorized,
			kind: errors.Unauthorized,
		},
		{
			name: "forbidden",
			code: http.StatusForbidden,
			kind: errors.Forbidden,
		},
		{
			name: "exists",
			code: http.StatusConflict,
			kind: errors.Exist,
		},
		{
			name: "not found",
			code: http.StatusNotFound,
			kind: errors.NotFound,
		},
		{
			name: "timeout",
			code: http.StatusRequestTimeout,
			kind: errors.Timeout,
		},
		{
			name: "internal",
			code: http.StatusInternalServerError,
			kind: errors.Internal,
		},
		{
			name: "service unavailable",
			code: http.StatusServiceUnavailable,
			kind: errors.ServiceUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kind := errors.GetKind(test.code)
			assert.Equal(t, test.kind, kind)
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		// input
		name       string
		err        error
		statusCode int

		// output
		responseCode  int
		responseError *errors.Error
	}{
		{
			name:         "error is nil",
			err:          nil,
			responseCode: http.StatusInternalServerError,
			responseError: &errors.Error{
				Kind: errors.Unknown,
			},
		},
		{
			name:         "simple text error",
			err:          fmt.Errorf("simple text error"),
			responseCode: http.StatusInternalServerError,
			responseError: &errors.Error{
				Kind:    errors.Unknown,
				Message: "",
			},
		},
		{
			name:         "structured error",
			err:          errors.New("structured error"),
			responseCode: http.StatusInternalServerError,
			responseError: &errors.Error{
				Kind:    errors.Unknown,
				Message: "structured error",
			},
		},
		{
			name:         "structured error with kind",
			err:          errors.New(errors.Forbidden, "structured error with kind"),
			responseCode: http.StatusForbidden,
			responseError: &errors.Error{
				Kind:    errors.Forbidden,
				Message: "structured error with kind",
			},
		},
		{
			name:         "structured error with kind and embedded error",
			err:          errors.New(errors.NotFound, "structured error with kind and embedded error", fmt.Errorf("embedded error")),
			responseCode: http.StatusNotFound,
			responseError: &errors.Error{
				Kind:    errors.NotFound,
				Message: "structured error with kind and embedded error",
				Err:     fmt.Errorf("embedded error"),
			},
		},
		{
			name:         "structured error with kind only",
			err:          errors.New(errors.Timeout),
			responseCode: http.StatusRequestTimeout,
			responseError: &errors.Error{
				Kind:    errors.Timeout,
				Message: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			if test.statusCode > 0 {
				errors.JSON(rr, test.err, test.statusCode)
			} else {
				errors.JSON(rr, test.err)
			}

			assert.Equal(t, test.responseCode, rr.Code)

			var responseError *errors.Error
			err := json.NewDecoder(rr.Body).Decode(&responseError)
			assert.NoError(t, err)
			assert.Equal(t, test.responseError.Kind, responseError.Kind)
			assert.Equal(t, test.responseError.Message, responseError.Message)
		})
	}
}
