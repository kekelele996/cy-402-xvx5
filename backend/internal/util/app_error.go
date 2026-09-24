package util

import (
	"errors"
	"fmt"

	"cylawcase/internal/constants"
	"cylawcase/internal/repository"
)

// AppError 业务错误，携带统一错误码。
type AppError struct {
	Code    int
	Message string
}

// Error 实现 error 接口。
func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

// NewAppError 构造业务错误。
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap 包装错误并附带上下文；仓储层哨兵错误在此统一映射为带错误码的 AppError，
// 保证 service/handler 层层透传后仍可被 errors.As 识别（如 ErrNotFound -> 404）。
func Wrap(err error, format string, args ...any) error {
	context := fmt.Sprintf(format, args...)
	wrapped := fmt.Errorf("%s: %w", context, err)
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	if errors.Is(err, repository.ErrNotFound) {
		return NewAppError(constants.CodeNotFound, context+": resource not found")
	}
	if errors.Is(err, repository.ErrDuplicate) {
		return NewAppError(constants.CodeConflict, context+": duplicate record")
	}
	return wrapped
}
