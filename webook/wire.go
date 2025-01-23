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
	ijwt "github.com/lalalalade/basic-go/webook/internal/web/jwt"
	"github.com/lalalalade/basic-go/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis, ioc.InitLogger,
		// 初始化DAO
		dao.NewUserDAO, dao.NewArticleDAO,
		// 初始化缓存
		cache.NewUserCache, cache.NewCodeCache,
		// 初始化repo
		repository.NewUserRepository, repository.NewCodeRepository, repository.NewArticleRepository,
		// 初始化service
		service.NewUserService, service.NewCodeService, service.NewArticleService,
		ioc.InitSMSService, ioc.InitWechatService,
		web.NewUserHandler, web.NewOAuth2WechatHandler, web.NewArticleHandler,
		ioc.NewWechatHandlerConfig, ijwt.NewRedisJWTHandler,
		ioc.InitWebServer, ioc.InitMiddlewares)
	return gin.Default()
}
