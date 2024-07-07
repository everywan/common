package main

import (
	"context"

	"github.com/everywan/common/application"
	"github.com/everywan/common/rpc"
	"github.com/everywan/common/rpc/grpc"
)

func main() {
	app := application.New()
	rpcBundle := rpc.NewGRPCBundle("example")
	app.AddBundle(rpcBundle)
	app.Run(context.Background())
}

// todo 添加 pb 文件补全 example
func NewClient() {
	client, err := grpc.NewClient(context.Background(), "0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
	defer client.Close()
}
