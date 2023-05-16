package component

import (
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-identity-api/errutil"
	"github.com/thteam47/go-identity-api/pkg/db"
)

type ComponentsContainer struct {
	userRepository        UserRepository
	authService           *grpcauth.AuthInterceptor
	handler               *db.Handler
	identityAuthenService IdentityAuthenService
}

func NewComponentsContainer(componentFactory ComponentFactory) (*ComponentsContainer, error) {
	inst := &ComponentsContainer{}

	var err error
	inst.authService = componentFactory.CreateAuthService()
	inst.userRepository, err = componentFactory.CreateUserRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateAuthenInfoRepository")
	}
	inst.identityAuthenService, err = componentFactory.CreateIdentityAuthenService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateIdentityAuthenService")
	}
	return inst, nil
}

func (inst *ComponentsContainer) AuthService() *grpcauth.AuthInterceptor {
	return inst.authService
}

func (inst *ComponentsContainer) UserRepository() UserRepository {
	return inst.userRepository
}

func (inst *ComponentsContainer) IdentityAuthenService() IdentityAuthenService {
	return inst.identityAuthenService
}

var errorCodeBadRequest = 400
