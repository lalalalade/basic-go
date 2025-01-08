package domain

import "time"

// User 领域对象，是 DDD 中的entity
// BO(business object)
type User struct {
	Email      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
}

//type Address struct {
//}
