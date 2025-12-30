package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal server error")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrValidation    = errors.New("validation error")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrap 包装错误
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

const (
	// 业务通用错误码
	CodeSuccess      = 0   // 操作成功
	CodeBadRequest   = 400 // 参数错误，非法请求
	CodeUnauthorized = 401 // 未授权，认证失败
	CodeForbidden    = 403 // 没有权限
	CodeNotFound     = 404 // 资源未找到
	CodeConflict     = 409 // 资源冲突，例如唯一性约束冲突
	CodeInternal     = 500 // 服务器内部错误
)

func ErrBadRequest(message string) *AppError {
	return NewAppError(CodeBadRequest, message, nil)
}

func ErrUnauthorizedError(message string) *AppError {
	return NewAppError(CodeUnauthorized, message, nil)
}

func ErrForbiddenError(message string) *AppError {
	return NewAppError(CodeForbidden, message, nil)
}

func ErrNotFoundError(message string) *AppError {
	return NewAppError(CodeNotFound, message, nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(CodeConflict, message, nil)
}

func ErrInternalError(message string) *AppError {
	return NewAppError(CodeInternal, message, nil)
}
