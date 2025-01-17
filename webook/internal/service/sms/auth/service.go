package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lalalalade/basic-go/webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	// 加盐字符串
	key []byte
}

//func (s *SMSService) GenerateToken(ctx context.Context, tplId string) (string, error) {
//
//}

// Send 发送， 其中biz必须是线下申请的一个代表业务方的 token
func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

	var tc Claims
	// 解析token
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return s.svc.Send(ctx, tc.TplId, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	TplId string
}
