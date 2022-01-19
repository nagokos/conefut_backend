package utils

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	// 認証エラー
	AppErrorCodeAuthenticationFailure AppErrorCode = "AUTHENTICATION_FAILURE"
	// 認可エラー
	AppErrorCodeAuthorizationFailure AppErrorCode = "AUTHORIZATION_FAILURE"
	// バリデーションエラー
	AppErrorCodeValidationFailure AppErrorCode = "VALIDATION_FAILURE"
	// その他の予期せぬエラー
	AppErrorCodeUnexpectedFailure AppErrorCode = "UNEXPECTED_FAILURE"
)

type AppError struct {
	httpStatusCode int
	appErorrCode   AppErrorCode
	attribute      string
	message        string
}

type AppErrorCode string

type AppErrorOption func(*AppError)

func (e *AppError) AddGraphQLError(ctx context.Context) {
	extensions := map[string]interface{}{
		"status_code": e.httpStatusCode,
		"error_code":  e.appErorrCode,
	}
	if e.attribute != "" {
		extensions["attribute"] = e.attribute
	}
	graphql.AddError(ctx, &gqlerror.Error{
		Message:    e.message,
		Extensions: extensions,
	})
}

func NewAppError(httpStatusCode int, appErrorCode AppErrorCode, message string, opts ...AppErrorOption) *AppError {
	e := &AppError{
		httpStatusCode: httpStatusCode,
		appErorrCode:   appErrorCode,
		message:        message,
	}

	for _, o := range opts {
		o(e)
	}

	return e
}

func NewValidationError(field, message string) *AppError {
	options := []AppErrorOption{WithField(field)}
	return NewAppError(http.StatusBadRequest, AppErrorCodeValidationFailure, message, options...)
}

func NewAuthenticationErorr(message string, opts ...AppErrorOption) *AppError {
	return NewAppError(http.StatusUnauthorized, AppErrorCodeAuthenticationFailure, message, opts...)
}

func WithField(v string) AppErrorOption {
	return func(a *AppError) {
		a.attribute = v
	}
}
