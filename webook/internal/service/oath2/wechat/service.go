package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lalalalade/basic-go/webook/internal/domain"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://seaandblue.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
}

func NewService(appId, appSecret string) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
	}
}

func (s *service) VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	resp, err := http.Get(target)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf(res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenID:  res.OpenId,
		UnionId: res.UnionId,
	}, nil
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
	Scope        string `json:"scope"`
}
