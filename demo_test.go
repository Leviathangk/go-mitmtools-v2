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
	content := []byte("ddl[j6(400)](null,n):k[l][j6(400)]")
	p := regexp.MustCompile("l\\[\\w+\\(\\d+\\)\\]\\(null,n\\):k\\[l\\]\\[\\w+\\(\\d+\\)\\]")
	t1 := p.ReplaceAll(content, []byte("2"))
	fmt.Println(t1)
	fmt.Println(string(t1))
}
