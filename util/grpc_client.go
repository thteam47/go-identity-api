package util

import (
	"context"
	"time"

	"github.com/thteam47/go-identity-api/pkg/configs"
	"google.golang.org/grpc"
)

type GrpcClientConn struct {
	config  configs.GrpcClientConnConfig
	conn    *grpc.ClientConn
	address string
}

var GrpcClientConnTimeout = 20 * time.Second

func NewGrpcClientConnWithConfig(config configs.GrpcClientConnConfig) (*GrpcClientConn, error) {
	if config.Timeout == 0 {
		config.Timeout = GrpcClientConnTimeout
	}

	return NewGrpcClientConn(config), nil
}

func NewGrpcClientConn(config configs.GrpcClientConnConfig) *GrpcClientConn {
	inst := &GrpcClientConn{
		config: config,
	}

	inst.address = inst.config.Address

	return inst
}

func (inst *GrpcClientConn) Context(ctx context.Context) context.Context {
	ctxx, cancel := context.WithTimeout(context.Background(), inst.config.Timeout)
	defer cancel()
	return ctxx
}

func (inst *GrpcClientConn) Stop() {
	if inst.conn != nil {
		inst.conn.Close()
		inst.conn = nil
	}
}