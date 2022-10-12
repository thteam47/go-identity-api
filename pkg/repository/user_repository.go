package repository

import (
	grpcauth "github.com/thteam47/go-identity-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-api/pkg/models"
)

type UserRepository interface {
	GetAll(userContext grpcauth.UserContext, number int32, limit int32) ([]*models.User, error)
	Count(userContext grpcauth.UserContext) (int32, error)
	GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.User, error)
	Create(userContext grpcauth.UserContext, user *models.User) (*models.User, error)
	UpdatebyId(userContext grpcauth.UserContext, user *models.User, id string) (*models.User, error)
	DeleteById(userContext grpcauth.UserContext, id string) error
	UpdateOneByAttr(id string, data map[string]interface{}) error 
}
