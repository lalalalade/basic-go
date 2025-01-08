package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lalalalade/basic-go/webook/internal/domain"
	"github.com/lalalalade/basic-go/webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx *gin.Context, u domain.User) error {
	return svc.repo.Create(ctx, u)
}
