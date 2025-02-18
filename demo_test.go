/*
下面的 handler 可以按需测试
命令：go test -v demo_test.go
*/
package main

import (
	"testing"
	"time"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools-v2/handler/req"
	"github.com/Leviathangk/go-mitmtools-v2/mitmtools"
)

const (
	Port     = 8866
	ProxyUrl = ""
)

func TestDemo(t *testing.T) {
	config := mitmtools.NewConfig(
		mitmtools.SetPort(Port),
		mitmtools.SetSslInsecure(true),
		mitmtools.SetProxy(ProxyUrl),
		mitmtools.SetShowLog(true),
		mitmtools.SetBackend(true), // 后台运行
		//mitmtools.SetCaRootPath("C:\\Users\\用户目录\\.mitmproxy"),	// windows 示例
	)

	// 打印请求
	config.AddHandler(&req.ShowReq{})

	proxy, err := mitmtools.Start(config)
	glog.DLogger.Debugln("程序已启动...")

	if err != nil {
		glog.DLogger.Fatalln(err)
	}

	glog.DLogger.Debugln(proxy.GetCertificate())

	for {
		time.Sleep(5 * time.Second)
		glog.DLogger.Debugln("程序正在关闭...")
		proxy.Close()
		break
	}
	glog.DLogger.Debugln(PortIsAvailable(8866))
}
