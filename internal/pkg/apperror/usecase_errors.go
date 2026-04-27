package apperror

import "errors"

var (
	ErrUsecase = errors.New("usecase internal error")
)

// claims
var (
	ErrUnknownClaimCategory       = errors.New("unknown claim category")
	ErrIsNotImage                 = errors.New("unknown type, want \"image/png\" or \"image/jpeg\"")
	ErrPointIsNotInPolygon        = errors.New("coordinates of the claim outside the permissible district")
	ErrUserCannotPerformOperation = errors.New("user cannot perform this operation")
	ErrUnknownClaimStatus         = errors.New("unknown claim status")
)

// categories
var (
	ErrBadJSON           = errors.New("bad json")
	ErrUnprocessableJSON = errors.New("unprocessable json")
)