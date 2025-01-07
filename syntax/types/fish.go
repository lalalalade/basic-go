package types

type Fish struct {
}

func (f Fish) Swim() {
	println("鱼在游")
}

// Yu 鱼
// 🐠是Fish的别名
type Yu = Fish

type FakeFish struct {
}

func (f FakeFish) FakeSwim() {
	println("假的鱼在游")
}

func UserFish() {
	f1 := Fish{}
	f1.Swim()
	f2 := FakeFish{}
	// f2将不能调用Fish上的方法，
	// 因为f2是一个全新的类型
	f2.FakeSwim()

	// 类型转换
	f3 := Fish(f2)
	f3.Swim()

	y := Yu{}
	y.Swim()
}
