# common
go 业务常用工具库.

[![go report card](https://goreportcard.com/badge/github.com/zhihu/norm "go report card")](https://goreportcard.com/report/github.com/zhihu/norm)
[![Go](https://github.com/zhihu/norm/actions/workflows/go.yml/badge.svg)](https://github.com/zhihu/norm/actions/workflows/go.yml)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/zhihu/norm)

## Overview

1. application: 管理程序和任务的生命周期. 任务包括 http、rpc、cron 等, 生命周期是指启动、终止、销毁, 支持 Hooks, 以及常用的工具注入如 sentry 等.
2. bundle: 任务.
   1. rest: http 服务
   2. rpc: grpc 服务
   3. cron: 定时任务


## Getting Started

Install:

```
go get github.com/everywan/common
```

use example: please go [use example](/examples/main.go)

## License

© Everywan, 2024~time.Now

Released under the [MIT License](/LICENSE)