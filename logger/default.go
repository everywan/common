package logger

import "log"

// 日志作为基本组件, 每个项目都需要. 通过 init 初始化好
func init() {
	var err error
	std, err = NewZapLogger()
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
}
