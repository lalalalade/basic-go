package repository

type UserRepository struct {
}

func (r *UserRepository) FindById(int64) {
	// 先从cache找
	// 再从 dao 找
	// 找到了回写cache
}
