package ioc

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(fmt.Errorf("redis初始化配置失败: %s \n", err))
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return redisClient
}
