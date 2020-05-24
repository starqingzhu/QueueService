package server

import (
	"QueueService/define"
	"QueueService/proto"
	"QueueService/queue"
)

func init() {
	RegisterParseHandle(define.CMD_QUERY_PLAYER_LOGIN_QUE_POS_REQ_NO, ParseAndHandleQueryLoginQuePosReq)
}

func ParseAndHandleQueryLoginQuePosReq(frame []byte, connInfo *define.ConnInfo) {
	queryReqInfo := proto.ParseToQueryPlayerLoginQuePosReq(frame)
	clientInfo := &define.ClientInfo{
		UserName: queryReqInfo.UserName,
		ConnAddr: (*connInfo.Conn).RemoteAddr().String(),
	}
	queue.QueryqueueChan <- *clientInfo
}
