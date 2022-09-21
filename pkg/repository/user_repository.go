package repository

import "github.com/thteam47/go-identity-api/pkg/models"

type UserRepository interface {
	GetAll(userContext *models.UserContext, number int64, limit int64) ([]*models.User, error)
	Count(userContext *models.UserContext) (int64, error)
	GetOneByAttr(userContext *models.UserContext, data map[string]string) (*models.User, error)
	Create(userContext *models.UserContext, user *models.User) (*models.User, error)
	UpdatebyId(userContext *models.UserContext, user *models.User, id string) (*models.User, error)
	DeleteById(userContext *models.UserContext, id string) error
}
