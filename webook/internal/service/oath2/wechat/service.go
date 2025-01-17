package wechat

import (
	"context"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

type service struct {
	appId string
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	//urlPattern := "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	//redirectURI := "https://seaandblue.com/oauth2/wechat/callback"
	//return fmt.Sprintf(urlPattern, s.appId)
}
