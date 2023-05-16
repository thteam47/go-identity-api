package defaultcomponent

import (
	"strings"

	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/mongoutil"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/common/pkg/mongorepository"
	"github.com/thteam47/go-identity-api/errutil"
)

type UserRepository struct {
	config         *UserRepositoryConfig
	baseRepository *mongorepository.BaseRepository
}

type UserRepositoryConfig struct {
	MongoClientWrapper *mongoutil.MongoClientWrapper `mapstructure:"mongo-client-wrapper"`
}

func NewUserRepositoryWithConfig(properties confg.Confg) (*UserRepository, error) {
	config := UserRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	mongoClientWrapper, err := mongoutil.NewBaseMongoClientWrapperWithConfig(properties.Sub("mongo-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewBaseMongoClientWrapperWithConfig")
	}
	return NewUserRepository(&UserRepositoryConfig{
		MongoClientWrapper: mongoClientWrapper,
	})
}

func NewUserRepository(config *UserRepositoryConfig) (*UserRepository, error) {
	inst := &UserRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = mongorepository.NewBaseRepository(&mongorepository.BaseRepositoryConfig{
		MongoClientWrapper: inst.config.MongoClientWrapper,
		Prototype:          entity.User{},
		MongoIdField:       "Id",
		IdField:            "UserId",
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *UserRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]entity.User, error) {
	result := []entity.User{}
	err := inst.baseRepository.FindAll(userContext, findRequest, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return result, nil
}

func (inst *UserRepository) Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error) {
	result, err := inst.baseRepository.Count(userContext, findRequest, &mongorepository.FindOptions{})
	if err != nil {
		return 0, errutil.Wrap(err, "baseRepository.Count")
	}
	return int32(result), nil
}

func (inst *UserRepository) FindById(userContext entity.UserContext, id string) (*entity.User, error) {
	result := &entity.User{}
	err := inst.baseRepository.FindOneByAttribute(userContext, "UserId", id, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByAttribute")
	}
	return result, nil
}

func (inst *UserRepository) FindByEmail(userContext entity.UserContext, email string) (*entity.User, error) {
	result := &entity.User{}
	email = strings.ToLower(strings.TrimSpace(email))
	err := inst.baseRepository.FindOneByAttribute(userContext, "Email", email, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByAttribute")
	}
	return result, nil
}

func (inst *UserRepository) FindByLoginName(userContext entity.UserContext, loginName string) (*entity.User, error) {
	result := &entity.User{}
	err := inst.baseRepository.FindOneByFindRequest(userContext, &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Operator: entity.FindRequestFilterOperatorOr,
				SubFilters: []entity.FindRequestFilter{
					entity.FindRequestFilter{
						Key:      "UserId",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    loginName,
					},
					entity.FindRequestFilter{
						Key:      "Email",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    loginName,
					},
					entity.FindRequestFilter{
						Key:      "Username",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    loginName,
					},
					entity.FindRequestFilter{
						Key:      "PhoneNumber",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    loginName,
					},
				},
			},
		},
	}, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByFindRequest")
	}
	return result, nil
}
func (inst *UserRepository) Create(userContext entity.UserContext, data *entity.User) (*entity.User, error) {
	if data.Meta == nil {
		data.Meta = map[string]string{}
	}
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	err := inst.baseRepository.Create(userContext, data, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.Create")
	}
	return data, nil
}

func (inst *UserRepository) Update(userContext entity.UserContext, data *entity.User, updateRequest *entity.UpdateRequest) (*entity.User, error) {
	if data.Meta == nil {
		data.Meta = map[string]string{}
	}
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))

	excludedProperties := []string{
		"CreatedTime", "Type",
	}

	err := inst.baseRepository.UpdateOneByAttribute(userContext, "UserId", data.UserId, data, updateRequest, &mongorepository.UpdateOptions{
		ExcludedProperties: excludedProperties,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.UpdateOneByAttribute")
	}
	return data, nil
}
func (inst *UserRepository) DeleteById(userContext entity.UserContext, id string) error {
	err := inst.baseRepository.DeleteOneByAttribute(userContext, "UserId", id)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteOneByAttribute")
	}
	return nil
}

// func (inst *UserRepository) ForgotPassword(userContext entity.UserContext, data string) (string, error) {
// 	userItem, err := inst.userService.FindByLoginName(userContext, data)
// 	if err != nil {
// 		return "", errutil.Wrapf(err, "userService.FindByLoginName")
// 	}
// 	if userItem == nil {
// 		return "", errutil.NewWithMessage("Username or password incorrect")
// 	}
// 	tokenInfo := &models.TokenInfo{
// 		AuthenticationDone: true,
// 		UserId:             userItem.UserId,
// 		Exp:                int32(time.Now().Add(5 * time.Minute).Unix()),
// 	}
// 	token, err := inst.jwtRepository.Generate(tokenInfo)
// 	if err != nil {
// 		return "", errutil.Wrapf(err, "jwtRepository.Generate")
// 	}
// 	dataMail, err := util.ParseTemplate("../util/template.html", map[string]string{
// 		"message":    "Click the link to change password.",
// 		"username":   userItem.FullName,
// 		"title":      "Forgot Password",
// 		"buttonText": "Change Passowrd Now",
// 		"link":       fmt.Sprintf("http://localhost:4200/update-password/%s", token),
// 	})
// 	if err != nil {
// 		return "", errutil.Wrapf(err, "util.ParseTemplate")
// 	}
// 	err = util.SendMail([]string{userItem.Email}, dataMail)
// 	if err != nil {
// 		return "", errutil.Wrapf(err, "util.SendMail")
// 	}
// 	return fmt.Sprintf("Click the link in your email %s to change your password", userItem.Email), nil
// }

//func (inst *UserRepository) RegisterUser(userContext entity.UserContext, username string, fullName string, email string) (string, error) {
//	userData := &v1.User{
//		FullName:   fullName,
//		Email:      email,
//		Username:   username,
//		Role:       "member",
//		CreateTime: int32(time.Now().Unix()),
//		Status:     "pending",
//	}
//	userItem, err := inst.userService.Create(userContext, userData)
//	if err != nil {
//		return "", errutil.Wrapf(err, "userService.FindByLoginName")
//	}
//	if userItem == nil {
//		return "", errutil.NewWithMessage("Username or password incorrect")
//	}
//	tokenInfo := &models.TokenInfo{
//		AuthenticationDone: true,
//		UserId:             userItem.UserId,
//		Exp:                int32(time.Now().Add(5 * time.Minute).Unix()),
//	}
//	token, err := inst.jwtRepository.Generate(tokenInfo)
//	if err != nil {
//		return "", errutil.Wrapf(err, "jwtRepository.Generate")
//	}
//	dataMail, err := util.ParseTemplate("../util/template.html", map[string]string{
//		"message":    "Click the link to verify account.",
//		"username":   fullName,
//		"title":      "Verify Account",
//		"buttonText": "Verify Now",
//		"link":       fmt.Sprintf("http://localhost:4200/verify-account/%s", token),
//	})
//	if err != nil {
//		return "", errutil.Wrapf(err, "util.ParseTemplate")
//	}
//	err = util.SendMail([]string{userItem.Email}, dataMail)
//	if err != nil {
//		return "", errutil.Wrapf(err, "util.SendMail")
//	}
//	return fmt.Sprintf("Click the link in your email %s to verify your account", userData.Email), nil
//}
