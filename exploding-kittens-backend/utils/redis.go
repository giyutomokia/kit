package utils

import (
    "context"
    "github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
}

func GetRedisClient() *redis.Client {
    return rdb
}

func Set(ctx context.Context, key string, value interface{}) error {
    return rdb.Set(ctx, key, value, 0).Err()
}

func Get(ctx context.Context, key string) (string, error) {
    val, err := rdb.Get(ctx, key).Result()
    if err != nil {
        return "", err
    }
    return val, nil
}
