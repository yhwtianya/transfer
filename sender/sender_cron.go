package sender

import (
	"log"
	"strings"
	"time"

	"github.com/open-falcon/transfer/proc"
	"github.com/toolkits/container/list"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second    //ProcCron的周期,默认1s
	DefaultLogCronPeriod  = time.Duration(3600) * time.Second //LogCron的周期,默认300s
)

// 定期任务，比如统计缓存数据量
// send_cron程序入口
func startSenderCron() {
	go startProcCron()
	go startLogCron()
}

// 定期统计缓存数据数量
func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

// 定期输出Graph连接池描述信息
func startLogCron() {
	for {
		time.Sleep(DefaultLogCronPeriod)
		logConnPoolsProc()
	}
}

// 统计缓存数据数量
func refreshSendingCacheSize() {
	// 统计未发送到Judge总数
	proc.JudgeQueuesCnt.SetCnt(calcSendCacheSize(JudgeQueues))
	// 统计未发送到Graph总数
	proc.GraphQueuesCnt.SetCnt(calcSendCacheSize(GraphQueues))
}

// 计算集群每个节点队列数据数量之和
func calcSendCacheSize(mapList map[string]*list.SafeListLimited) int64 {
	var cnt int64 = 0
	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}

// 输出Graph连接池描述信息
func logConnPoolsProc() {
	log.Printf("connPools proc: \n%v", strings.Join(GraphConnPools.Proc(), "\n"))
}
