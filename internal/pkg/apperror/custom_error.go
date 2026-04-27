package apperror

import "context"

type appErrorWithoutMsg struct {
	next error
	ctx  context.Context
}

func NewErrorCtxWithoutMsg(ctx context.Context, err error) error {
	return &appErrorWithoutMsg{
		ctx:  ctx,
		next: err,
	}
}

func (e appErrorWithoutMsg) Error() string {
	return e.next.Error()
}

func (e appErrorWithoutMsg) Unwrap() error {
	return e.next
}

func (e appErrorWithoutMsg) Ctx() context.Context {
	return e.ctx
}

type appError struct {
	appErrorWithoutMsg
	msg string
}

func NewErrorCtx(ctx context.Context, msg string, err error) error {
	return &appError{
		appErrorWithoutMsg: appErrorWithoutMsg{
			next: err,
			ctx:  ctx,
		},
		msg: msg,
	}
}

func (e appError) Error() string {
	return e.msg + "; " + e.next.Error()
}