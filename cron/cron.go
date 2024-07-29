package cron

import (
	"context"
	"errors"
	"fmt"

	"github.com/everywan/common/application"
	"github.com/everywan/common/logger"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)

type CronBundle struct {
	xxl.Executor
	registryKey string
}

var _ application.IBundle = new(CronBundle)

func NewCronBundle(_opts ...Option) (*CronBundle, error) {
	opts := &options{}
	for _, opt := range _opts {
		opt(opts)
	}
	opts.LoadDefault()
	if opts.registryKey == "" {
		return &CronBundle{}, errors.New("must have registry_key")
	}
	client := xxl.NewExecutor(
		xxl.ServerAddr(opts.serverAddr),
		xxl.AccessToken(opts.accessToken),
		xxl.ExecutorIp(opts.executorIP),
		xxl.ExecutorPort(fmt.Sprint(opts.port)),
		xxl.RegistryKey(opts.registryKey),
		xxl.SetLogger(opts.logger),
	)
	client.Init()
	client.Use(LogMiddleware, RecoverMiddleware)
	logger.Info(context.TODO(), "create new xxljob executor. name:%s, opts:%+v", opts.registryKey, opts)
	return &CronBundle{
		Executor:    client,
		registryKey: opts.registryKey,
	}, nil
}

func (s *CronBundle) GetName() string {
	return s.registryKey
}

func (s *CronBundle) Run(ctx context.Context) {
	s.Executor.Run()
}

func (s *CronBundle) Stop(ctx context.Context) {
	s.Executor.Stop()
}
