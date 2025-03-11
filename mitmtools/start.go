package mitmtools

import (
	"bytes"
	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools-v2/handler"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"os/exec"
	"strconv"
	"strings"
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
		opts.Handler = NewHandler()
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
	return m.Config.Handler.AddHandler(h)
}

// RemoveHandler 移除配置
func (m *MitmWorker) RemoveHandler(handlerIndex int) {
	m.Config.Handler.RemoveHandler(handlerIndex)
}

// Start 启动
func (m *MitmWorker) Start() error {
	// kill 所有关联 job
	m.KillAll()

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

// KillAll 会杀掉所有与当前端口关联的进程，因为 ws 的连接会一直保持，stop 方法只会改变新连接
func (m *MitmWorker) KillAll() {
	killPort(m.Config.Port)
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

// containsStr 是否包含指定字符串
func containsStr(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// killPort 所有和指定端口关联的程序
func killPort(port int) {
	// 执行 netstat 命令并获取输出
	cmd := exec.Command("netstat", "-ano")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		glog.DLogger.Warnf("端口获取失败: %v\n", err)
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
				glog.DLogger.Debugln("找到 pid:", cols[4])
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
