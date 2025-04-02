/*
下面的 handler 可以按需测试
命令：go test -v demo_test.go
*/
package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestDemo(t *testing.T) {
	sign1 := "d"
	sign2 := "d"
	toContent := fmt.Sprintf("console.log(`%s：$${%s} %s：$${%s}`),console.log(`%s：$${%s} %s：$${%s}`)", sign1, sign1, sign2, sign2, sign1, sign1, sign2, sign2)
	p := regexp.MustCompile("function")
	r := p.ReplaceAll([]byte("function工"), []byte(toContent))
	fmt.Println(string(r))
}
