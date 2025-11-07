package main

import "fmt"

func Cal(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			fmt.Println("非法运算")
			return 0
		} else {
			return a / b
		}
	default:
		fmt.Println("非法运算符")
		return 0
	}
}
func main() {
	var a, b int
	var op string
	fmt.Println("请输入两个数和运算符")
	_, err := fmt.Scanf("%d %s %d", &a, &op, &b)
	if err != nil {
		fmt.Println(err)
	}
	sum := Cal(a, b, op)
	fmt.Printf("%d %s %d = %d\n", a, op, b, sum)

}
