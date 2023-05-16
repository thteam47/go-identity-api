package grpcapp

import (
	"context"

	"github.com/icrowley/fake"
	pb "github.com/thteam47/common/api/identity-api"
	"github.com/thteam47/common/entity"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-identity-api/errutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (inst *IdentityService) FakeUsers(ctx context.Context, req *pb.FakeUserRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "create", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, http.StatusText(http.StatusForbidden))
	// }
	for i := 1; i <= int(req.NumberUser); i++ {
		user := &entity.User{
			FullName: fake.FullName(),
			Email:    fake.EmailAddress(),
			Username: fake.UserName(),
			DomainId: req.Ctx.DomainId,
			Status:   "approved",
			Position: int32(i),
		}
		result, err := inst.componentsContainer.UserRepository().Create(userContext, user)
		if err != nil {
			return nil, errutil.Wrap(err, "UserRepository.Create")
		}
		err = inst.componentsContainer.IdentityAuthenService().UpdatePassword(userContext, result.UserId, req.Password)
		if err != nil {
			return nil, errutil.Wrap(err, "IdentityAuthenService.UpdatePassword")
		}
	}

	return &pb.StringResponse{}, nil
}
