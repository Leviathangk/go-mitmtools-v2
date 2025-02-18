package main

import (
	"fmt"
	"github.com/Leviathangk/go-mitmtools-v2/handler"
	"testing"
)

type Handler struct {
	HandlerIndex int                   // 计数
	Handlers     map[int]handler.Addon // 方便添加和移除
}

func TestIp(t *testing.T) {
	h := &Handler{}
	fmt.Println(h.HandlerIndex)
	fmt.Println(h.Handlers == nil)
}
