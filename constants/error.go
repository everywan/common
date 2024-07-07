package constants

import "github.com/pkg/errors"

// base error 封装, 并附带错误码, 以 http.statusCode 为标准
// 上游程序通过 errors.Cause 获取判断, 参考 error_test.go
var (
	ErrorBadRequest     = errors.New("400_error_bad_request") // 请求错误, 如参数错误等.
	ErrorForbidden      = errors.New("403_error_forbidden")
	ErrorNotFound       = errors.New("404_error_not_found")
	ErrorTimeout        = errors.New("408_error_timeout")
	ErrorContenxtCancel = errors.New("408_error_context_cancel")
	ErrorDeprecated     = errors.New("410_error_deprecated") // 被弃用
	ErrorLocked         = errors.New("423_error_locked")
	ErrorTooManyRequest = errors.New("429_error_too_many_request")
	ErrorUnkownFound    = errors.New("500_error_unkown")
	ErrorNotImplemented = errors.New("501_error_not_implemented")
	ErrorUnavailable    = errors.New("503_error_unavailable") // 服务不可用, 如被断路器拦截等
)

// 获取错误码.
// 当 error==nil 时返回 0, 当 err 不是 base_error 时返回 defaultCode.
func GetErrorCode(err error, defaultCode int) int {
	if err == nil {
		return 0
	}
	str := errors.Cause(err).Error()
	if len(str) < 3 {
		return defaultCode
	}
	code := 0
	for _, c := range str[:3] {
		if c < '0' || c > '9' {
			return defaultCode
		}
		code = code*10 + int(c-'0')
	}
	return code
}
