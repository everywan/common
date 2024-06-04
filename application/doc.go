/*
app 用于管理程序的生命周期, 目前主要是启动和销毁. 用户编写 run() 和 stop() 方法,
app 会在程序启动时调用 run, 并在收到终止信号或正常结束后, 调用 stop 方法.

另外, 对于常见的 rpc/http 服务, 定时任务和消费任务等, common 提供相应的 bundle,
用户只需输入配置便可启动服务, 不必编写初始化代码.
*/

package application
