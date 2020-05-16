package main

import (
	"QueueService/define"
	"QueueService/proto"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/smallnest/goframe"
	"net"
	"strconv"
	"sync"
	"testing"
)

func run(userName string, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover err: %v ,userName: %s---->>>>\n", err, userName)
		}
	}()
	defer wg.Done()
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

	loginInfo := proto.NewLoginReq(define.CMD_LOGIN_REQ_NO, define.PROTO_VERSION, userName)

	fc := goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)

	loginInfoBuf := &bytes.Buffer{}
	binary.Write(loginInfoBuf, binary.BigEndian, loginInfo.ToBytes())
	err = fc.WriteFrame(loginInfoBuf.Bytes())
	if err != nil {
		panic(err)
	}

	//解析login回复
	buf, err := fc.ReadFrame()
	if err != nil {
		panic(err)
	}
	printProtoInfo(buf)

	//查询位置
	queryInfo := proto.NewQueryPlayerLoginQuePosReq(define.CMD_QUERY_PLAYER_LOGIN_QUE_POS_REQ_NO,
		define.PROTO_VERSION,
		userName)
	queryInfoBuf := &bytes.Buffer{}
	binary.Write(queryInfoBuf, binary.BigEndian, queryInfo.ToBytes())
	err = fc.WriteFrame(queryInfoBuf.Bytes())
	if err != nil {
		panic(err)
	}

	//读取两次，一次异步通知、一次位置查询
	notifyBuf, err := fc.ReadFrame()
	if err != nil {
		panic(err)
	}
	printProtoInfo(notifyBuf)

	notifyBuf, err = fc.ReadFrame()
	if err != nil {
		panic(err)
	}
	printProtoInfo(notifyBuf)

}

func printProtoInfo(info []byte) {
	headInfo := proto.ParseToReqHead(info)
	switch headInfo.CmdNo {
	case define.CMD_LOGIN_RES_NO:
		resp := proto.ParseToLoginRes(info)
		fmt.Printf("received res: %+v\n", resp)

	case define.CMD_LOGIN_NOTIFY_NO:
		notify := proto.ParseToLoginNotify(info)
		fmt.Printf("received notify: %+v\n", notify)

	case define.CMD_QUERY_PLAYER_LOGIN_QUE_POS_RSP_NO:
		resp := proto.ParseToQueryPlayerLoginQuePosRes(info)
		fmt.Printf("received res: %+v\n", resp)

	default:

	}
}

func TestClient(t *testing.T) {
	var goNum int = 10
	var wg sync.WaitGroup
	wg.Add(goNum)
	for i := 1; i <= goNum; i++ {
		go run(strconv.Itoa(i), &wg)
	}
	wg.Wait()
}
