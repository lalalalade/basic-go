package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

// MiddlewareBuilder
// 1. 小心日志内容过多，url可能过长，请求、响应体都可能很大
// 2. 用户可能换用不同的日志框架，所以要有足够的灵活性
// 3. 考虑动态开关，结合监听配置文件，小心并发安全
type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}
func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		url := c.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: c.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && c.Request.Body != nil {
			body, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			c.Writer = responseWriter{
				al:             al,
				ResponseWriter: c.Writer,
			}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(c, al)
		}()
		// 执行到业务逻辑
		c.Next()

	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	Method   string
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
