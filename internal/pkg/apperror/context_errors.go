package apperror

import "errors"

var (
	ErrCtxEmptyUserID = errors.New("empty userID in context")
	ErrCtxEmptyRole   = errors.New("empty role in context")
)
