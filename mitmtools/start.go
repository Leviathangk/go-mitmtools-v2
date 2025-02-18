package mitmtools

import (
	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools/handler"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"time"
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
	if opts.Backend {
		glog.DLogger.Debugln("程序正在后台运行！")
		runStatus, runErr := waitStart(p, opts.Port) // 这里是非阻塞式运行
		if runStatus {
			return p, nil
		} else {
			return nil, runErr
		}
	} else {
		err = p.Start() // 这里是阻塞式运行
		if err != nil {
			return nil, err
		} else {
			return p, nil
		}
	}
}

// waitStart 等待启动完成
func waitStart(p *proxy.Proxy, port int) (bool, error) {
	// 这里启动等待错误
	startCh := make(chan error, 1)
	go func(c *proxy.Proxy) {
		startErr := p.Start()
		startCh <- startErr
	}(p)

	// 这里等待错误和端口占用
	for {
		select {
		case errMsg := <-startCh:
			return false, errMsg
		case <-time.After(1 * time.Second):
			glog.DLogger.Debugf("正在等待端口 %d\n", port)
			s := PortIsAvailable(port)
			glog.DLogger.Debugln(s)
			if !PortIsAvailable(port) {
				return true, nil
			}
		}
	}
}
