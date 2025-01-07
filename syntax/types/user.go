package types

import "fmt"

type User struct {
	Name string
	Age  int
}

func (u User) ChangeName(name string) {
	fmt.Printf("u address %p \n", &u)
	u.Name = name
}

func (u User) ChangeAge(age int) {
	u.Age = age
}

func ChangeUser() {
	u1 := User{
		Name: "Tom",
		Age:  18,
	}
	fmt.Printf("%+v \n", u1)
	fmt.Printf("u1 address %p \n", &u1)

	u1.ChangeName("Jerry")
	u1.ChangeAge(24)
	fmt.Printf("%+v \n", u1)

	u2 := &User{
		Name: "小海",
		Age:  18,
	}
	fmt.Printf("%+v \n", u2)
	fmt.Printf("u2 address %p \n", &u2)
	u2.ChangeName("Jerry")
	u2.ChangeAge(24)
	fmt.Printf("%+v \n", u2)
}

func NewUser() {
	// u1是指向一个User对象的指针
	u1 := &User{}
	println(u1)

	// u2中的字段都是零值
	u2 := User{}
	println(u2)
	// 修改u2的字段
	u2.Name = "Jerry"

	// u3中的字段也都是零值
	var u3 User
	println(u3)

	// 初始化的同时，还赋值了 Name
	var u4 User = User{
		Name: "Tom",
	}
	println(u4)

	// 没有指定字段名，按照字段顺序赋值
	// 必须全部赋值
	var u5 User = User{"Tom", 18}
	println(u5)
}
