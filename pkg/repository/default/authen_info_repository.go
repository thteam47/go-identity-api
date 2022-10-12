package repoimpl

import (
	"github.com/thteam47/go-identity-api/errutil"
	v1 "github.com/thteam47/go-identity-api/pkg/api-client/identity-authen-api"
	grpcauth "github.com/thteam47/go-identity-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-api/pkg/repository"
	"github.com/thteam47/go-identity-api/util"
	"google.golang.org/grpc"
)

type AuthenInfoImpl struct {
	config *util.GrpcClientConn
	client v1.IdentityAuthenServiceClient
}

func NewAuthenInfoRepo(config *util.GrpcClientConn) (repository.AuthenInfoRepository, error) {
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		return &AuthenInfoImpl{}, errutil.Wrapf(err, "grpc.Dial")
	}
	client := v1.NewIdentityAuthenServiceClient(conn)
	return &AuthenInfoImpl{
		config: config,
		client: client,
	}, nil
}

func (inst *AuthenInfoImpl) requestCtx() *v1.Context {
	return &v1.Context{
		AccessToken: inst.config.Config.AccessToken,
	}
}

func (inst *AuthenInfoImpl) UpdatePassword(userContext grpcauth.UserContext, userId string, password string) error {
	_, err := inst.client.UpdatePassword(inst.config.Context(), &v1.UpdatePasswordRequest{
		Ctx:      inst.requestCtx(),
		UserId:   userId,
		Password: password,
	})
	if err != nil {
		return errutil.Wrapf(err, "client.UpdatePassword")
	}
	return nil
}
