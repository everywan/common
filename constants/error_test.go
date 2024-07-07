package constants

import (
	"testing"

	"github.com/pkg/errors"
)

func TestError(t *testing.T) {
	err := errors.Wrap(ErrorNotFound, "id=101")
	t.Run("get_base_error", func(t *testing.T) {
		switch errors.Cause(err) {
		case ErrorNotFound:
		default:
			t.Error("base error not match")
		}
	})
	t.Run("get_error_code", func(t *testing.T) {
		if GetErrorCode(err, 0) != 404 {
			t.Error("error code not match")
		}
		// err==nil 时返回0
		if GetErrorCode(nil, 500) != 0 {
			t.Error("error nil code not 0")
		}
		// 非 base_error 返回默认code
		if GetErrorCode(errors.New("udf error"), 500) != 500 {
			t.Error("error code not match")
		}
	})
}
