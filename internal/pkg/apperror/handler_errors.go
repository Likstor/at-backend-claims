package apperror

import "errors"

var (
	ErrHandler = errors.New("handler internal error")
)