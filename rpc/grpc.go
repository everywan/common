package rpc

import (
	"context"
	"net"

	"github.com/everywan/common/application"
	xgrpc "github.com/everywan/common/rpc/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GRPCBundle struct {
	name       string
	Server     *grpc.Server
	listenAddr string

	serverFactory func() *grpc.Server
}

var _ application.IBundle = new(GRPCBundle)

func NewGRPCBundle(name string, opts ...GRPCOption) *GRPCBundle {
	s := &GRPCBundle{
		name:          name,
		listenAddr:    "0.0.0.0:8000",
		serverFactory: func() *grpc.Server { return xgrpc.NewServer() },
	}
	for _, opt := range opts {
		opt(s)
	}

	s.Server = s.serverFactory()

	return s
}

func (s *GRPCBundle) GetName() string {
	return s.name
}

func (s *GRPCBundle) Run(ctx context.Context) {
	addr, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		err := errors.Wrap(err, "listen failed")
		panic(err)
	}
	err = s.Server.Serve(addr)
	if err != nil {
		panic(err)
	}
}

func (s *GRPCBundle) Stop(ctx context.Context) {
	s.Server.Stop()
}
