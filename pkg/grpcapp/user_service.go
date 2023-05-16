package grpcapp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/thteam47/common/entity"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/pkg/adapter"
	"github.com/thteam47/go-identity-api/errutil"
	"github.com/thteam47/go-identity-api/pkg/component"
	pb "github.com/thteam47/common/api/identity-api"
	"github.com/thteam47/go-identity-api/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityService struct {
	pb.IdentityServiceServer
	componentsContainer *component.ComponentsContainer
}

func NewIdentityService(componentsContainer *component.ComponentsContainer) *IdentityService {
	return &IdentityService{
		componentsContainer: componentsContainer,
	}
}
func getUser(item *pb.User) (*entity.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &entity.User{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

func getUsers(items []*pb.User) ([]*entity.User, error) {
	users := []*entity.User{}
	for _, item := range items {
		user, err := getUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func makeUser(item *entity.User) (*pb.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &pb.User{}
	err := util.ToMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return user, nil
}

func makeUsers(items []entity.User) ([]*pb.User, error) {
	users := []*pb.User{}
	for _, item := range items {
		user, err := makeUser(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func (inst *IdentityService) Create(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "create", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	user, err := getUser(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getUser")
	}
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Username = strings.TrimSpace(strings.ToLower(user.Username))
	result, err := inst.componentsContainer.UserRepository().Create(userContext, user)
	if err != nil {
		return nil, errutil.Wrap(err, "UserRepository.Create")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return &pb.UserResponse{
		Data: item,
	}, nil
}

func (inst *IdentityService) GetById(ctx context.Context, req *pb.StringRequest) (*pb.UserResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	result, err := inst.componentsContainer.UserRepository().FindById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.FindById")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return &pb.UserResponse{
		Data: item,
	}, nil
}

func (inst *IdentityService) GetByLoginName(ctx context.Context, req *pb.StringRequest) (*pb.UserResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "create", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	result, err := inst.componentsContainer.UserRepository().FindByLoginName(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.FindByLoginName")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return &pb.UserResponse{
		Data: item,
	}, nil
}

func (inst *IdentityService) GetByEmail(ctx context.Context, req *pb.StringRequest) (*pb.UserResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	result, err := inst.componentsContainer.UserRepository().FindByEmail(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.GetById")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return &pb.UserResponse{
		Data: item,
	}, nil
}

func (inst *IdentityService) GetAll(ctx context.Context, req *pb.ListRequest) (*pb.ListUserResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, "Fobbiden!")
	// }
	findRequest, err := adapter.GetFindRequest(req, req.RequestPayload)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprint(err))
	}
	result, err := inst.componentsContainer.UserRepository().FindAll(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.FindAll")
	}
	item, err := makeUsers(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUsers")
	}
	count, err := inst.componentsContainer.UserRepository().Count(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.Count")
	}

	return &pb.ListUserResponse{
		Data:  item,
		Total: count,
	}, nil
}

func (inst *IdentityService) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "update", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	user, err := getUser(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getUser")
	}
	result, err := inst.componentsContainer.UserRepository().Update(userContext, user, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.UpdatebyId")
	}
	item, err := makeUser(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return &pb.UserResponse{
		Data: item,
	}, nil
}
func (inst *IdentityService) DeleteById(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "delete", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.componentsContainer.UserRepository().DeleteById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.DeleteById")
	}
	return &pb.StringResponse{}, nil
}

func (inst *IdentityService) ApproveUser(ctx context.Context, req *pb.ApproveUserRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "identity-api", "update", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}

	user, err := inst.componentsContainer.UserRepository().FindById(userContext, req.UserId)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.FindById")
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, http.StatusText(http.StatusNotFound))
	}
	user.Status = req.Status
	_, err = inst.componentsContainer.UserRepository().Update(userContext, user, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.Update")
	}
	return &pb.StringResponse{}, nil
}
