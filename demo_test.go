/*
下面的 handler 可以按需测试
命令：go test -v demo_test.go
*/
package main

import (
	"testing"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools/handler/req"
	"github.com/Leviathangk/go-mitmtools/mitmtools"
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
		mitmtools.SetBackend(true),
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
}
