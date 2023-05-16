package component

import (
	"github.com/thteam47/common/entity"
)

type UserRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]entity.User, error)
	Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error)
	FindById(userContext entity.UserContext, id string) (*entity.User, error)
	FindByEmail(userContext entity.UserContext, email string) (*entity.User, error)
	FindByLoginName(userContext entity.UserContext, loginName string) (*entity.User, error)
	Create(userContext entity.UserContext, user *entity.User) (*entity.User, error)
	Update(userContext entity.UserContext, data *entity.User, updateRequest *entity.UpdateRequest) (*entity.User, error)
	DeleteById(userContext entity.UserContext, id string) error
}
