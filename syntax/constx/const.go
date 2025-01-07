package constx

const internal = "包内可访问"
const External = "包外可访问"

func Const() {
	const a = "你好"
	print(a)
}

const (
	StatusA = iota
	StatusB
	StatusC
	StatusD
	StatusE
	StatusF
	StatusG = 6
	StatusH
)

const (
	One = iota << 1
	Two
	Three
	Four
)
