package repoimpl

import (
	"context"
	"log"
	"time"

	"github.com/thteam47/go-identity-api/errutil"
	grpcauth "github.com/thteam47/go-identity-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-api/pkg/models"
	"github.com/thteam47/go-identity-api/pkg/repository"
	"github.com/thteam47/go-identity-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryImpl struct {
	MongoDB *mongo.Collection
}

func NewUserRepo(mongodb *mongo.Collection) repository.UserRepository {
	return &UserRepositoryImpl{
		MongoDB: mongodb,
	}
}

func (inst *UserRepositoryImpl) Create(userContext grpcauth.UserContext, user *models.User) (*models.User, error) {
	user.CreateTime = int32(time.Now().Unix())
	result, err := inst.MongoDB.InsertOne(context.Background(), user)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.InsertOne")
	}
	user.UserId = result.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (inst *UserRepositoryImpl) GetAll(userContext grpcauth.UserContext, number int32, limit int32) ([]*models.User, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	if number == 1 {
		findOptions.SetSkip(0)
		findOptions.SetLimit(int64(limit))
	} else {
		findOptions.SetSkip(int64((number - 1) * limit))
		findOptions.SetLimit(int64(limit))
	}

	cur, err := inst.MongoDB.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.Find")
	}
	var users []*models.User
	for cur.Next(context.TODO()) {
		var elem models.User
		err = cur.Decode(&elem)
		elem.UserId = elem.ID.Hex()
		if err != nil {
			return nil, errutil.Wrap(err, "Decode")
		}
		users = append(users, &elem)
	}
	return users, nil
}

func (inst *UserRepositoryImpl) Count(userContext grpcauth.UserContext) (int32, error) {
	findOptions := options.Find()

	cur, err := inst.MongoDB.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return 0, err
	}
	var users []*models.User
	for cur.Next(context.TODO()) {
		var elem *models.User
		er := cur.Decode(&elem)
		if er != nil {
			log.Fatal(err)
		}
		users = append(users, elem)
	}
	return int32(len(users)), nil
}

func (inst *UserRepositoryImpl) GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.User, error) {
	if data == nil {
		return nil, nil
	}
	findquery := []bson.M{}
	for _, key := range util.Keys[string, string](data) {
		value := ""
		if item, ok := data[key]; ok {
			value = item
			if key == "_id" {
				id, err := primitive.ObjectIDFromHex(value)
				if err != nil {
					return nil, errutil.Wrap(err, "primitive.ObjectIDFromHex")
				}
				findquery = append(findquery, bson.M{
					key: id,
				})
				continue
			}
			findquery = append(findquery, bson.M{
				key: value,
			})
		}
	}

	filters := bson.M{
		"$or": findquery,
	}

	user := &models.User{}
	err := inst.MongoDB.FindOne(context.Background(), filters).Decode(user)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.FindOne")
	}
	user.UserId = user.ID.Hex()
	return user, nil
}

func (inst *UserRepositoryImpl) UpdateOneByAttr(id string, data map[string]interface{}) error {
	idPri, _ := primitive.ObjectIDFromHex(id)

	filterSurvey := bson.M{"_id": idPri}
	dataUpdate := bson.M{}
	for _, key := range util.Keys[string, interface{}](data) {
		if value, ok := data[key]; ok {
			dataUpdate[key] = value
		}
	}
	opts := options.Update().SetUpsert(true)

	updateUser := bson.M{"$set": dataUpdate}
	_, err := inst.MongoDB.UpdateOne(context.Background(), filterSurvey, updateUser, opts)
	if err != nil {
		return errutil.Wrapf(err, "MongoDB.UpdateOne")
	}
	return nil
}

// func (inst *UserRepositoryImpl) ChangeActionUser(idUser string, role string, a []string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	var roleUser string
// 	if role == "" {
// 		roleUser = "staff"
// 	} else {
// 		roleUser = role
// 	}
// 	var actionList []string
// 	if roleUser == "admin" {
// 		actionList = append(actionList, "All Rights")
// 	} else if roleUser == "assistant" {
// 		actionList = []string{"Add Server", "Update Server", "Detail Status", "Export", "Connect", "Disconnect", "Delete Server", "Change Password"}
// 	} else {
// 		actionList = a
// 	}
// 	id, _ := primitive.ObjectIDFromHex(idUser)
// 	filterUser := bson.M{"_id": id}
// 	updateUser := bson.M{"$set": bson.M{
// 		"role":   roleUser,
// 		"action": actionList,
// 	}}
// 	_, err := inst.MongoDB.Collection(vi.GetString("collectionUser")).UpdateOne(ctx, filterUser, updateUser)
// 	if err != nil {
// 		return err
// 	}

//		return nil
//	}
func (inst *UserRepositoryImpl) UpdatebyId(userContext grpcauth.UserContext, user *models.User, id string) (*models.User, error) {
	user.UpdateTime = int32(time.Now().Unix())
	user_id, _ := primitive.ObjectIDFromHex(id)

	filterUser := bson.M{"_id": user_id}

	pByte, err := bson.Marshal(user)
	if err != nil {
		return nil, errutil.Wrap(err, "bson.Marshal")
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return nil, errutil.Wrap(err, "bson.Unmarshal")
	}
	_, err = inst.MongoDB.ReplaceOne(context.Background(), filterUser, update)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.ReplaceOne")
	}
	user.UserId = id
	return user, nil
}

// func (u *UserRepositoryImpl) ChangePassUser(idUser string, pass string) error {
// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()
// var id primitive.ObjectID
// id, _ = primitive.ObjectIDFromHex(idUser)

// passHash, _ := drive.HashPassword(pass)
// filterUser := bson.M{"_id": id}
//
//	updateUser := bson.M{"$set": bson.M{
//		"password": passHash,
//	}}
//
// _, err := u.MongoDB.UpdateOne(ctx, filterUser, updateUser)
//
//	if err != nil {
//		return err
//	}
//
//		return nil
//	}
func (inst *UserRepositoryImpl) DeleteById(userContext grpcauth.UserContext, id string) error {
	user_id, _ := primitive.ObjectIDFromHex(id)
	_, err := inst.MongoDB.DeleteOne(context.Background(), bson.M{"_id": user_id})

	if err != nil {
		return errutil.Wrap(err, "MongoDB.DeleteOne")
	}
	return nil
}
