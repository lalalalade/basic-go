package article

import (
	"context"
	"github.com/lalalalade/basic-go/webook/internal/domain"
)

type ArticleReaderRepository interface {
	// Save 有就更新 没有就新建
	Save(ctx context.Context, art domain.Article) (int64, error)
}
