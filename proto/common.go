package proto

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type (
	ProtoHeader struct {
		CmdNo     int64
		HeaderLen int32
		BodyLen   int32
		Version   string //固定长度
	}
)

func NewReqHead(cmdNo int64, version string, bodyLen int32) *ProtoHeader {
	info := &ProtoHeader{}
	info.CmdNo = cmdNo
	info.Version = version
	info.HeaderLen = int32(unsafe.Sizeof(info.CmdNo)+unsafe.Sizeof(info.BodyLen)+unsafe.Sizeof(info.HeaderLen)) + int32(len(info.Version))
	info.BodyLen = bodyLen

	return info
}

func ParseToReqHead(res []byte) *ProtoHeader {
	info := &ProtoHeader{}

	cmdNoTypeLen := int(unsafe.Sizeof(info.CmdNo))
	info.CmdNo = int64(binary.BigEndian.Uint64(res[:cmdNoTypeLen]))

	curLen := cmdNoTypeLen
	headerTypeLen := int(unsafe.Sizeof(info.HeaderLen))
	info.HeaderLen = int32(binary.BigEndian.Uint32(res[curLen : curLen+headerTypeLen]))

	curLen += headerTypeLen
	bodyTypeLen := int(unsafe.Sizeof(info.BodyLen))
	info.BodyLen = int32(binary.BigEndian.Uint32(res[curLen : curLen+bodyTypeLen]))

	curLen += bodyTypeLen
	info.Version = string(res[curLen:info.HeaderLen])

	return info
}

func (info *ProtoHeader) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.CmdNo)
	binary.Write(resBuf, binary.BigEndian, info.HeaderLen)
	binary.Write(resBuf, binary.BigEndian, info.BodyLen)
	binary.Write(resBuf, binary.BigEndian, []byte(info.Version))

	//log.Printf("ProtoHeader ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}
