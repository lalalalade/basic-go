package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}

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

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 算出签名， 返回字符串
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

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
	// 算出签名， 返回字符串
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	claims := ctx.MustGet("claims").(*UserClaims)
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid),
		"", time.Hour*24*7).Err()
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err
}

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}
