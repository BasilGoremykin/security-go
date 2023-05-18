package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-security/internal/model"
	"go-security/internal/persistence"
	"go-security/internal/persistence/cache"
	"log"
)

type redisCache struct {
	Client *redis.Client
}

var ctx = context.Background()

func NewRedisCache(addr, password string, db int) (cache.UserPermissionCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &redisCache{Client: client}, nil
}

func (r *redisCache) GetPermissions(userId int64) map[model.Permission]struct{} {
	data, err := r.Client.Get(ctx, fmt.Sprintf("%d", userId)).Bytes()
	if err != nil {
		log.Println(err)
		return nil
	}
	var permissionsArr []model.Permission
	err = json.Unmarshal(data, &permissionsArr)
	if err != nil {
		log.Println(err)
		return nil
	}
	return persistence.MapModelArrayToModelSet(permissionsArr)
}

func (r *redisCache) SetPermissions(userId int64, permissions map[model.Permission]struct{}) error {
	permissionsJSON, err := json.Marshal(persistence.MapModelSetToArray(permissions))
	if err != nil {
		log.Println(err)
		return err
	}
	err = r.Client.Set(ctx, fmt.Sprintf("%d", userId), permissionsJSON, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *redisCache) DeletePermissions(userId int64) error {
	delCmd := r.Client.Del(ctx, fmt.Sprintf("%d", userId))
	if delCmd != nil {
		return delCmd.Err()
	}
	return nil
}
