package main

import (
	"QueueService/define"
	"QueueService/proto"
	"bytes"
	"encoding/binary"
	"fmt"
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
	loginInfo :=  proto.NewLoginReq(define.CMD_LOGIN_REQ_NO, define.PROTO_VERSION,userName)

	fc := goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)

	loginInfoBuf := &bytes.Buffer{}

	binary.Write(loginInfoBuf,binary.BigEndian,loginInfo.ToBytes())
	err = fc.WriteFrame(loginInfoBuf.Bytes())
	if err != nil {
		panic(err)
	}

	buf, err := fc.ReadFrame()
	if err != nil {
		panic(err)
	}

	resp := proto.ParseToLoginRes(buf)

	fmt.Printf("received: %+v\n", resp)

}
