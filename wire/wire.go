//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lalalalade/basic-go/wire/repository"
	"github.com/lalalalade/basic-go/wire/repository/dao"
)

func InitRepository() *repository.UserRepository {
	// 传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
