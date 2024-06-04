package main

import (
	"context"
	"fmt"

	"github.com/everywan/common/application"
)

func main() {
	// 启动 app 服务
	app := application.New(
		application.OptionName("example"),
		application.OptionWithProfiler(9999),
		application.OptionBeforeStart(func(ctx context.Context) error {
			fmt.Println("example app start...")
			return nil
		}),
		application.OptionAfterStop(func(ctx context.Context) error {
			fmt.Println("example app stop...")
			return nil
		}),
	)
	// app.AddBundle()
	app.Run(context.TODO())
}
