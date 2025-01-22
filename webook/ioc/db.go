package ioc

import (
	"fmt"
	"github.com/lalalalade/basic-go/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化MySQL服务
func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:13316)/webook_default",
	}
	// remote下注意key不能含.
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(fmt.Errorf("mysql初始化配置失败: %v \n", err))
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
