package main

import (
	"bytes"
	"fmt"
	"github.com/Leviathangk/go-glog/glog"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func containsStr(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func killPort(port int) {
	// 执行 netstat 命令并获取输出
	cmd := exec.Command("netstat", "-ano")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing netstat: %v", err)
		return
	}

	// 查找所有匹配的PID
	lines := strings.Split(out.String(), "\n")
	var pidArr []string
	portStr := strconv.Itoa(port)
	for _, line := range lines {
		if strings.Contains(line, portStr) {
			cols := strings.Fields(line)
			if len(cols) != 5 {
				continue
			}
			if !strings.Contains(cols[2], portStr) {
				continue
			}
			if strings.Contains(cols[3], "ESTABLISHED") || strings.Contains(cols[3], "LISTENING") {
				if containsStr(pidArr, cols[4]) {
					continue
				}
				fmt.Println("找到 pid:", cols[4])
				pidArr = append(pidArr, cols[4])
			}
		}
	}

	// kill 所有 pid
	for _, pid := range pidArr {
		killCommand := exec.Command("taskkill", "/PID", pid, "/F")
		if err := killCommand.Run(); err != nil {
			glog.DLogger.Warnf("%s kill 失败\n", pid)
		} else {
			glog.DLogger.Infof("%s kill 成功\n", pid)
		}
	}
}

func TestPort(t *testing.T) {
	//port := 8866
	killPort(8866)
}
