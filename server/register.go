package server

import "QueueService/define"

//解析函数签名
type ParseAndBusinessHandle func(req []byte, connInfo *define.ConnInfo)

var MapParseHandle = make(map[int64]ParseAndBusinessHandle)

func RegisterParseHandle(cmdNo int64, handle ParseAndBusinessHandle) {
	MapParseHandle[cmdNo] = handle
}
