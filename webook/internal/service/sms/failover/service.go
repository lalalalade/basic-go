package failover

import (
	"context"
	"errors"
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) *FailoverSMSService {
	return &FailoverSMSService{}
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {

	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		// 输出日志 做好监控
		log.Println(err)
	}
	return errors.New("全部服务商都失败")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {

	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+uint64(length); i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			// 输出日志
		}
	}
	return errors.New("全部服务商都失败")
}
