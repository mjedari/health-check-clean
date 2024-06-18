package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/infra/utils"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strconv"
)

// todo: retry pattern for storage

type Redis struct {
	Client *redis.Client
	Config config.RedisConfig
}

func (r *Redis) CheckHealth(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) ResetConnection(ctx context.Context) error {
	newClient, err := NewRedis(r.Config)
	if err != nil {
		return err
	}
	r.Client = newClient.Client
	return nil
}

func NewRedis(conf config.RedisConfig) (*Redis, error) {
	ctx := context.TODO()

	redisRetry, err := utils.Retry(func(ctx context.Context) (any, error) {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", conf.Host, conf.Port),
			Username: conf.User,
			Password: conf.Pass,
		})

		_, err := client.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}

		return client, nil
	}, utils.RetryTimes, utils.RetryDelay)(ctx)

	if err != nil {
		return nil, err
	}
	// here we convert interface datatype to redis.Client
	client := redisRetry.(*redis.Client)

	return &Redis{Client: client}, nil
}

func (r *Redis) Add(ctx context.Context, key uint, item any) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, strconv.Itoa(int(key)), data, 0).Err()
}

func (r *Redis) Get(ctx context.Context, key uint, out any) error {
	data, err := r.Client.Get(ctx, strconv.Itoa(int(key))).Result()
	if err == redis.Nil {
		return errors.New("not found")
	} else if err != nil {
		return err
	}

	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("pointer needed")
	}

	return json.Unmarshal([]byte(data), out)
}

func (r *Redis) Remove(ctx context.Context, key uint) error {
	_, err := r.Client.Del(ctx, strconv.Itoa(int(key))).Result()
	if err == redis.Nil {
		return errors.New("not found")
	}
	return err
}

func (r *Redis) Exist(ctx context.Context, key uint) bool {
	_, err := r.Client.Get(ctx, strconv.Itoa(int(key))).Result()
	if err == redis.Nil {
		return false
	}
	return true
}
