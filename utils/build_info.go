package utils

import (
	"encoding/base64"
	"log"
	"strings"
)

var (
	AppName      string // 应用名称
	AppVersion   string // 应用版本
	BuildVersion string // 编译版本
	BuildTime    string // 编译时间
	GitRevision  string // Git版本
	GitBranch    string // Git分支
	GoVersion    string // Golang信息
)

// todo demo 替换为实际名.
// 打印程序构建信息, 需要在打包时注入变量. 程序实例参考 demo.
/*
Need add build flags
	APP_NAME := commons
	APP_VERSION     := $(shell git describe --abbrev=0 --tags)
	BUILD_VERSION   := $(shell git log -1 --oneline | base64)
	BUILD_TIME      := $(shell date "+%FT%T%z")
	GIT_REVISION    := $(shell git rev-parse --short HEAD)
	GIT_BRANCH      := $(shell git name-rev --name-only HEAD)
	GO_VERSION      := $(shell go version)
	go build -ldflags "
		-X 'github.com/everywan/common/utils.AppName=${APP_NAME}'             \
		-X 'github.com/everywan/common/utils.AppVersion=${APP_VERSION}'       \
		-X 'github.com/everywan/common/utils.BuildVersion=${BUILD_VERSION}'   \
		-X 'github.com/everywan/common/utils.BuildTime=${BUILD_TIME}'         \
		-X 'github.com/everywan/common/utils.GitRevision=${GIT_REVISION}'     \
		-X 'github.com/everywan/common/utils.GitBranch=${GIT_BRANCH}'         \
		-X 'github.com/everywan/common/utils.GoVersion=${GO_VERSION}'         \
		-s -w
	" -mod vendor -v -o  $(NAME) ${MAIN}
*/
func PrintBuildInfo() {
	v, _ := base64.StdEncoding.DecodeString(BuildVersion)
	BuildVersion = strings.TrimSpace(string(v))

	log.Println("App Name:", AppName)
	log.Println("App Version:", AppVersion)
	log.Println("Build version:", BuildVersion)
	log.Println("Build time:", BuildTime)
	log.Println("Git revision:", GitRevision)
	log.Println("Git branch:", GitBranch)
	log.Println("Golang Version:", GoVersion)
}
