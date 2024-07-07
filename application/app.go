package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/everywan/common/logger"
)

type Application interface {
	Name() string
	// 启动程序. 程序会在收到终止信号 或 所有任务执行完毕 或 context.Done() 时终止任务.
	// 终止信号: syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT
	Run(context.Context)
	// 添加任务
	AddBundle(...IBundle)
}

// 任务, http、rest 服务, corn job 都是任务. 一个程序由多个任务组成.
type IBundle interface {
	GetName() string
	Run(context.Context)
	Stop(context.Context)
}

// Application 的默认实现
type defaultApplication struct {
	name   string
	config *appConfig
	logger logger.Logger

	// bundles
	bundles []IBundle

	// app lifecycle hooks
	beforeStart []func(ctx context.Context) error
	afterStart  []func(ctx context.Context) error
	beforeStop  []func(ctx context.Context) error
	afterStop   []func(ctx context.Context) error
}

var _ Application = new(defaultApplication)

func New(_opts ...Option) Application {
	opts := &options{}
	for _, opt := range _opts {
		opt(opts)
	}
	opts.LoadDefault()

	return &defaultApplication{
		name: opts.appName,
		config: &appConfig{
			stopTimeout: opts.stopTimeout,
			profilePort: opts.profilePort,
		},
		logger:      opts.logger,
		beforeStart: opts.beforeStart,
		afterStart:  opts.afterStart,
		beforeStop:  opts.beforeStop,
		afterStop:   opts.afterStop,
	}

}

func (app *defaultApplication) Name() string {
	return app.name
}

func (app *defaultApplication) Run(ctx context.Context) {
	app.logger.Infof(ctx, "Run application [%s]", app.name)

	// 设置 sentry

	if app.config.profilePort > 0 {
		go func() {
			err := http.ListenAndServe(fmt.Sprintf(":%d", app.config.profilePort), nil)
			if err != nil {
				app.logger.Error(ctx, "Start pprof error: ", err)
			}
		}()
	}

	app.runBeforeStart(ctx)
	bundleFinishCtx := app.runAllBundles(ctx)
	app.runAfterStart(ctx)
	app.logger.Infof(ctx, "Run application [%s] success", app.name)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	select {
	case <-bundleFinishCtx.Done():
		app.logger.Info(ctx, "All bundle finished.")
	case <-quit:
		app.logger.Info(ctx, "Shutdown signal received.")
	}

	app.runBeforeStop(ctx)
	bundleStopCtx := app.stopAllBundles(ctx)
	shutdownTimeout := time.After(app.config.stopTimeout)
	select {
	case <-bundleStopCtx.Done():
		app.logger.Info(ctx, "Application stopped")
	case <-shutdownTimeout:
		app.logger.Info(ctx, "Shutdown timeout, force stop application")
	}
	app.runAfterStop(ctx)

	// 清除其他资源(如有)

	app.logger.Info(ctx, "Bye!")
}

func (app *defaultApplication) runBeforeStart(ctx context.Context) {
	for _, fn := range app.beforeStart {
		if err := fn(ctx); err != nil {
			app.logger.Fatalf(ctx, "before bundles start fn fatal:%v", err)
		}
	}
}

func (app *defaultApplication) runAfterStart(ctx context.Context) {
	for _, fn := range app.afterStart {
		if err := fn(ctx); err != nil {
			app.logger.Errorf(ctx, "after bundles start fn error:%v", err)
		}
	}
}

func (app *defaultApplication) runBeforeStop(ctx context.Context) {
	for _, fn := range app.beforeStop {
		if err := fn(ctx); err != nil {
			app.logger.Errorf(ctx, "before bundles stop fn error:%v", err)
		}
	}
}

func (app *defaultApplication) runAfterStop(ctx context.Context) {
	for _, fn := range app.afterStop {
		if err := fn(ctx); err != nil {
			app.logger.Errorf(ctx, "after bundles stop fn error:%v", err)
		}
	}
}

func (app *defaultApplication) runAllBundles(ctx context.Context) context.Context {
	finishCtx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for _, bundle := range app.bundles {
		bundle := bundle
		wg.Add(1)
		app.logger.Infof(ctx, "bundle %s starting..", bundle.GetName())
		// 由 bundle 自己决定是否捕获异常, app 不做处理.
		go func() {
			defer wg.Done()
			bundle.Run(finishCtx)
			app.logger.Infof(ctx, "bundle %s start suceess.", bundle.GetName())
		}()
	}

	go func() {
		wg.Wait()
		cancel()
	}()

	return finishCtx
}

func (app *defaultApplication) stopAllBundles(ctx context.Context) context.Context {
	// todo add trace
	finishCtx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for _, bundle := range app.bundles {
		bundle := bundle
		wg.Add(1)
		app.logger.Infof(ctx, "bundle %s stoping.", bundle.GetName())
		// 由 bundle 自己决定是否捕获异常, app 不做处理.
		go func() {
			defer wg.Done()
			bundle.Stop(finishCtx)
			app.logger.Infof(ctx, "bundle %s stop suceess.", bundle.GetName())
		}()
	}

	go func() {
		wg.Wait()
		cancel()
		app.logger.Info(ctx, "All bundle stopped")
	}()

	return finishCtx
}

func (app *defaultApplication) AddBundle(bundles ...IBundle) {
	app.bundles = append(app.bundles, bundles...)
}
