package main

import "fmt"

/*
func intsum(a int, b int) int {
ret := a + b
return ret
*/
/*func intsum2(a int, b int) (ret int) {
	ret = a + b
	return
}*/
/*func intsum3(d int, c ...int) int { //d为固定参数，c为可变参数***********************
ret := 0
for _, v := range c {
	ret += v
}
return ret*/
/*func intsum5(a, b int) (sum, sub int) {                  //返回多个变量*************
sum = a + b
sub = a * b
return*/
/*func fff() {
m := 10
n := "啦啦啦"
fmt.Println(m)
fmt.Println(n)}*/
func add(x, y int) int { //函数作为变量***********************
	return x + y
}
func sub(x, y int) int {
	return x - y
}
func cal(x, y int, op func(int, int) int) int {
	return op(x, y)
}

func main() {
	/*r1 := intsum3(1, 2, 3)
	r3 := intsum3(1, 2)
	fmt.Println(r1, r3)
	x, y := intsum5(1, 2)
	defer fmt.Println(x)
	defer fmt.Println(y)
	fff()*/
	r1 := cal(3, 4, add)
	fmt.Println(r1)
	r2 := cal(3, 4, sub)
	fmt.Println(r2)
}
