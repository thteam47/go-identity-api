package db

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/thteam47/go-identity-api/pkg/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Handler struct {
	MongoDB    *mongo.Collection
	RedisCache *cache.Cache
	JwtKey     string
}

func NewHandlerWithConfig(config *configs.Config) (*Handler, error) {
	mongodb, err := connectMongo(config.MongoDb)
	if err != nil {
		return nil, errors.WithMessage(err, "connectMongo")
	}
	redisCache, err := connectRedis(config.RedisCache)
	if err != nil {
		return nil, errors.WithMessage(err, "connectRedis")
	}
	return &Handler{
		MongoDB:    mongodb,
		RedisCache: redisCache,
		JwtKey:     config.KeyJwt,
	}, nil
	// return &Handler{
	// 	MongoDB:    &mongo.Collection{},
	// 	RedisCache: &cache.Cache{},
	// 	JwtKey:     config.KeyJwt,
	// }, nil
}

func connectMongo(cf configs.MongoDB) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cf.Url))
	if err != nil {
		return nil, errors.WithMessage(err, "mongo.Connect")
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, errors.WithMessage(err, "client.Ping")
	}
	// password := "admin"
	// passHash, _ := HashPassword(password)
	// user := User {
	// 	Username: "admin",
	// 	Password: passHash,
	// }
	// DB.Collection("user").InsertOne(ctx, user)
	// if err != nil {
	// 	panic(err)
	// }
	return client.Database(cf.DbName).Collection(cf.Collection), nil
}
func connectRedis(cf configs.Redis) (*cache.Cache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			cf.Address: cf.Url,
		},
	})

	err := ring.Ping(context.Background()).Err()
	if err != nil {
		return nil, errors.WithMessage(err, "ring.Ping")
	}
	memCache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return memCache, nil
}
