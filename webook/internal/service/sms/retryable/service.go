package retryable

import (
	"context"
	"errors"
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
)

var _ sms.Service = (*Service)(nil)

type Service struct {
	svc      sms.Service
	retryMax int
}

func NewService(svc sms.Service, retryMax int) *Service {
	return &Service{
		svc:      svc,
		retryMax: retryMax,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	if err == nil {
		return nil
	}
	cnt := 1
	for cnt < s.retryMax {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试都失败了")
}

// 设计并实现了一个高可用的短信平台
// 1.提高可用性：重试机制、客户端限流、failover（轮询、实时检测）
// 1.1 实时检测：
// 1.1.1 基于超时的实时检测（连续超时）
// 1.1.2 基于响应时间的实时检测（平均响应时间上升 20%）
// 1.1.3 基于长尾请求的实时检测（比如说，响应时间超过 1s 的请求占比超过10%）
// 1.1.4 错误率
// 2.提高安全性
// 2.1 完整的资源申请与审批流程
// 2.2 鉴权：
// 2.2.1 静态token
// 2.2.2 动态token
// 3.提高可观测性：日志、metrics、tracing、丰富完善的排查手段
