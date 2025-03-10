package atomicx

import "sync/atomic"

func Atomic() {
	var val int32 = 12
	// 原子读，你不会读到修改了一般的数据
	val = atomic.LoadInt32(&val)
	println(val)
	// 原子写，即便别的goroutine 在别的CPU核上，也能立刻看到
	atomic.StoreInt32(&val, 13)
	// 原子自增，返回的是自增后的结果
	newVal := atomic.AddInt32(&val, 1)
	println(newVal)
	// CAS操作
	// 如果 val 的值是13，就修改为 15
	swapped := atomic.CompareAndSwapInt32(&val, 13, 15)
	println(swapped)
}
