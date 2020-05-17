package proto

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type (
	QuitLoginQueReq struct {
		ProtoHeader
		QuitLoginQueReqBody
	}
	QuitLoginQueRes struct {
		ProtoHeader
		QuitLoginQueResBody
	}

	QuitLoginQueReqBody struct {
		UserName string
	}
	QuitLoginQueResBody struct {
		Status   uint16
		UserName string
	}
)

func NewQuitLoginQueReq(cmdNo int64, version string, userName string) *QuitLoginQueReq {
	info := &QuitLoginQueReq{}

	bodyLen := int32(len(userName))
	info.ProtoHeader = *NewReqHead(cmdNo, version, bodyLen)

	info.UserName = userName

	return info
}

func ParseToQuitLoginQueReq(req []byte) *QuitLoginQueReq {
	info := &QuitLoginQueReq{}

	//包头
	info.ProtoHeader = *ParseToReqHead(req)

	//包体
	curLen := int(info.HeaderLen)
	info.UserName = string(req[curLen:])

	return info
}

func (info *QuitLoginQueReq) ToBytes() []byte {
	resBuf := &bytes.Buffer{}

	//包头
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	//包体
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	return resBuf.Bytes()
}

func NewQuitLoginQueRes(cmdNo int64, version string, userName string, status uint16) *QuitLoginQueRes {
	info := &QuitLoginQueRes{}

	bodyLen := int32(len(userName) + int(unsafe.Sizeof(info.Status)))
	//包头
	info.ProtoHeader = *NewReqHead(cmdNo, version, bodyLen)

	//包体
	info.UserName = userName
	info.Status = status

	return info
}

func ParseToQuitLoginQueRes(res []byte) *QuitLoginQueRes {
	info := &QuitLoginQueRes{}

	//包头
	info.ProtoHeader = *ParseToReqHead(res)
	curLen := int(info.HeaderLen)

	//包体
	statusTypeLen := int(unsafe.Sizeof(info.Status))
	endLen := curLen + statusTypeLen
	info.Status = binary.BigEndian.Uint16(res[curLen:endLen])

	curLen = endLen
	info.UserName = string(res[curLen:])

	return info
}

func (info *QuitLoginQueRes) ToBytes() []byte {
	resBuf := &bytes.Buffer{}

	//包头
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())

	//包体
	binary.Write(resBuf, binary.BigEndian, info.Status)
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	return resBuf.Bytes()
}
