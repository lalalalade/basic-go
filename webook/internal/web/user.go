package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lalalalade/basic-go/webook/internal/domain"
	"github.com/lalalalade/basic-go/webook/internal/service"
	"net/http"
	"time"
)

const biz = "login"

// 确保 UserHandler 实现了 handler 接口
var _ handler = &UserHandler{}

// UserHandler 定义用户相关路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	jwtHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+(-+.\\w+)*@\\w+(-.\\w+)*.\\w+(-.\\w+)*$"
		// 强密码(必须包含大小写字母和数字的组合，可以使用特殊字符，长度在8-10之间)：
		passwordRegexPattern = "^(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).{8,10}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) SendLoginSMSCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	err := u.codeSvc.Send(c, biz, req.Phone)
	switch err {
	case nil:
		c.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		c.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后重试",
		})
	default:
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err = u.setJWTToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "验证码校验成功",
	})
}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
		Password        string `json:"password" binding:"required"`
	}

	var req SignUpReq
	// Bind方法会根据Content-Type来解析数据到req里面
	// 解析错了，就会直接写回一个400错误
	if err := c.Bind(&req); err != nil {
		return
	}
	isEmail, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isEmail {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "邮箱格式不正确",
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "两次输入的密码不一致",
		})
		return
	}
	isPassword, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isPassword {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "密码必须包含大小写字母和数字的组合，可以使用特殊字符，长度在8-10之间",
		})
		return
	}
	// 调用service方法
	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicate) {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "邮箱冲突",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Msg: "注册成功",
	})
}

func (u *UserHandler) LoginJWT(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "用户名或密码不对",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 登录成功了
	// 用jwt设置登录态
	// 生成一个jwt token
	if err = u.setJWTToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
	return
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "用户名或密码不对",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 登录成功了
	// 设置session
	sess := sessions.Default(c)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 30,
	})
	sess.Save()
	c.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
	return
}

func (u *UserHandler) Logout(c *gin.Context) {
	sess := sessions.Default(c)

	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	c.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}
func (u *UserHandler) Edit(c *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname" binding:"required"`
		Birthday string `json:"birthday" binding:"required"`
		Info     string `json:"info" binding:"required"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	if req.Nickname == "" {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "昵称不能为空",
		})
		return
	}
	if len(req.Info) > 1024 {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "简介过长",
		})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "日期格式不对",
		})
		return
	}
	uc := c.MustGet("claims").(*UserClaims)
	err = u.svc.UpdateNoneSensitiveInfo(c, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		Info:     req.Info,
	})
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (u *UserHandler) ProfileJWT(c *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		Info     string
	}
	uc := c.MustGet("claims").(UserClaims)
	user, err := u.svc.Profile(c, uc.Uid)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Profile{
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Birthday: user.Birthday.Format(time.DateOnly),
		Info:     user.Info,
	})
}
