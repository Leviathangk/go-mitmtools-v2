package mitmtools

import (
	"strconv"

	"github.com/Leviathangk/go-mitmtools-v2/handler"
)

const (
	defaultPort              = 8866
	defaultStreamLargeBodies = 1024 * 1024 * 5
)

type Handler struct {
	HandlerIndex int                   // 计数
	Handlers     map[int]handler.Addon // 方便添加和移除
}
type Config struct {
	Debug             int
	Addr              string
	Port              int   // 记录端口，方便获取
	StreamLargeBodies int64 // 当请求或响应体大于此字节时，转为 stream 模式
	SslInsecure       bool  // 为 true 时不验证上游服务器的 SSL/TLS 证书
	CaRootPath        string
	Upstream          string
	ShowLog           bool    // 是否打印日志
	Backend           bool    // 是否后台运行
	Handler           Handler // 处理器
}
type SetFunc func(c *Config)

// NewConfig 新建配置
func NewConfig(opt ...SetFunc) *Config {
	config := new(Config)

	for _, o := range opt {
		o(config)
	}

	// 参数检查
	if config.Addr == "" {
		config.Port = defaultPort
		config.Addr = ":" + strconv.Itoa(defaultPort)
	}
	if config.StreamLargeBodies == 0 {
		config.StreamLargeBodies = defaultStreamLargeBodies
	}

	config.Handler.HandlerIndex = 0
	config.Handler.Handlers = make(map[int]handler.Addon)

	return config
}
