package apperror

import "errors"

// file manager
var (
	ErrFileManager   = errors.New("file manager internal error")
	ErrFileNotExists = errors.New("file or directory does not exists")
	ErrFileExists    = errors.New("file exists")
)