package cron

import (
	"context"

	"github.com/everywan/common/logger"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)

func LogMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(ctx context.Context, param *xxl.RunReq) string {
		logger.Info(ctx, "task start %s, prams:%+v", param.ExecutorHandler, param)
		res := tf(ctx, param)
		logger.Info(ctx, "task finished %s, result:%s, prams:%+v", param.ExecutorHandler, res, param)
		return res
	}
}

func RecoverMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(ctx context.Context, param *xxl.RunReq) string {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(ctx, "task %s panic. param:%+v", param.ExecutorHandler, param)
			}
		}()
		res := tf(ctx, param)
		return res
	}
}
