package grpcapp

import (
	"context"

	"github.com/thteam47/go-identity-api/errutil"
	grpcauth "github.com/thteam47/go-identity-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-api/pkg/models"
	"github.com/thteam47/go-identity-api/pkg/pb"
	repository "github.com/thteam47/go-identity-api/pkg/repository"
	"github.com/thteam47/go-identity-api/util"
)

type IdentityService struct {
	pb.IdentityServiceServer
	userRepository repository.UserRepository
	authRepository *grpcauth.AuthInterceptor
}

func NewIdentityService(userRepository repository.UserRepository, authRepository *grpcauth.AuthInterceptor) *IdentityService {
	return &IdentityService{
		userRepository: userRepository,
		authRepository: authRepository,
	}
}

func getUser(item *pb.User) (*models.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &models.User{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

func getUsers(items []*pb.User) ([]*models.User, error) {
	users := []*models.User{}
	for _, item := range items {
		user, err := getUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func makeUser(item *models.User) (*pb.User, error) {
	user := &pb.User{}
	err := util.ToMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return user, nil
}

func makeUsers(items []*models.User) ([]*pb.User, error) {
	users := []*pb.User{}
	for _, item := range items {
		user, err := makeUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func (inst *IdentityService) Create(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "create")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	user, err := getUser(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getUser")
	}
	result, err := inst.userRepository.Create(nil, user)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.Create")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return item, nil
}

func (inst *IdentityService) GetById(ctx context.Context, req *pb.StringRequest) (*pb.User, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "get")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	result, err := inst.userRepository.GetOneByAttr(nil, map[string]string{
		"_id": req.Value,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.GetById")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return item, nil
}

func (inst *IdentityService) GetByLoginName(ctx context.Context, req *pb.StringRequest) (*pb.User, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "get")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	result, err := inst.userRepository.GetOneByAttr(nil, map[string]string{
		"username": req.Value,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.GetById")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return item, nil
}

func (inst *IdentityService) GetByEmail(ctx context.Context, req *pb.StringRequest) (*pb.User, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "get")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	result, err := inst.userRepository.GetOneByAttr(nil, map[string]string{
		"email": req.Value,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.GetById")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return item, nil
}

func (inst *IdentityService) GetAll(ctx context.Context, req *pb.ListRequest) (*pb.ListUserResponse, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "get")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	number := 1
	limit := 10
	if req != nil && req.Data != nil {
		if req.Data.Limit > 0 {
			limit = int(req.Data.Limit)
		}
		if req.Data.Number >= 1 {
			number = int(req.Data.Number)
		}
	}
	result, err := inst.userRepository.GetAll(nil, int64(number), int64(limit))
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.GetAll")
	}
	item, err := makeUsers(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUsers")
	}
	count, err := inst.userRepository.Count(nil)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.Count")
	}

	return &pb.ListUserResponse{
		Data:  item,
		Total: count,
	}, nil
}

func (inst *IdentityService) UpdatebyId(ctx context.Context, req *pb.UpdateUserRequest) (*pb.StringResponse, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "update")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	user, err := getUser(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getUser")
	}
	_, err = inst.userRepository.UpdatebyId(nil, user, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.UpdatebyId")
	}
	return &pb.StringResponse{}, nil
}

func (inst *IdentityService) DeleteById(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-api:user", "delete")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	err := inst.userRepository.DeleteById(nil, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.DeleteById")
	}
	return &pb.StringResponse{}, nil
}
