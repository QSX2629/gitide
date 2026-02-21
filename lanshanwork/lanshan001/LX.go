package main

import "fmt"

func average(sum int, count int) float64 {
	var aver float64 = float64(sum) / float64(count)
	return aver
}
func main() {
	count := 0
	sum := 0

	var n int
	for {
		fmt.Println("请输入函数(0时停止)")
		fmt.Scanf("%d", &n)
		if n == 0 {
			break
		} else {
			sum = sum + n

			count++
		}
	}
	if count == 0 {
		fmt.Println("非法输入")
	}
	aver := average(sum, count)
	if aver < 60 {
		fmt.Println("不及格", aver)
	} else {
		fmt.Printf("及格", aver)
	}
}
