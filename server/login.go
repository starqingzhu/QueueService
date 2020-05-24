package server

import (
	"QueueService/define"
	"QueueService/proto"
	"QueueService/queue"
	"log"
)

func init() {
	RegisterParseHandle(define.CMD_LOGIN_REQ_NO, ParseAndHandleLoginReq)
}

func ParseAndHandleLoginReq(frame []byte, connInfo *define.ConnInfo) {
	loginReqInfo := proto.ParseToLoginReq(frame)

	var status uint16 = define.STATUS_LOGIN_PRE_SUCCESS
	inWaitQueFlag := queue.CheckPlayerIsInWaitQue(loginReqInfo.UserName)

	if inWaitQueFlag {
		status = define.STATUS_LOGIN_ING
	}

	loginResInfo := proto.NewLoginRes(define.CMD_LOGIN_RES_NO, loginReqInfo.Version, loginReqInfo.UserName, status)
	//预返回成功
	if err := (*connInfo.Conn).AsyncWrite(loginResInfo.ToBytes()); err != nil {
		log.Printf("sendto %s content %v ,err: %v\n", (*connInfo.Conn).RemoteAddr().String(), loginResInfo, err)
	}

	//进入缓存chan
	if !inWaitQueFlag {
		clientInfo := &define.ClientInfo{
			UserName: loginReqInfo.UserName,
			ConnAddr: (*connInfo.Conn).RemoteAddr().String(),
		}
		queue.WaitNumMap.Store(clientInfo.UserName, queue.IncrLoginCurNum())
		queue.EnqueueChan <- *clientInfo
	}
}
