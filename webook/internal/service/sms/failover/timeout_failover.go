package failover

import (
	"context"
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
	"sync/atomic"
)

var _ sms.Service = (*TimeoutFailoverSMSService)(nil)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	idx  int32
	// 连续超时个数
	cnt int32

	// 阈值
	threshold int32
}

func NewTimeoutFailoverSMSService() *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		// 触发切换 计算新的下标
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[int(idx)%len(t.svcs)]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		return err
	}
}
