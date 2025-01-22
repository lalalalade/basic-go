package ioc

import (
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
	"github.com/lalalalade/basic-go/webook/internal/service/sms/memory"
	"github.com/lalalalade/basic-go/webook/internal/service/sms/ratelimit"
	"github.com/lalalalade/basic-go/webook/internal/service/sms/retryable"
	limiter "github.com/lalalalade/basic-go/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

// InitSMSService 初始化短信服务
func InitSMSService(cmd redis.Cmdable) sms.Service {
	svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
		limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	return retryable.NewService(svc, 3)
}
