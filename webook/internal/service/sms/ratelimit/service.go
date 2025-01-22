package ratelimit

import (
	"context"
	"fmt"
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
	"github.com/lalalalade/basic-go/webook/pkg/ratelimit"
)

var _ sms.Service = (*RatelimitSMSService)(nil)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (r *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, "sms")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题, %w", err)
	}
	if limited {
		return errLimited
	}
	err = r.svc.Send(ctx, tpl, args, numbers...)
	return err
}
