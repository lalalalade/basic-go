package ioc

import (
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
	"github.com/lalalalade/basic-go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
