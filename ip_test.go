package main

import (
	"fmt"
	"testing"
)

func TestIp(t *testing.T) {
	var a []int
	a = append(a, 1)
	a = append(a, 1)
	fmt.Println(a)
	fmt.Println(len(a))
	a = make([]int, 0)
	fmt.Println(a)
	fmt.Println(len(a))
	a = append(a, 1)
	a = append(a, 1)
	fmt.Println(a)
	fmt.Println(len(a))
}
