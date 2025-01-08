package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要登录校验
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(c)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
