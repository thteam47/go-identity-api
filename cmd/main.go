package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/go-redis/cache/v8"
	"github.com/golang/glog"
	"github.com/soheilhy/cmux"
	clienthttp "github.com/thteam47/go-identity-api/cmd/client"
	servergrpc "github.com/thteam47/go-identity-api/cmd/server"
	"github.com/thteam47/go-identity-api/pkg/configs"
	"github.com/thteam47/go-identity-api/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func client(lis net.Listener, grpc_port string, http_port string) error {
	flag.Parse()
	defer glog.Flush()
	return clienthttp.Run(lis, grpc_port, http_port)
}

func serverGrpc(lis net.Listener, handler *db.Handler) error {
	flag.Parse()
	defer glog.Flush()
	return servergrpc.Run(lis, handler)
}

func main() {
	cf, err := configs.LoadConfig()
	if err != nil {
		fmt.Println("Failed at config", err)
	}
	handler, err := db.NewHandlerWithConfig(cf)
	if err != nil {
		fmt.Println("NewHandlerWithConfig", err)
		handler = &db.Handler{
			MongoDB:    &mongo.Collection{},
			RedisCache: &cache.Cache{},
			JwtKey:     cf.KeyJwt,
		}
	}

	lis, err := net.Listen("tcp", cf.GrpcPort)
	if err != nil {
		fmt.Println("Failed to listing:", err)
	}

	fmt.Println("Server run on", cf.GrpcPort)

	m := cmux.New(lis)
	// a different listener for HTTP1
	httpL := m.Match(cmux.HTTP1Fast())
	grpcL := m.Match(cmux.HTTP2())
	go serverGrpc(grpcL, handler)
	go client(httpL, cf.GrpcPort, cf.HttpPort)
	m.Serve()
}
