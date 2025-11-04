package main

import "fmt"

func a() {
	fmt.Println("hello")
}
func b() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	panic("panic in b")
}

func main() {
	a()
	b()

}
