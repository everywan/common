package main

import (
	"fmt"

	"github.com/everywan/common/configs"
)

func main() {
	// nacos 配置中心
	// 从默认分组获取数据
	fmt.Printf("get config type_string:%s\n", configs.Get("test_string"))
	fmt.Printf("get config type_string:%d\n", configs.GetInt("test_int"))
	out := struct {
		Addr string `json:"addr"`
	}{}
	configs.GetJson("test_json", &out)
	fmt.Printf("get config type_string:%+v\n", out)

	// 从指定分组获取数据
	fmt.Printf("get config from group:test, key:type_string:%s\n", configs.GetByGroup("test", "test_string"))

	// 监听配置
	configs.MonitorChange("", "test_monitor", func(newvalue string) {
		fmt.Printf("key test_monitor value change to %s\n", newvalue)
	})

	fmt.Println("监听配置变更中. 按任意键结束")
	a := ""
	fmt.Scan(&a)
}
