package main

import "fmt"

func Array() {
	// 直接初始化一个三个元素的数组，大括号里面只能少，不能多
	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1: %v, len: %d, cap: %d", a1, len(a1), cap(a1))

	// 少了的部分就是默认零值， 等价于9， 8， 0
	a2 := [3]int{9, 8}
	fmt.Printf("a2: %v, len: %d, cap: %d", a2, len(a2), cap(a2))

	// 虽然没有显式初始化，但是实际上内存已经分配好，等价于0，0，0
	var a3 [3]int
	fmt.Printf("a3: %v, len: %d, cap: %d", a3, len(a3), cap(a3))

	// 数组不支持append操作
	//a3 = append(a3, 1)

	// 按下标索引，如果编译器能判断出来下标越界，那么会有编译错误，
	// 如果不能，那么运行时会报错，出现panic
	fmt.Printf("a1[1]:%d", a1[1])
}
