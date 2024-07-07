package grpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc 基础客户端, 预设配置和工具.
type IClient interface {
	GetConn() *grpc.ClientConn
	Close() error
}

type client struct {
	conn *grpc.ClientConn
}

var _ IClient = (*client)(nil)

func NewClient(ctx context.Context, target string, opts ...grpc.DialOption) (*client, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 500 * time.Millisecond,
		}),
	}
	opts = append(opts, defaultOpts...)
	overwriteOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(withClientTelemetry),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	}
	opts = append(opts, overwriteOpts...)

	// todo 未来自定义地址时, 支持 resolver 解析 target
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return &client{}, errors.Wrapf(err, "grpc.NewClient error. target:%s", target)
	}
	return &client{conn: conn}, nil
}

// todo 支持切面方法
func withClientTelemetry(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (client *client) GetConn() *grpc.ClientConn {
	return client.conn
}

func (client *client) Close() error {
	err := client.conn.Close()
	if err != nil {
		return errors.Wrap(err, "close grpc conn error")
	}
	return nil
}
