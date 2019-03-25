package sender

import (
	"sort"

	"github.com/open-falcon/transfer/g"
	"github.com/toolkits/consistent"
)

func initNodeRings() {
	cfg := g.Config()

	// Judge一致性哈希环
	JudgeNodeRing = newConsistentHashNodesRing(cfg.Judge.Replicas, KeysOfMap(cfg.Judge.Cluster))
	// Graph一致性哈希环
	GraphNodeRing = newConsistentHashNodesRing(cfg.Graph.Replicas, KeysOfMap(cfg.Graph.Cluster))
}

// 返回排序的Key列表
// TODO 考虑放到公共组件库,或utils库
func KeysOfMap(m map[string]string) []string {
	// sort.StringSlice排序的[]string
	keys := make(sort.StringSlice, len(m))
	i := 0
	for key, _ := range m {
		keys[i] = key
		i++
	}

	keys.Sort()
	return []string(keys)
}

// 一致性哈希环,用于管理服务器节点.
type ConsistentHashNodeRing struct {
	ring *consistent.Consistent
}

// 创建一致性哈希环
func newConsistentHashNodesRing(numberOfReplicas int, nodes []string) *ConsistentHashNodeRing {
	ret := &ConsistentHashNodeRing{ring: consistent.New()}
	ret.SetNumberOfReplicas(numberOfReplicas)
	ret.SetNodes(nodes)
	return ret
}

// 根据pk,获取node节点. chash(pk) -> node
func (this *ConsistentHashNodeRing) GetNode(pk string) (string, error) {
	return this.ring.Get(pk)
}

// 添加节点
func (this *ConsistentHashNodeRing) SetNodes(nodes []string) {
	for _, node := range nodes {
		this.ring.Add(node)
	}
}

// 设置副本数量
func (this *ConsistentHashNodeRing) SetNumberOfReplicas(num int) {
	this.ring.NumberOfReplicas = num
}
