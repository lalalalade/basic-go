package memory

import (
	"context"
	"fmt"
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
)

var _ sms.Service = (*Service)(nil)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
