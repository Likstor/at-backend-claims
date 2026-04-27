package apperror

import (
	"context"
	"errors"
)

func GetErrorCtx(ctx context.Context, err error) context.Context {
	if c, ok := errors.AsType[interface {
		Ctx() context.Context
		Error() string
	}](err); ok {
		return c.Ctx()
	}

	return ctx
}
