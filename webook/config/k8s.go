//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(webook-test-mysql:11309)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-test-redis:11479",
	},
}
