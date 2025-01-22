package ioc

import (
	"github.com/lalalalade/basic-go/webook/pkg/logger"
	"go.uber.org/zap"
)

// InitLogger 初始化日志
func InitLogger() logger.LoggerV1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
