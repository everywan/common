package utils

import "os"

var (
	isDebug   = false
	app       = os.Getenv("APP_NAME")
	env       = os.Getenv("APP_ENV")
	sandboxID = os.Getenv("SANDBOX_ID")
)

const (
	EnvProd    = "production" // 生产环境
	EnvStaging = "staging"    // 预发环境
	EnvTesting = "testing"    // 测试环境
	EnvCI      = "test_ci"    // CI 环境
	EnvDevelop = "develop"    // 开发环境

	// 沙盒环境
	NonSandboxDeployID = "0" // 非联调环境部署Id
)

func App() string {
	return app
}

func Env() string {
	return env
}

func IsProdEnv() bool {
	return Env() == EnvProd
}

func IsStagingEnv() bool {
	return Env() == EnvStaging
}

func IsTestingEnv() bool {
	return Env() == EnvTesting
}

func IsCIEnv() bool {
	return Env() == EnvCI
}

func IsDevelopEnv() bool {
	return Env() == EnvDevelop
}

func IsSandBoxEnv() bool {
	return sandboxID != "" && sandboxID != NonSandboxDeployID
}

func SandBoxID() string {
	return sandboxID
}

func IsDebug() bool {
	return isDebug
}

func EnableDebug() {
	isDebug = true
}
