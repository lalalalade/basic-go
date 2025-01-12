//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "localhost:13316",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
