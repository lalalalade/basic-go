package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "github.com/lalalalade/basic-go/webook/internal/web/jwt"
	"net/http"
)

// LoginJWTMiddlewareBuilder JWT登录校验中间件
type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
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
		tokenStr := l.ExtractToken(c)
		claims := ijwt.UserClaims{}
		// 解析token
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(ijwt.AtKey), nil
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
		err = l.CheckSession(c, claims.Ssid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
