/*
下面的 handler 可以按需测试
命令：go test -v demo_test.go
*/
package main

import (
	"log"
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
	worker, err := mitmtools.NewWorker(mitmtools.NewConfig(
		mitmtools.SetPort(Port),
		mitmtools.SetSslInsecure(true),
		mitmtools.SetProxy(ProxyUrl),
		mitmtools.SetShowLog(true),
		//mitmtools.SetBackend(true), // 后台运行
		//mitmtools.SetCaRootPath("C:\\Users\\用户目录\\.mitmproxy"),	// windows 示例
	))
	if err != nil {
		log.Fatalln(err)
	}

	worker.AddHandler(&req.ShowReq{})

	err = worker.Start()
	if err != nil {
		log.Fatalln(err)
	}
	glog.DLogger.Debugln("程序已启动...")

	glog.DLogger.Debugln(worker.Proxy.GetCertificate())

	for {
		time.Sleep(5 * time.Second)
		glog.DLogger.Debugln("程序正在关闭...")
		worker.Stop()
		break
	}
}
