package proto

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type (
	QueryPlayerLoginQuePosReq struct {
		ProtoHeader
		QueryPlayerLoginQuePosReqBody
	}

	QueryPlayerLoginQuePosRes struct {
		ProtoHeader
		QueryPlayerLoginQuePosResBody
	}

	QueryPlayerLoginQuePosReqBody struct {
		UserName string
	}

	QueryPlayerLoginQuePosResBody struct {
		QueWaitPlayersNum int32 //玩游戏人数
		QuePlayerPos      int32 //在队列中位置
		PlayersGameIngNum int32 //玩游戏人数
	}
)

func NewQueryPlayerLoginQuePosReq(cmdNo int64, version string, userName string) *QueryPlayerLoginQuePosReq {
	info := &QueryPlayerLoginQuePosReq{}

	bodyLen := int32(len(userName))
	info.ProtoHeader = *NewReqHead(cmdNo, version, bodyLen)

	info.UserName = userName

	//log.Printf("NewQueryPlayerLoginQuePosReq %+v\n", info)

	return info
}

func ParseToQueryPlayerLoginQuePosReq(req []byte) *QueryPlayerLoginQuePosReq {
	info := &QueryPlayerLoginQuePosReq{}

	info.ProtoHeader = *ParseToReqHead(req)
	curLen := int(info.HeaderLen)
	info.UserName = string(req[curLen:])

	//log.Printf("ParseToQueryPlayerLoginQuePosReq %+v\n", info)

	return info
}

func (info *QueryPlayerLoginQuePosReq) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	//log.Printf("QueryPlayerLoginQuePosReq ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}

func NewQueryPlayerLoginQuePosRes(cmdNo int64, version string, queWaitPlayersNum, playersGameIngNum, quePlayerPos int32) *QueryPlayerLoginQuePosRes {
	info := &QueryPlayerLoginQuePosRes{}

	bodyLen := int32(unsafe.Sizeof(info.QueryPlayerLoginQuePosResBody.QueWaitPlayersNum) +
		unsafe.Sizeof(info.QueryPlayerLoginQuePosResBody.PlayersGameIngNum) +
		unsafe.Sizeof(info.QueryPlayerLoginQuePosResBody.QuePlayerPos))
	info.ProtoHeader = *NewReqHead(cmdNo, version, bodyLen)

	info.QueWaitPlayersNum = queWaitPlayersNum
	info.PlayersGameIngNum = playersGameIngNum
	info.QuePlayerPos = quePlayerPos

	//log.Printf("NewQueryPlayerLoginQuePosReq %+v\n", info)

	return info
}

func ParseToQueryPlayerLoginQuePosRes(res []byte) *QueryPlayerLoginQuePosRes {
	info := &QueryPlayerLoginQuePosRes{}

	info.ProtoHeader = *ParseToReqHead(res)

	curLen := int(info.HeaderLen)
	endLen := curLen + int(unsafe.Sizeof(info.QueWaitPlayersNum))
	info.QueWaitPlayersNum = int32(binary.BigEndian.Uint32(res[curLen:endLen]))

	curLen = endLen
	endLen = curLen + int(unsafe.Sizeof(info.QuePlayerPos))
	info.QuePlayerPos = int32(binary.BigEndian.Uint32(res[curLen:endLen]))

	curLen = endLen
	endLen = curLen + int(unsafe.Sizeof(info.PlayersGameIngNum))
	info.PlayersGameIngNum = int32(binary.BigEndian.Uint32(res[curLen:endLen]))

	//log.Printf("ParseToQueryPlayerLoginQuePosRes %+v\n", info)

	return info
}

func (info *QueryPlayerLoginQuePosRes) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	binary.Write(resBuf, binary.BigEndian, info.QueWaitPlayersNum)
	binary.Write(resBuf, binary.BigEndian, info.QuePlayerPos)
	binary.Write(resBuf, binary.BigEndian, info.PlayersGameIngNum)

	//log.Printf("QueryPlayerLoginQuePosRes ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}
