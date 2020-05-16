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

	//log.Printf("NewQueryPlayerLoginQuePosReq %+v\n", info)

	return info
}

func ParseToQuitLoginQueReq(req []byte) *QuitLoginQueReq {
	info := &QuitLoginQueReq{}

	info.ProtoHeader = *ParseToReqHead(req)
	curLen := int(info.HeaderLen)
	info.UserName = string(req[curLen:])

	//log.Printf("ParseToQueryPlayerLoginQuePosReq %+v\n", info)

	return info
}

func (info *QuitLoginQueReq) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	//log.Printf("QueryPlayerLoginQuePosReq ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}

func NewQuitLoginQueRes(cmdNo int64, version string, userName string, status uint16) *QuitLoginQueRes {
	info := &QuitLoginQueRes{}

	bodyLen := int32(len(userName) + int(unsafe.Sizeof(info.Status)))
	info.ProtoHeader = *NewReqHead(cmdNo, version, bodyLen)

	info.UserName = userName
	info.Status = status

	//log.Printf("NewQuitLoginQueRes %+v\n", info)

	return info
}

func ParseToQuitLoginQueRes(res []byte) *QuitLoginQueRes {
	info := &QuitLoginQueRes{}

	info.ProtoHeader = *ParseToReqHead(res)
	curLen := int(info.HeaderLen)

	statusTypeLen := int(unsafe.Sizeof(info.Status))
	endLen := curLen + statusTypeLen
	info.Status = binary.BigEndian.Uint16(res[curLen:endLen])

	curLen = endLen
	info.UserName = string(res[curLen:])

	//log.Printf("ParseToQuitLoginQueRes %+v\n", info)

	return info
}

func (info *QuitLoginQueRes) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	binary.Write(resBuf, binary.BigEndian, info.Status)
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	//log.Printf("QuitLoginQueRes ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}
