package main

import (
	"QueueService/preload"
	"QueueService/queue"
	"QueueService/server"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func main() {
	if preload.Conf.Stress.Switch {
		go func() {
			http.ListenAndServe(preload.Conf.Stress.HttpAddr, nil)
		}()
	}

	runtime.GOMAXPROCS(preload.Conf.Server.GoMaxProcsNum)
	//初始化工作
	queue.Init()
	queue.OperateWaitList()

	addr := fmt.Sprintf("tcp://:%d", preload.Conf.TcpServer.TcpPort)
	server.CodecServerRun(addr, preload.Conf.Server.Multicore, true, nil)
}
