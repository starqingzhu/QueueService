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

func run(userName string,wg *sync.WaitGroup){
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

	packVer := "v1.0.0"
	loginInfo :=  proto.NewLoginReq(define.CMD_LOGIN_REQ_NO, packVer,userName)

	fc := goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)

	loginInfoBuf := &bytes.Buffer{}

	binary.Write(loginInfoBuf,binary.BigEndian,loginInfo.ToBytes())
	err = fc.WriteFrame(loginInfoBuf.Bytes())
	if err != nil {
		panic(err)
	}

	//解析login回复
	buf, err := fc.ReadFrame()
	if err != nil {
		panic(err)
	}
	resp := proto.ParseToLoginRes(buf)
	fmt.Printf("received res: %+v\n", resp)
	//login notify
	notifyBuf, err := fc.ReadFrame()
	if err != nil {
		panic(err)
	}
	notify := proto.ParseToLoginNotify(notifyBuf)
	fmt.Printf("received notify: %+v\n", notify)
}

func TestClient(t *testing.T) {
	var goNum int = 280
	var wg sync.WaitGroup
	wg.Add(goNum)
	for i := 1; i <= goNum; i++ {
		go run(strconv.Itoa(i),&wg)
	}
	wg.Wait()
}