package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lalalalade/basic-go/webook/internal/web"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// LoginJWTMiddlewareBuilder JWT登录校验中间件
type LoginJWTMiddlewareBuilder struct {
	paths []string
	cmd   redis.Cmdable
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要登录校验
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		// 现在使用jwt来校验
		tokenStr := web.ExtractToken(c)
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// err 为nil, token 不为nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != c.Request.UserAgent() {
			// 严重安全问题
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		cnt, err := l.cmd.Exists(c, fmt.Sprintf("users:ssid:%s", claims.Ssid)).Result()
		if err != nil || cnt > 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("claims", claims)
	}
}
