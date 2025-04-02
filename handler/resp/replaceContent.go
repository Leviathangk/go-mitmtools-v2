package resp

import (
	"fmt"
	"github.com/Leviathangk/go-glog/glog"
	"github.com/Leviathangk/go-mitmtools-v2/handler"
	"regexp"
	"strings"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

type ReplaceContent struct {
	handler.BaseHandler
	Pattern     string // url 匹配规则
	FindContent string // 查找的需要替换的内容：针对部分内容
	ToContent   string // 替换后的新内容
	Times       int    // 替换次数（0 代表无限）
	IsRegexp    bool   // 是否正则匹配
	IsNoCookie  bool   // 只有当没有 cookie 的时候才替换
	timesRecord int    // 记录当前次数
}

func (r *ReplaceContent) Response(f *proxy.Flow) {

	// 替换响应
	if handler.IsMatch(r.Pattern, f.Request.URL.String()) {
		if r.IsNoCookie && handler.CookieExists(f) {
			glog.DLogger.Warnf("当前存在 cookie 不进行替换：%s\n", f.Request.Header.Get("cookie"))
			return
		}

		if r.Times != 0 {
			if r.timesRecord >= r.Times {
				glog.DLogger.Warnf("当前替换已达到上限：%d\n", r.Times)
				return
			}
			r.timesRecord += 1
			glog.DLogger.Debugf("当前替换次数：%d-%d\n", r.Times, r.timesRecord)
		}

		if r.IsRegexp {
			p := regexp.MustCompile(r.FindContent)
			f.Response.Body = p.ReplaceAll(f.Response.Body, []byte(r.ToContent))
		} else {
			f.Response.Body = []byte(strings.ReplaceAll(string(f.Response.Body), r.FindContent, r.ToContent))
		}

		if handler.ShowLog || r.ShowLog {
			glog.DLogger.Debugf("ReplaceContent 已修改响应结果：%s\n", f.Request.URL)
		}
	}
}

// Check 检查是否符合启动要求
func (r *ReplaceContent) Check() error {

	if r.FindContent == "" {
		return fmt.Errorf("ReplaceContent 未接收到需要替换的内容！")
	}

	return nil
}

// CustomizeReplaceFunc 自定义的替换函数
type CustomizeReplaceFunc func(body []byte) []byte

// ReplaceContentCustomize 自定义如何处理响应体
type ReplaceContentCustomize struct {
	handler.BaseHandler
	Pattern     string               // url 匹配规则
	ReplaceFunc CustomizeReplaceFunc // 自定义替换函数
	Times       int                  // 替换次数（0 代表无限）
	IsNoCookie  bool                 // 只有当没有 cookie 的时候才替换
	timesRecord int                  // 记录当前次数
}

func (r *ReplaceContentCustomize) Response(f *proxy.Flow) {

	// 替换响应
	if handler.IsMatch(r.Pattern, f.Request.URL.String()) {
		if r.IsNoCookie && handler.CookieExists(f) {
			glog.DLogger.Warnf("当前存在 cookie 不进行替换：%s\n", f.Request.Header.Get("cookie"))
			return
		}

		if r.Times != 0 {
			if r.timesRecord >= r.Times {
				glog.DLogger.Warnf("当前替换已达到上限：%d\n", r.Times)
				return
			}
			r.timesRecord += 1
			glog.DLogger.Debugf("当前替换次数：%d-%d\n", r.Times, r.timesRecord)
		}

		f.Response.Body = r.ReplaceFunc(f.Response.Body)

		if handler.ShowLog || r.ShowLog {
			glog.DLogger.Debugf("ReplaceContent 已修改响应结果：%s\n", f.Request.URL)
		}
	}
}

// Check 检查是否符合启动要求
func (r *ReplaceContentCustomize) Check() error {
	if r.ReplaceFunc == nil {
		return fmt.Errorf("未定义修改 body 函数！")
	}
	return nil
}
