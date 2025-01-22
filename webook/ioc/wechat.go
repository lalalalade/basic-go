package ioc

import (
	"github.com/lalalalade/basic-go/webook/internal/service/oauth2/wechat"
	"github.com/lalalalade/basic-go/webook/internal/web"
	logger2 "github.com/lalalalade/basic-go/webook/pkg/logger"
)

// InitWechatService 初始化微信服务
func InitWechatService(l logger2.LoggerV1) wechat.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_ID")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_SECRET")
	//}
	appId := "12341"
	appSecret := "sdfskm"
	return wechat.NewService(appId, appSecret, l)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
