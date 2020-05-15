package main

import (
	"QueueService/connManager"
	"QueueService/define"
	"QueueService/proto"
	"QueueService/queue"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
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
	connManager.ConnManager.Store(c.RemoteAddr().String(), c)

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

	infoBytes := HandleReqInfoParse(frame, c)
	if cs.async {
		data := append([]byte{}, infoBytes...)
		_ = cs.workerPool.Submit(func() {
			c.AsyncWrite(data)
		})
		return
	}
	//out = loginResInfo.ToBytes()
	return
}

/*
解析各种协议
先解析协议头、根据协议头确定协议，然后解析协议体
*/
func HandleReqInfoParse(frame []byte, c gnet.Conn) (res []byte) {
	info := proto.ParseToReqHead(frame)
	switch info.CmdNo {
	case define.CMD_LOGIN_REQ_NO:
		loginReqInfo := proto.ParseToLoginReq(frame)
		loginResInfo := proto.NewLoginRes(define.CMD_LOGIN_RES_NO, loginReqInfo.Version, loginReqInfo.UserName, define.STATUS_LOGIN_ING)
		res = loginResInfo.ToBytes()
		clientInfo := &define.ClientInfo{
			UserName: loginReqInfo.UserName,
			ConnAddr: c.RemoteAddr().String(),
		}
		//log.Printf("queue.EnqueueChan clientInfo:%+v\n", clientInfo)
		queue.EnqueueChan <- *clientInfo

	case define.CMD_LOGIN_QUIT_REQ_NO:

	default:

	}
	return
}

func codecServerRun(addr string, multicore, async bool, codec gnet.ICodec) {
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

func main() {
	var port int
	var multicore bool

	flag.IntVar(&port, "port", 9000, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()

	runtime.GOMAXPROCS(2)
	//初始化工作
	queue.Init()
	go queue.ListenChanges()
	go queue.HandleLogin()
	go queue.OperateWaitList()

	addr := fmt.Sprintf("tcp://:%d", port)
	codecServerRun(addr, multicore, true, nil)
}
