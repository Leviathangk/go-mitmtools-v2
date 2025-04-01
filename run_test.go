package main

import (
	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools-v2/handler"
	"github.com/Leviathangk/go-mitmtools-v2/handler/req"
	"github.com/Leviathangk/go-mitmtools-v2/handler/resp"
	"github.com/Leviathangk/go-mitmtools-v2/mitmtools"
	"log"
	"testing"
)

const (
	Port     = 8866
	ProxyUrl = ""
)

func TestRun(t *testing.T) {
	worker, err := mitmtools.NewWorker(mitmtools.NewConfig(
		mitmtools.SetPort(Port),
		mitmtools.SetSslInsecure(true),
		mitmtools.SetProxy(ProxyUrl),
		mitmtools.SetShowLog(true),
		mitmtools.SetBackend(false), // true 为后台运行
		//mitmtools.SetCaRootPath("C:\\Users\\用户目录\\.mitmproxy"),	// windows 示例
	))
	if err != nil {
		log.Fatalln(err)
	}

	worker.AddHandler(&req.ShowReq{})
	worker.AddHandler(&resp.ReplaceContent{
		BaseHandler: handler.BaseHandler{},
		Pattern:     "baidu",
		FindContent: "百度",
		ToContent:   "千度",
		Times:       0,
		IsNoCookie:  false, // 有没有 cookie 都启用
		IsRegexp:    true,  // FindContent 使用正则匹配
	})

	err = worker.Start()
	if err != nil {
		log.Fatalln(err)
	}
	glog.DLogger.Debugln("程序已启动...")

	//glog.DLogger.Debugln(worker.Proxy.GetCertificate())
	//
	//for {
	//	time.Sleep(5 * time.Second)
	//	glog.DLogger.Debugln("程序正在关闭...")
	//	worker.Stop()
	//	break
	//}
}
