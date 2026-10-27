package errs

import (
	"errors"
	"runtime"
	"runtime/debug"
)

var ErrAborted = errors.New("build aborted")

type Error struct {
	Err      error
	Stack    []byte
	Function string
	File     string
	Line     int
}

func (e *Error) Error() string { return e.Err.Error() }
func (e *Error) Unwrap() error { return e.Err }

func throw(err error) {
	if err == nil {
		return
	}
	var existing *Error
	if errors.As(err, &existing) {
		panic(existing)
	}

	e := &Error{
		Err:   err,
		Stack: debug.Stack(),
	}

	pcs := make([]uintptr, 1)
	if n := runtime.Callers(3, pcs); n > 0 {
		frame, _ := runtime.CallersFrames(pcs).Next()
		e.Function = frame.Function
		e.File = frame.File
		e.Line = frame.Line
	}

	panic(e)
}

func Handle(err *error, fn ...func(*error)) {
	r := recover()
	if r == nil {
		return
	}
	e, ok := r.(*Error)
	if !ok {
		panic(r)
	}
	*err = e
	if len(fn) > 0 && fn[0] != nil {
		fn[0](err)
	}
}

func resolve(err error, h []any) error {
	if len(h) == 0 {
		return err
	}
	switch f := h[0].(type) {
	case func() error:
		if f != nil {
			return f()
		}
	case func(error) error:
		if f != nil {
			return f(err)
		}
	}
	return err
}

func CheckE(err error, h ...any) {
	if err == nil {
		return
	}
	if e := resolve(err, h); e != nil {
		throw(e)
	}
}

func CheckV[T any](v T, err error, h ...any) T {
	if err == nil {
		return v
	}
	if e := resolve(err, h); e != nil {
		throw(e)
	}
	return v
}
