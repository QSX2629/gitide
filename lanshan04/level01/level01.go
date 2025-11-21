package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func TimewithoutBuf(file *os.File, data []byte, count int) time.Duration {
	start := time.Now()
	for i := 0; i < count; i++ {
		_, err := file.Write(data)
		if err != nil {
			panic(err)
		}
	}
	return time.Since(start)
}
func TimewithBuf(file *os.File, data []byte, count int) time.Duration {
	start := time.Now()
	writer := bufio.NewWriter(file)
	for i := 0; i < count; i++ {
		_, err := writer.Write(data)
		if err != nil {
			panic(err)
		}
	}
	_ = writer.Flush()
	return time.Since(start)
}
func main() {
	testdata := []byte("对比时间")
	count := 10000000
	withoutBuf, err := os.Create("without_buf.txt")
	if err != nil {
		panic(err)
	}
	defer withoutBuf.Close()
	withoutBuf01 := TimewithoutBuf(withoutBuf, testdata, count)
	withBuf, err := os.Create("with_buf.txt")
	if err != nil {
		panic(err)
	}
	defer withBuf.Close()
	withBuf01 := TimewithBuf(withBuf, testdata, count)
	fmt.Println("Without Buf =", withoutBuf01)
	fmt.Println("With Buf =", withBuf01)
	fmt.Printf("两者的对比：%v", float64(withoutBuf01)/float64(withBuf01))
}
