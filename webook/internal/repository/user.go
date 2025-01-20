package repository

import (
	"context"
	"database/sql"
	"github.com/lalalalade/basic-go/webook/internal/domain"
	"github.com/lalalalade/basic-go/webook/internal/repository/cache"
	"github.com/lalalalade/basic-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToPo(u))
}

func (r *CachedUserRepository) Update(ctx context.Context, u domain.User) error {
	err := r.dao.UpdateNonZeroFields(ctx, r.domainToPo(u))
	if err != nil {
		return err
	}
	return r.cache.Delete(ctx, u.Id)
}
func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	// 缓存有数据
	if err == nil {
		return u, nil
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.poToDomain(ue)
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志，做监控
		}
	}()
	return u, err
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.poToDomain(u), err
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.poToDomain(u), err
}

func (r *CachedUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return r.poToDomain(u), err
}

func (r *CachedUserRepository) poToDomain(u dao.User) domain.User {
	var birthday time.Time
	if u.Birthday.Valid {
		birthday = time.UnixMilli(u.Birthday.Int64)
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		WechatInfo: domain.WechatInfo{
			UnionId: u.WechatUnionID.String,
			OpenID:  u.WechatOpenID.String,
		},
		Info:       u.Info.String,
		Nickname:   u.Nickname.String,
		Birthday:   birthday,
		CreateTime: time.UnixMilli(u.CreateTime),
	}
}
func (r *CachedUserRepository) domainToPo(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		Birthday: sql.NullInt64{
			Int64: u.Birthday.UnixMilli(),
			Valid: !u.Birthday.IsZero(),
		},
		Nickname: sql.NullString{
			String: u.Nickname,
			Valid:  u.Nickname != "",
		},
		Info: sql.NullString{
			String: u.Info,
			Valid:  u.Info != "",
		},
		Password: u.Password,
	}
}
