//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/lalalalade/basic-go/webook/internal/repository"
	"github.com/lalalalade/basic-go/webook/internal/repository/cache"
	"github.com/lalalalade/basic-go/webook/internal/repository/dao"
	"github.com/lalalalade/basic-go/webook/internal/service"
	"github.com/lalalalade/basic-go/webook/internal/web"
	"github.com/lalalalade/basic-go/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// 初始化DAO
		dao.NewUserDAO,
		// 初始化缓存
		cache.NewUserCache, cache.NewCodeCache,
		// 初始化repo
		repository.NewUserRepository, repository.NewCodeRepository,
		// 初始化service
		service.NewUserService, service.NewCodeService, ioc.InitSMSService, ioc.InitWechatService,
		web.NewUserHandler, web.NewOAuth2WechatHandler, ioc.NewWechatHandlerConfig,
		ioc.InitWebServer, ioc.InitMiddlewares)
	return gin.Default()
}
