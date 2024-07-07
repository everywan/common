package rpc

import "google.golang.org/grpc"

type GRPCOption func(bundle *GRPCBundle)

func GRPCListen(listenAddr string) GRPCOption {
	return func(s *GRPCBundle) {
		s.listenAddr = listenAddr
	}
}

func GRPCServerFactory(f func() *grpc.Server) GRPCOption {
	return func(s *GRPCBundle) {
		s.serverFactory = f
	}
}
