package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	name := "邱双喜"
	n, err := fmt.Fprintf(file, "姓名：%v", name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("成功输入%d个字节", n)
}
