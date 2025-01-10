package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// LoginMiddlewareBuilder 登录验证中间件，建造者模式实现
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
	gob.Register(time.Now())
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

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		// 说明还没有刷新过,刚登录，还没刷新
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		// updateTime是有的
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > 10*time.Second {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
