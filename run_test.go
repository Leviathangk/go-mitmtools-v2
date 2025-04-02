package main

import (
	"fmt"
	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools-v2/handler/req"
	"github.com/Leviathangk/go-mitmtools-v2/handler/resp"
	"github.com/Leviathangk/go-mitmtools-v2/mitmtools"
	"log"
	"regexp"
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
	//worker.AddHandler(&resp.ReplaceContent{
	//	BaseHandler: handler.BaseHandler{},
	//	Pattern:     "baidu",
	//	FindContent: "百度",
	//	ToContent:   "千度",
	//	Times:       0,
	//	IsNoCookie:  false, // 有没有 cookie 都启用
	//	IsRegexp:    true,  // FindContent 使用正则匹配
	//})

	// 正则匹配查找进行替换
	//worker.AddHandler(&resp.ReplaceContent{
	//	BaseHandler: handler.BaseHandler{},
	//	Pattern:     "chl_api",
	//	FindContent: "l\\[\\w+\\(\\d+\\)\\]\\(null,n\\):k\\[l\\]\\[\\w+\\(\\d+\\)\\]\\(k,n\\)",
	//	//ToContent:   "(console.log(\"直接执行 \", l,n),l[\"apply\"](null, n)) : (console.log(\"正在执行函数 \",k,l,n),k[l][\"apply\"](k, n))",
	//	ToContent:  "(xxx=l[\"apply\"](null, n),console.log(\"直接执行 \", l,n,\" 结果 \",xxx),xxx) : (xxx=k[l][\"apply\"](k, n),console.log(\"正在执行函数 \",k,l,n,\" 结果 \",xxx),xxx)",
	//	Times:      0,
	//	IsNoCookie: false, // 有没有 cookie 都启用
	//	IsRegexp:   true,  // FindContent 使用正则匹配
	//})

	// 支持自定义如何查找，如何替换
	worker.AddHandler(&resp.ReplaceContentCustomize{
		Pattern:    "chl_api",
		Times:      0,
		IsNoCookie: false, // 有没有 cookie 都启用
		ReplaceFunc: func(body []byte) []byte {
			p := regexp.MustCompile("function \\w+\\((\\w+),\\w,\\w,(\\w),\\w,\\w,\\w,\\w,\\w\\)")
			f := p.FindAllSubmatch(body, -1)
			if len(f) == 1 && len(f[0]) == 3 {
				return body
			}

			sign1 := string(f[0][1])
			sign2 := string(f[0][2])

			glog.DLogger.Infof("成功匹配 sign1：%s sign2：%s\n", sign1, sign2)

			findContent := regexp.MustCompile("l\\[\\w+\\(\\d+\\)\\]\\(null,n\\):k\\[l\\]\\[\\w+\\(\\d+\\)\\]\\(k,n\\)")
			toContent := fmt.Sprintf("(xxx=l[\"apply\"](null, n),console.log(`%s：${%s} %s：${%s}`),console.log(\"直接执行 \", l,n,\" 结果 \",xxx),xxx) : (xxx=k[l][\"apply\"](k, n),console.log(`%s：${%s} %s：${%s}`),console.log(\"正在执行函数 \",k,l,n,\" 结果 \",xxx),xxx)", sign1, sign1, sign2, sign2, sign1, sign1, sign2, sign2)
			return findContent.ReplaceAll(body, []byte(toContent))
		},
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
