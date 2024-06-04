package application

import (
	"context"
	"time"

	"github.com/everywan/common/logger"
	"github.com/everywan/common/utils"
)

type options struct {
	appName      string        // app name
	logger       logger.Logger // 日志
	profilePort  int           // profile port
	stopTimeout  time.Duration // app stop timeout
	includePaths []string      // sentry path

	// lifecycle hooks
	beforeStart []func(ctx context.Context) error
	afterStart  []func(ctx context.Context) error
	beforeStop  []func(ctx context.Context) error
	afterStop   []func(ctx context.Context) error
}

func (opts *options) LoadDefault() {
	if opts.appName == "" {
		opts.appName = utils.App()
	}
	if opts.logger == nil {
		opts.logger = logger.GetLogger()
	}
	if opts.stopTimeout < 0 {
		opts.stopTimeout = DefaultStopTimeout
	}
}

type Option func(*options)

func OptionName(name string) Option {
	return func(o *options) {
		o.appName = name
	}
}
func OptionWithLogger(logger logger.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// 启用 pprof 分析
func OptionWithProfiler(profilePort int) Option {
	return func(o *options) {
		o.profilePort = profilePort
	}
}

// 收到终止信号后, 如果任务无法正常结束时强制结束的时间.
func OptionStopTimeout(timeout time.Duration) Option {
	if timeout < 0 {
		timeout = 0
	}
	return func(o *options) {
		o.stopTimeout = timeout
	}
}

// sentry path
func OptionSentryIncludePaths(paths ...string) Option {
	return func(o *options) {
		o.includePaths = paths
	}
}

func OptionBeforeStart(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.beforeStart = append(opts.beforeStart, fn)
	}
}

func OptionAfterStart(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.afterStart = append(opts.afterStart, fn)
	}
}

func OptionBeforeStop(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.beforeStop = append(opts.beforeStop, fn)
	}
}

func OptionAfterStop(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.afterStop = append(opts.afterStop, fn)
	}
}
