package ioc

import (
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
	"github.com/lalalalade/basic-go/webook/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	return memory.NewService()
}
