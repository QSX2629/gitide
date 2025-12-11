package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

func main() {
	//map的初始化*******************
	a := make(map[string]int, 8)
	fmt.Println(a == nil)
	//添加键值对*********************
	a["哈哈哈哈"] = 100
	a["啦啦啦"] = 200
	fmt.Println(a)
	fmt.Printf("type:%T\n", a)
	//判断是否有键值对************
	b := make(map[string]int)
	b["嗡嗡嗡"] = 100
	b["呃呃呃"] = 200
	value, ok := b["aaa"]
	if ok {
		fmt.Println("在", value)
	} else {
		fmt.Println("不在")
	}
	//遍历键值对*************for range
	c := make(map[int]int, 6)
	c[1] = 100
	c[2] = 200
	for k, v := range c {
		fmt.Println(k, v)
	}
	//用delete删除键值对**********************
	delete(b, "嗡嗡嗡")
	fmt.Println(b)
	//用某个顺序遍历键值对**************************
	u := make(map[string]int, 100)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("stu%02d", i)
		value := rand.Intn(100) //随机数
		u[key] = value
	}
	for key, value := range u {
		fmt.Println(key, value)
	}
	//将键值对放入切片中，把切片进行排序
	keys := make([]string, 0, 100)
	for key := range u {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key, u[key])
	}
	//元素类型为map的切片************
	ms := make([]map[string]int, 8, 8) //切片初始化
	ms[0] = make(map[string]int, 8)    //map初始化
	ms[0]["啦啦啦"] = 100
	fmt.Println(ms[0])
	//值为切片的map*********************************
	sm := make(map[string][]int, 8)
	sm["哈哈哈"] = make([]int, 8, 8)
	sm["哈哈哈"][0] = 100
	sm["哈哈哈"][1] = 200
	fmt.Println(sm)
	for k, v := range sm {
		fmt.Println(k, v)
	}
	//练习**********************************
	var s = "How do you do"
	words := strings.Split(s, " ")
	z := make(map[string]int, 8)
	for _, word := range words {
		v, ok := z[word]
		if ok {
			z[word] = v + 1
		} else {
			z[word] = 1
		}

	}
	for k, v := range z {
		fmt.Println(k, v)
	}

}
