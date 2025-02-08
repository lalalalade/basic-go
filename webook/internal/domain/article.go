package domain

// Article 可以同时表达制作库和线上库的概念？
type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
}

// Author 在帖子这个领域内是一个值对象
type Author struct {
	Id   int64
	Name string
}
