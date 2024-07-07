package utils

import (
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"
)

// SafeRun sync run a func. If the func panics, the panic value is returned as an error.
func SafeRun(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = errors.Wrap(e, string(debug.Stack()))
			} else {
				err = fmt.Errorf("unknown panic %v\n%s", r, string(debug.Stack()))
			}
		}
	}()
	fn()
	return nil
}

// SafeRun sync run a func with error.
// If the func panics, the panic value is returned as an error.
func SafeRunWithError(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = errors.Wrap(e, string(debug.Stack()))
			} else {
				err = fmt.Errorf("unknown panic %v\n%s", r, string(debug.Stack()))
			}
		}
	}()
	return fn()
}

// SafeGo async run a func.
// If the func panics, the panic value will be handle by errHandler.
func SafeGo(fn func(), errHandler func(error)) {
	go func() {
		if err := SafeRun(fn); err != nil {
			errHandler(err)
		}
	}()
}

// SafeGoWithError async run a func with error.
// If the func panics, the panic value will be handle by errHandler.
func SafeGoWithError(fn func() error, errHandler func(error)) {
	go func() {
		if err := SafeRunWithError(fn); err != nil {
			errHandler(err)
		}
	}()
}
