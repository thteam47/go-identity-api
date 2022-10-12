package servergrpc

import (
	"net"

	"github.com/thteam47/go-identity-api/pkg/db"
	"github.com/thteam47/go-identity-api/pkg/grpcapp"
	grpcauth "github.com/thteam47/go-identity-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-api/pkg/pb"
	repoimpl "github.com/thteam47/go-identity-api/pkg/repository/default"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(lis net.Listener, handler *db.Handler) error {
	userRepository := repoimpl.NewUserRepo(handler.MongoDB)
	authRepository := grpcauth.NewAuthInterceptor(handler.JwtKey)
	serverOptions := []grpc.ServerOption{}
	s := grpc.NewServer(serverOptions...)
	pb.RegisterIdentityServiceServer(s, grpcapp.NewIdentityService(userRepository, authRepository))
	reflection.Register(s)
	return s.Serve(lis)
}
