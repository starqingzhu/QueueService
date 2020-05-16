package define

import (
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/ringbuffer"
)

const (
	CONN_BUFF_MAX_SIZE = 1024 //每条连接缓存 1k
)

type (
	ConnInfo struct {
		Conn *gnet.Conn
		Buff *ringbuffer.RingBuffer
	}
)

func NewConnInfo(conn *gnet.Conn) *ConnInfo {
	connInfo := &ConnInfo{
		Conn: conn,
		Buff: ringbuffer.New(CONN_BUFF_MAX_SIZE),
	}
	return connInfo
}

//func (conn *ConnInfo)Destory(){
//	if conn.Buff != nil {
//		conn.Buff.Reset()
//	}
//	conn.Buff = nil
//}
