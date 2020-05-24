package server

import (
	"QueueService/connManager"
	"QueueService/define"
	"QueueService/proto"
	"bytes"
	"encoding/binary"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"

	"time"
)

type codecServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
}

func (cs *codecServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (cs *codecServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("client %s connect)\n", c.RemoteAddr().String())

	connInfo := define.NewConnInfo(&c)

	connManager.ConnManager.Store(c.RemoteAddr().String(), connInfo)

	return
}

func (cs *codecServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		log.Printf("client %s disconnect err: %v)\n", c.RemoteAddr().String(), err)
	} else {
		log.Printf("client %s disconnect)\n", c.RemoteAddr().String())
	}
	connManager.ConnManager.Delete(c.RemoteAddr().String())

	return
}

func (cs *codecServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

	value, ok := connManager.ConnManager.Load(c.RemoteAddr().String())
	if !ok {
		return
	}
	connInfo, ok := value.(*define.ConnInfo)
	if !ok {
		return
	}

	HandleReqInfoParse(frame, connInfo)
	//if cs.async {
	//	data := append([]byte{}, []byte{}...)
	//	_ = cs.workerPool.Submit(func() {
	//		c.AsyncWrite(data)
	//	})
	//	return
	//}
	//out = loginResInfo.ToBytes()
	return
}

/*
解析各种协议
先解析协议头、根据协议头确定协议，然后解析协议体
*/
func HandleReqInfoParse(frame []byte, connInfo *define.ConnInfo) error {
	frameBuff := bytes.NewBuffer(frame) //Next
	for frameBuff.Len() > 0 || connInfo.Buff.Free() > 0 {
		if connInfo.Buff.Free() >= frameBuff.Len() {
			if _, err := connInfo.Buff.Write(frameBuff.Next(frameBuff.Len())); err != nil {
				log.Printf("HandleReqInfoParse Write err: %v \n", err)
				return err
			}
		} else {
			if _, err := connInfo.Buff.Write(frameBuff.Next(connInfo.Buff.Free())); err != nil {
				log.Printf("HandleReqInfoParse Write err: %v \n", err)
				return err
			}
		}

		//如果连接缓存中连最低协议字符长度都没有，直接返回
		minProtoLen := proto.MinProtoLen()
		if minProtoLen >= connInfo.Buff.Length() {
			return nil
		}

		head, tail := connInfo.Buff.LazyRead(minProtoLen)
		head = append(head, tail...)

		lenInfo := proto.ParseToProtoLen(head)

		protoLen := int(lenInfo.HeaderLen + lenInfo.BodyLen)
		if protoLen > connInfo.Buff.Length() {
			return nil
		}

		head, tail = connInfo.Buff.LazyRead(protoLen)
		head = append(head, tail...)
		connInfo.Buff.Shift(protoLen)

		if handleFun, ok := MapParseHandle[lenInfo.CmdNo]; ok {
			handleFun(frame, connInfo)
		} else {
			log.Printf("HandleReqInfoParse not support cmdNo %d\n", lenInfo.CmdNo)
		}

	}

	return nil

}

func CodecServerRun(addr string, multicore, async bool, codec gnet.ICodec) {
	var err error
	if codec == nil {
		encoderConfig := gnet.EncoderConfig{
			ByteOrder:                       binary.BigEndian,
			LengthFieldLength:               4,
			LengthAdjustment:                0,
			LengthIncludesLengthFieldLength: false,
		}
		decoderConfig := gnet.DecoderConfig{
			ByteOrder:           binary.BigEndian,
			LengthFieldOffset:   0,
			LengthFieldLength:   4,
			LengthAdjustment:    0,
			InitialBytesToStrip: 4,
		}
		codec = gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
	}
	cs := &codecServer{addr: addr, multicore: multicore, async: async, codec: codec, workerPool: goroutine.Default()}
	err = gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}
