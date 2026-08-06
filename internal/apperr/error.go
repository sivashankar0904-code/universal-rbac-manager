package apperr

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

type Kind string

const (
	KindInternal        Kind = "internal"
	KindValidation      Kind = "validation"
	KindInputBody       Kind = "input_body"
	KindConflict        Kind = "conflict"
	KindNotFound        Kind = "not_found"
	KindUnauthenticated Kind = "unauthenticated"
	KindUnauthorized    Kind = "unauthorized"
	KindMaxAttempts     Kind = "max_attempts"
)

type Error struct {
	kind       Kind
	message    string
	original   error
	fileLine   string
	fields     map[string]string
	retryAfter time.Duration
}

func (e *Error) Error() string {
	if e.original != nil {
		return e.fileLine + ": " + e.original.Error()
	}
	return e.fileLine + ": " + e.message
}

func (e *Error) Unwrap() error {
	return e.original
}

func (e *Error) Kind() Kind {
	return e.kind
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Fields() map[string]string {
	return e.fields
}

func (e *Error) RetryAfter() time.Duration {
	return e.retryAfter
}

func (e *Error) HTTPStatusCode() int {
	switch e.kind {
	case KindValidation:
		return http.StatusUnprocessableEntity
	case KindInputBody:
		return http.StatusBadRequest
	case KindConflict:
		return http.StatusConflict
	case KindNotFound:
		return http.StatusNotFound
	case KindUnauthenticated:
		return http.StatusUnauthorized
	case KindUnauthorized:
		return http.StatusForbidden
	case KindMaxAttempts:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func (e *Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("kind", string(e.kind)),
		slog.String("msg", e.message),
		slog.String("at", e.fileLine),
	)
}

func callerFileLine() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func NotFound(msg string) *Error {
	return &Error{kind: KindNotFound, message: msg, fileLine: callerFileLine()}
}

func Conflict(msg string) *Error {
	return &Error{kind: KindConflict, message: msg, fileLine: callerFileLine()}
}

func InputBody(msg string) *Error {
	return &Error{kind: KindInputBody, message: msg, fileLine: callerFileLine()}
}

func Unauthenticated(msg string) *Error {
	return &Error{kind: KindUnauthenticated, message: msg, fileLine: callerFileLine()}
}

func Unauthorized(msg string) *Error {
	return &Error{kind: KindUnauthorized, message: msg, fileLine: callerFileLine()}
}

func Internal(msg string) *Error {
	return &Error{kind: KindInternal, message: msg, fileLine: callerFileLine()}
}

func Validation(msg string, fields map[string]string) *Error {
	return &Error{kind: KindValidation, message: msg, fields: fields, fileLine: callerFileLine()}
}

func MaxAttempts(msg string, retryAfter time.Duration) *Error {
	return &Error{kind: KindMaxAttempts, message: msg, retryAfter: retryAfter, fileLine: callerFileLine()}
}

func Wrap(err error, kind Kind, msg string) *Error {
	return &Error{kind: kind, message: msg, original: err, fileLine: callerFileLine()}
}
