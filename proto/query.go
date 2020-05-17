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

	return info
}

func ParseToQueryPlayerLoginQuePosReq(req []byte) *QueryPlayerLoginQuePosReq {
	info := &QueryPlayerLoginQuePosReq{}

	//包头
	info.ProtoHeader = *ParseToReqHead(req)

	//包体
	curLen := int(info.HeaderLen)
	info.UserName = string(req[curLen:])

	return info
}

func (info *QueryPlayerLoginQuePosReq) ToBytes() []byte {
	resBuf := &bytes.Buffer{}

	//包头
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	//包体
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	return resBuf.Bytes()
}

func NewQueryPlayerLoginQuePosRes(cmdNo int64, version string, queWaitPlayersNum, playersGameIngNum, quePlayerPos int32) *QueryPlayerLoginQuePosRes {
	info := &QueryPlayerLoginQuePosRes{}

	bodyLen := int32(unsafe.Sizeof(info.QueryPlayerLoginQuePosResBody.QueWaitPlayersNum) +
		unsafe.Sizeof(info.QueryPlayerLoginQuePosResBody.PlayersGameIngNum) +
		unsafe.Sizeof(info.QueryPlayerLoginQuePosResBody.QuePlayerPos))
	//包头
	info.ProtoHeader = *NewReqHead(cmdNo, version, bodyLen)

	//包体
	info.QueWaitPlayersNum = queWaitPlayersNum
	info.PlayersGameIngNum = playersGameIngNum
	info.QuePlayerPos = quePlayerPos

	return info
}

func ParseToQueryPlayerLoginQuePosRes(res []byte) *QueryPlayerLoginQuePosRes {
	info := &QueryPlayerLoginQuePosRes{}

	//包头
	info.ProtoHeader = *ParseToReqHead(res)

	//包体
	curLen := int(info.HeaderLen)
	endLen := curLen + int(unsafe.Sizeof(info.QueWaitPlayersNum))
	info.QueWaitPlayersNum = int32(binary.BigEndian.Uint32(res[curLen:endLen]))

	curLen = endLen
	endLen = curLen + int(unsafe.Sizeof(info.QuePlayerPos))
	info.QuePlayerPos = int32(binary.BigEndian.Uint32(res[curLen:endLen]))

	curLen = endLen
	endLen = curLen + int(unsafe.Sizeof(info.PlayersGameIngNum))
	info.PlayersGameIngNum = int32(binary.BigEndian.Uint32(res[curLen:endLen]))

	return info
}

func (info *QueryPlayerLoginQuePosRes) ToBytes() []byte {
	resBuf := &bytes.Buffer{}

	//包头
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())

	//包体
	binary.Write(resBuf, binary.BigEndian, info.QueWaitPlayersNum)
	binary.Write(resBuf, binary.BigEndian, info.QuePlayerPos)
	binary.Write(resBuf, binary.BigEndian, info.PlayersGameIngNum)

	return resBuf.Bytes()
}
