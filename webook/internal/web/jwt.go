package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

type jwtHandler struct {
	// access_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
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

func newJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		rtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx"),
	}
}

func (h jwtHandler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(ctx, uid, ssid)
	return err
}

func (h jwtHandler) setJWTToken(c *gin.Context, uid int64, ssid string) error {
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
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (h jwtHandler) setRefreshToken(c *gin.Context, uid int64, ssid string) error {
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
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}
