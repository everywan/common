package utils

import (
	"errors"
	"testing"
)

func TestSafeRun(t *testing.T) {
	t.Run("test_panic_err", func(t *testing.T) {
		expectErr := errors.New("should error")
		fn := func() { panic(expectErr) }
		err := SafeRun(fn)
		if !errors.Is(err, expectErr) {
			t.Error("SafeRun")
		}
	})
	t.Run("test_panic_other", func(t *testing.T) {
		fn := func() { panic("panic") }
		_ = SafeRun(fn)
	})

}
