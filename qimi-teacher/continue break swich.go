package main

import "fmt"

/*
	func main() {
		age := 10
		switch {
		case age > 10:
			fmt.Println("你可以上了")
		case age < 10:
			fmt.Println("你依然可以上")
		default:
			fmt.Println("are you pig?")
		}
	}
*/
/*func main() {
	i := 5main
	for i < 15 {
		i++
		if i == 10 {
			continue
		}
		fmt.Println(i)
	}

}*/
func main() {
	for i := 10; i < 20; i++ {
		if i == 15 {
			break
		}
		fmt.Println(i)
	}
}
