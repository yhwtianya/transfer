package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/open-falcon/transfer/sender"
)

// 获取连接池描述信息
func configDebugHttpRoutes() {
	// conn pools
	http.HandleFunc("/debug/connpool/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/debug/connpool/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		if argsLen < 1 {
			w.Write([]byte(fmt.Sprintf("bad args\n")))
			return
		}

		var result string
		receiver := args[0]
		switch receiver {
		case "judge":
			// 获取Judge连接池描述信息
			result = strings.Join(sender.JudgeConnPools.Proc(), "\n")
		case "graph":
			// 获取graph连接池描述信息
			result = strings.Join(sender.GraphConnPools.Proc(), "\n")
		default:
			result = fmt.Sprintf("bad args, module not exist\n")
		}
		w.Write([]byte(result))
	})
}
