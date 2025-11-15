package main

import "fmt"

func level1(array []int) map[int]int {
	MAP := make(map[int]int)
	for _, v := range array {
		MAP[v]++
	}
	return MAP
}
func main() {
	testarray := []int{1, 2, 3, 4, 5, 1, 5, 3, 4}
	fmt.Println(level1(testarray))

}
