package apperror

import "errors"

var (
	ErrRepository = errors.New("repository internal error")
	
	ErrClaimNotExists = errors.New("claim does not exists")

	ErrCategoryNotExists = errors.New("category or subcategory does not exists")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrSubcategoryAlreadyExists = errors.New("subcategory already exists")
)