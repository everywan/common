/*
更稳定, 更高效的 grpc sdk.

client:
1. 负载均衡优化: 默认生效, 替换为 round_robin 轮询机制.
2. 退避机制: 默认生效, 避免立即重试导致的失败.
3. todo: 链路追踪 traceid, header
4. todo 打点(qps、rt) -> 监控报警
4. todo: 断路器插件
5. todo: 限流插件
7. todo: ip 自定义解析


server
3. todo: 链路追踪 traceid, header
4. todo 打点(qps、rt) -> 监控报警
1. http 检活
2. 默认超时参数
3. get caller, 根据 caller 做限速、监控追踪

*/

package grpc
