package main

import (
	"fmt"
	"strings"
)

func makesuffixfunc(suffix string) func(string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}
func hanshu3(base int) (func(int) int, func(int) int) {
	add := func(i int) int {
		base += i
		return base
	}
	sub := func(i int) int {
		base -= i
		return base
	}
	return add, sub
}

func main() {
	r := makesuffixfunc(".txt")
	ret := r("啦啦啦")
	fmt.Println(ret)
	x, y := hanshu3(33)
	ret2 := x(100)
	ret1 := y(200)
	fmt.Println(ret2)
	fmt.Println(ret1)
}
