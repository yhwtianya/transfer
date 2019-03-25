package receiver

import (
	"github.com/open-falcon/transfer/receiver/rpc"
	"github.com/open-falcon/transfer/receiver/socket"
)

func Start() {
	// 启动Rpc服务
	go rpc.StartRpc()
	// 处理原始socket的上报数据请求
	go socket.StartSocket()
}
