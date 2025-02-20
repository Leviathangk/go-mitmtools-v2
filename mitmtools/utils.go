package mitmtools

import (
	"fmt"
	"net"
)

// GetFreePort 获取随机可用端口
func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0") // 使用 ":0" 作为地址来让系统自动分配可用端口
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil // 获取分配到的端口号
}

// PortIsAvailable 判断端口是否被占用
func PortIsAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
