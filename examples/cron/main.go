package main

import (
	"context"
	"time"

	"github.com/everywan/common/application"
	"github.com/everywan/common/cron"
	"github.com/everywan/common/logger"
	"github.com/everywan/common/utils"
	"github.com/xxl-job/xxl-job-executor-go"
)

func main() {
	// 创建服务, 注册 xxljob 执行器
	cronBundle, err := cron.NewCronBundle(
		cron.RegistryKey(utils.App() + "_" + utils.Env()),
	)
	utils.HandleInitError("cron-bundle", err)

	// 注册定时任务
	cronBundle.RegTask("example_task", func(ctx context.Context, param *xxl.RunReq) string {
		// todo ctx with trace
		logger.Info(ctx, "example_task start, param:%+v", param)
		params := map[string]interface{}{}
		if err := cron.ParseJsonParams(param.ExecutorParams, &params); err != nil {
			logger.Error(ctx, "parse params error, err:%+v", err)
			return "task failed"
		}

		// do something
		time.Sleep(time.Second * 2)

		logger.Info(ctx, "example_task done, param:%+v", param)
		return "task done"
	})

	app := application.New()
	app.AddBundle()
}
