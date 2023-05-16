package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/handler"
	"github.com/thteam47/go-identity-api/errutil"
	"github.com/thteam47/go-identity-api/pkg/component"
)

type ComponentFactory struct {
	properties confg.Confg
	handle     *handler.Handler
}

func NewComponentFactory(properties confg.Confg, handle *handler.Handler) (*ComponentFactory, error) {
	inst := &ComponentFactory{
		properties: properties,
		handle:     handle,
	}

	return inst, nil
}

func (inst *ComponentFactory) CreateAuthService() *grpcauth.AuthInterceptor {
	authService := grpcauth.NewAuthInterceptor(inst.handle)
	return authService
}

func (inst *ComponentFactory) CreateUserRepository() (component.UserRepository, error) {
	userRepository, err := NewUserRepositoryWithConfig(inst.properties.Sub("user-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewAuthenInfoRepositoryWithConfig")
	}
	return userRepository, nil
}

func (inst *ComponentFactory) CreateIdentityAuthenService() (component.IdentityAuthenService, error) {
	identityAuthenService, err := NewIdentityAuthenServiceWithConfig(inst.properties.Sub("identity-authen-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewIdentityAuthenServiceWithConfig")
	}
	return identityAuthenService, nil
}