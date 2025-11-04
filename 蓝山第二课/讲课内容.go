package main

import "fmt"

/*
	func main() {
		var x, y []int
		for i := 0; i < 10; i++ {
			y = append(x, i)
			fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
			x = y
		}
	}
*/
type Printer func(contents string) (n int, err error)

func printToStd(contents string) (bytesNum int, err error) {
	return fmt.Println(contents)
}

func main() {
	var p Printer
	p = printToStd
	p("something")
}
