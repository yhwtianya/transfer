package sender

import (
	"github.com/open-falcon/transfer/g"
	nlist "github.com/toolkits/container/list"
)

// 初始缓存队列
func initSendQueues() {
	cfg := g.Config()
	// judge缓存队列
	for node, _ := range cfg.Judge.Cluster {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}

	// graph缓存队列
	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
			GraphQueues[node+addr] = Q
		}
	}

	// tsdb缓存队列
	if cfg.Tsdb.Enabled {
		TsdbQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}
