package main

import "math"

func main() {
	var a int = 456
	var b int = 123
	println(a + b)
	println(a - b)
	println(a * b)
	println(a / b)
	String()
	//Byte()
	//Bool()
	//Extremum()
}

// Extremum 极值
func Extremum() {
	println("float64 最大值", math.MaxFloat64)
	// 没有float64 最小值
	println("float64 最小的正数", math.SmallestNonzeroFloat64)

	println("float32 最大值", math.MaxFloat32)
	// 没有float32 最小值
	println("float32 最小的正数", math.SmallestNonzeroFloat32)
}
