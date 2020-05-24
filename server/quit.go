package server

import (
	"QueueService/define"
	"QueueService/proto"
	"QueueService/queue"
)

func init() {
	RegisterParseHandle(define.CMD_LOGIN_QUIT_REQ_NO, ParseAndHandleLoginQuitReq)
}

func ParseAndHandleLoginQuitReq(frame []byte, connInfo *define.ConnInfo) {
	quitReqInfo := proto.ParseToQuitLoginQueReq(frame)
	clientInfo := &define.ClientInfo{
		UserName: quitReqInfo.UserName,
		ConnAddr: (*connInfo.Conn).RemoteAddr().String(),
	}
	queue.QuitQueueChan <- *clientInfo
}
