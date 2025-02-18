package mitmtools

import (
	"strconv"

	"github.com/Leviathangk/go-mitmtools-v2/handler"
)

const (
	defaultPort              = 8866
	defaultStreamLargeBodies = 1024 * 1024 * 5
	defaultSslInsecure       = true
	defaultShowLog           = true
)

type Config struct {
	Debug             int
	Addr              string
	Port              int   // 记录端口，方便获取
	StreamLargeBodies int64 // 当请求或响应体大于此字节时，转为 stream 模式
	SslInsecure       bool
	CaRootPath        string
	Upstream          string
	ShowLog           bool // 是否打印日志
	Backend           bool // 是否后台运行
	handlers          []handler.Addon
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

	return config
}
