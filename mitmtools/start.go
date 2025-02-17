package mitmtools

import (
	"errors"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools/handler"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

// Start 启动入口
func Start(opts *Config, handlers ...handler.Addon) (*proxy.Proxy, error) {
	p, err := proxy.NewProxy(&proxy.Options{
		Debug:             opts.Debug,
		Addr:              opts.Addr,
		StreamLargeBodies: opts.StreamLargeBodies,
		SslInsecure:       opts.SslInsecure,
		CaRootPath:        opts.CaRootPath,
		Upstream:          opts.Upstream,
	})
	if err != nil {
		return nil, err
	}

	// 修改配置
	handler.ShowLog = opts.ShowLog
	glog.DLogger.ShowCaller = false

	// 添加解析响应体
	p.AddAddon(new(handler.DecodeRule))

	// 添加规则
	for _, h := range handlers {
		err = h.Check()
		if err != nil {
			return nil, err
		}

		p.AddAddon(h)
	}

	// 添加规则
	for _, h := range opts.handlers {
		err = h.Check()
		if err != nil {
			return nil, err
		}

		p.AddAddon(h)
	}

	// 添加响应体重新计算
	p.AddAddon(new(handler.RecalculateRule))

	// 执行
	startCh := make(chan string, 1)
	go func(c *proxy.Proxy) {
		startErr := p.Start()
		if startError!=nil{
			startCh <- startError.Error()
		}else{
			startCh<-"success"
		}
	}(p)

	select {
	case startMsg:=<-startCh:
		if startMsg=="success"{
			return p,nil
		}else{
			return nil, errors.New(startMsg)
		}
	}
}
