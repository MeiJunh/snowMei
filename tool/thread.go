package tool

import (
	"fmt"
)

// 以for range需要在协程使用的话参照
// test，需要先重新赋给一个新变量在进行使用
func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer PanicFunc(
		func(panic string) {
			fmt.Printf("panic,info:%s\n", panic)
		})

	fn()
}
