package main

import (
	"fmt"
	"sort"
)

func main() {
	a := [5]int{1, 2, 3, 4, 5}
	b := a[1:3]
	fmt.Println(b)
	fmt.Println(a)
	c := b[1:2]
	fmt.Println(c)
	//用make函数构造切片********************
	d := make([]int, 5, 10)
	fmt.Printf("切片d是:%v\n", d)
	//用len函数获取切片长度
	fmt.Println(len(d))
	//用cap函数获取切片容量
	fmt.Println(cap(d))
	//切片的索引*************
	e := []string{"重庆", "深圳", "杭州", "上海"}
	for i := 0; i < len(e); i++ {
		fmt.Println(i+1, e[i])
	}
	for _, value := range e {
		fmt.Println(value)
	}
	//切片的扩容**************************
	f := []int{}
	for i := 0; i < 5; i++ {
		f = append(f, i)
		fmt.Printf("%v len:%d cap:%d ptr:%p\n", f, len(f), cap(f), f) //特别注意printf！！！
	}
	//切片的copy函数*************************
	n := []int{1, 2, 3, 4, 5}
	m := make([]int, 5, 5)
	copy(m, n)
	m[3] = 10
	l := m
	fmt.Println(m)
	fmt.Println(n)
	fmt.Println(l)
	//切片的删减*************************舍弃第二个k的前括号的数
	k := []int{1, 2, 3, 4, 5}
	k = append(k[:2], k[3:]...)
	fmt.Println(k)
	//用sort包进行排序
	h := []int{5, 4, 2, 3}
	sort.Ints(h[:])
	fmt.Println(h)
}
