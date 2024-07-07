package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	overwriteOpts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      5 * time.Minute,
			MaxConnectionAgeGrace: 10 * time.Second,
			Time:                  time.Second,
			Timeout:               time.Millisecond * 100,
		}),
		grpc.ChainUnaryInterceptor(WithServerTelemetry),
	}
	opts = append(opts, overwriteOpts...)

	s := grpc.NewServer(opts...)
	// 允许客户端查询服务端支持的服务和方法
	reflection.Register(s)
	// k8s 检活服务
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	return s
}

// todo 支持切面方法
func WithServerTelemetry(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	return handler(ctx, req)
}
