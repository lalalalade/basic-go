package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var _ Handler = (*RedisJWTHandler)(nil)

// 用于签名的字符串
var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

// UserClaims 短token对应的claims
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}

// RefreshClaims 长token对应的claims
type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

// SetLoginToken 同时创建长短token
func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

// SetJWTToken 创建短token
func (h *RedisJWTHandler) SetJWTToken(c *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		// 实际就是Payload（负载）部分
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: c.Request.UserAgent(),
	}
	// 生成token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

// SetRefreshToken 创建长token
func (h *RedisJWTHandler) SetRefreshToken(c *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		// 实际就是Payload（负载）部分
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

// ClearToken 删除token
func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	claims := ctx.MustGet("claims").(UserClaims)
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid),
		"", time.Hour*24*7).Err()
}

// CheckSession 判断ssid是否在redis中 在->用户已退出登录
func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("session已经无效了")
	default:
		return err
	}
}

// ExtractToken 提取token字符串
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}
