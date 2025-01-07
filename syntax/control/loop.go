package main

import "fmt"

func Loop1() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	// 这样也可以
	for i := 0; i < 10; {
		println(i)
		i++
	}
}

func Loop2() {
	i := 0
	for i < 10 {
		println(i)
		i++
	}
}

// Loop3 是无限循环
func Loop3() {
	for {
		println("hello")
	}
}

func LoopBreak() {
	i := 0
	for {
		if i >= 10 {
			break
		}
		println(i)
		i++
	}
}

func LoopContinue() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		println(i)
	}
}

func LoopRange() {
	println("遍历数组")
	arr := [3]string{"11", "12", "13"}
	for i, val := range arr {
		println(i, val)
	}
	for i := range arr {
		println(i, arr[i])
	}

	println("遍历切片")
	slice := []string{"a", "b", "c"}
	for i, val := range slice {
		println(i, val)
	}
	for i := range slice {
		println(i, slice[i])
	}

	println("遍历map")
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	for k, v := range m {
		println(k, v)
	}
	for k := range m {
		println(k, m[k])
	}
}

type User struct {
	Name string
}

// LoopBug 不算BUG, 算是Go循环里的踩坑点
func LoopBug() {
	users := []User{
		{
			Name: "Tom",
		},
		{
			Name: "Jerry",
		},
	}
	m := make(map[string]*User, 2)
	for _, u := range users {
		m[u.Name] = &u
	}

	for k, v := range m {
		fmt.Printf("name: %s, user: %v", k, v)
	}
}
