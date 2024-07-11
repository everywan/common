package feature

import (
	"math/rand"

	"github.com/everywan/common/configs"
)

// 开关控制, 配置 1、true 时表示打开, 否则关闭
func CheckSwitch(key string) bool {
	toggle := configs.Get(key)
	return toggle == "1" || toggle == "true"
}

// 放量开关. 配置 30 时, 表示放量 30%
func CheckRateLimit(key string, defaultPercent int) bool {
	percent := configs.GetIntOrDefault(key, defaultPercent)
	return percent > rand.Intn(100)
}
