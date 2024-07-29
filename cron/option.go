package cron

import (
	"strings"

	"github.com/everywan/common/utils"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)

type (
	options struct {
		serverAddr  string
		accessToken string
		executorIP  string
		registryKey string
		logger      xxl.Logger
		port        int
	}
	Option func(*options)
)

func (opts *options) LoadDefault() {
	if opts.serverAddr == "" {
		opts.serverAddr = "http://xxl-job-admin.devops.svc:8080/xxl-job-admin"
	}
	if opts.accessToken == "" {
		opts.accessToken = "default_token"
	}
	// if opts.executorIP == "" {
	// 	opts.executorIP = ipv4.LocalIP()
	// }
	if opts.port == 0 {
		opts.port = 9999
	}
	if opts.registryKey == "" {
		opts.registryKey = strings.Join([]string{utils.App(), utils.Env()}, ":")
	}
	if opts.logger == nil {
		opts.logger = &Logger{}
	}
}

func ServerAddr(addr string) Option {
	return func(o *options) {
		o.serverAddr = addr
	}
}

// AccessToken 请求令牌
func AccessToken(token string) Option {
	return func(o *options) {
		o.accessToken = token
	}
}

// ExecutorIp 设置执行器IP
func ExecutorIp(ip string) Option {
	return func(o *options) {
		o.executorIP = ip
	}
}

// ExecutorPort 设置执行器端口
func ExecutorPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

// RegistryKey 设置执行器标识
func RegistryKey(registryKey string) Option {
	return func(o *options) {
		o.registryKey = registryKey
	}
}
