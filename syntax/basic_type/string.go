package main

import "unicode/utf8"

func String() {
	// He said:"hello, go!"
	println("He said:\"hello, go!\"")
	println(`我可以换行
这是新的行
但是这里不能有反引号`)
	// 字符串拼接
	println("hello, " + "go")

	println(len("你好"))                      // 输出6
	println(utf8.RuneCountInString("你好"))   // 输出2
	println(utf8.RuneCountInString("你好ab")) // 输出4

}
