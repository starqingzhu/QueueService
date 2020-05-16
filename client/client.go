package main

import (
	"QueueService/define"
	"QueueService/proto"
	"bytes"
	"encoding/binary"
	"github.com/smallnest/goframe"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	encoderConfig := goframe.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}

	decoderConfig := goframe.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}

	userName := "sunbin"

	loginInfo := proto.NewLoginReq(define.CMD_LOGIN_REQ_NO, define.PROTO_VERSION, userName)

	fc := goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)

	loginInfoBuf := &bytes.Buffer{}
	binary.Write(loginInfoBuf, binary.BigEndian, loginInfo.ToBytes())

	//查询位置
	queryInfo := proto.NewQueryPlayerLoginQuePosReq(define.CMD_QUERY_PLAYER_LOGIN_QUE_POS_REQ_NO,
		define.PROTO_VERSION,
		userName)
	binary.Write(loginInfoBuf, binary.BigEndian, queryInfo.ToBytes())
	err = fc.WriteFrame(loginInfoBuf.Bytes())
	if err != nil {
		panic(err)
	}

	////解析login回复
	//buf, err := fc.ReadFrame()
	//if err != nil {
	//	panic(err)
	//}
	//printProtoInfo(buf)
	//
	////读取两次，一次异步通知、一次位置查询
	//notifyBuf, err := fc.ReadFrame()
	//if err != nil {
	//	panic(err)
	//}
	//printProtoInfo(notifyBuf)
	//
	//notifyBuf, err = fc.ReadFrame()
	//if err != nil {
	//	panic(err)
	//}
	//printProtoInfo(notifyBuf)

}
