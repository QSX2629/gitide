package main

import (
	"fmt"
	"os"
)

func main() {
	/*fmt.Println("请输入年龄和名字")
	var name string
	var age int
	_, err := fmt.Scan(&name, &age)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("yse")
	}
	file, err := os.Open("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	var name01 string
	for {
		_, err = fmt.Fscanf(file, "%s", &name01)
		if err != nil {
			break
		}
		fmt.Println(name01)
	}*/
	/*filepath := "io_test.txt"
		WriteFile, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer WriteFile.Close()
		WriteContent := "这是通过os.file(io.writer)写入的内容"
		WriteN, Writeerr := io.WriteString(WriteFile, WriteContent)
		if Writeerr != nil {
			fmt.Println(Writeerr)
		} else {
			fmt.Printf("成功写入%d个字节\n", WriteN)
		}
		readFile, err := os.Open(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer readFile.Close()
		ReadContent, err := io.ReadAll(readFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(ReadContent))
		readFile.Seek(0, 0)
		buf := make([]byte, 1024)
		n, err := readFile.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Read bytes:", n, "Contene:", string(buf[:n]))
	}*/
	var args []string = os.Args
	if len(args) < 2 {
		fmt.Printf("没有输入名字")
		return
	}
	name := args[1]
	fmt.Printf("你好%s\n", name)
}
