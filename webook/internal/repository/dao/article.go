package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// Article 制作库的
type Article struct {
	Id      int64  `gorm:"primaryKey;autoIncrement"`
	Title   string `gorm:"type:varchar(1024);"`
	Content string `gorm:"type=BLOB"`
	// 如何设计索引
	// 查询场景？
	// 对于创作者来说，看草稿箱，看到所有自己的文章 where author_id = 123
	// 单独查询某一篇   where id = 1
	// 最佳选择 在author_id 和 ctime 上创建联合索引
	AuthorId int64 `gorm:"index=aid_ctime"`
	Ctime    int64 `gorm:"index=aid_ctime"`
	Utime    int64
}
