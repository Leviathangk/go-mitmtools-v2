package mitmtools

import (
	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools-v2/handler"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"strconv"
	"time"
)

type MitmWorker struct {
	Proxy  *proxy.Proxy
	Config *Config
}

// NewWorker 创建一个处理器
func NewWorker(opts *Config) (*MitmWorker, error) {
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
	if opts.Handler == nil {
		opts.Handler = new(Handler)
		opts.Handler.HandlerIndex = 0
		opts.Handler.Handlers = make(map[int]handler.Addon)
	}
	if opts.Addr == "" {
		opts.Port = defaultPort
		opts.Addr = ":" + strconv.Itoa(defaultPort)
	}
	if opts.StreamLargeBodies == 0 {
		opts.StreamLargeBodies = defaultStreamLargeBodies
	}

	return &MitmWorker{
		Config: opts,
		Proxy:  p,
	}, nil
}

// AddHandler 添加配置
func (m *MitmWorker) AddHandler(h handler.Addon) int {
	index := m.Config.Handler.HandlerIndex + 1
	m.Config.Handler.Handlers[index] = h
	return index
}

// RemoveHandler 移除配置
func (m *MitmWorker) RemoveHandler(handlerIndex int) {
	if _, ok := m.Config.Handler.Handlers[handlerIndex]; ok {
		delete(m.Config.Handler.Handlers, handlerIndex)
	}
}

// Start 启动
func (m *MitmWorker) Start() error {
	// 清空原始 handlers
	m.Proxy.Addons = make([]proxy.Addon, 0)

	// 添加解析响应体
	m.Proxy.AddAddon(new(handler.DecodeRule))

	// 添加 handlers
	for _, h := range m.Config.Handler.Handlers {
		err := h.Check()
		if err != nil {
			return err
		}
		m.Proxy.AddAddon(h)
	}

	// 添加响应体重新计算
	m.Proxy.AddAddon(new(handler.RecalculateRule))

	glog.DLogger.Debugf("启动地址 %s\n", m.Config.Addr)

	// 启动
	if m.Config.Backend {
		glog.DLogger.Debugln("正在后台运行...")
		runStatus, runErr := waitStart(m.Proxy, m.Config.Port) // 这里是非阻塞式运行
		if runStatus {
			return nil
		} else {
			return runErr
		}
	} else {
		err := m.Proxy.Start() // 这里是阻塞式运行
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

// Stop 停止
func (m *MitmWorker) Stop() error {
	err := m.Proxy.Close()
	if err != nil {
		return err
	}
	glog.DLogger.Debugln("关闭完成...")

	return nil
}

// ReStart 重启
func (m *MitmWorker) ReStart() error {
	var err error

	if !PortIsAvailable(m.Config.Port) {
		err = m.Stop()
		if err != nil {
			return err
		}
	}

	err = m.Start()
	if err != nil {
		return err
	}
	glog.DLogger.Debugln("重启完成...")

	return nil
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
			glog.DLogger.Debugf("正在等待端口启动 %d\n", port)
			if !PortIsAvailable(port) {
				glog.DLogger.Debugf("端口已启动 %d\n", port)
				return true, nil
			}
		}
	}
}
