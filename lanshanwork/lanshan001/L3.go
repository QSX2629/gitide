package main

import "fmt"

func main() {
	var ji int = 1
	var n int
	fmt.Println("请输入一个整数")
	fmt.Scanf("%d", &n)
	for i := n; i > 0; i-- {
		ji = ji * i
	}
	fmt.Println(ji)
}
