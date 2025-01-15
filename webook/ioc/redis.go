package ioc

import (
	"github.com/lalalalade/basic-go/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
