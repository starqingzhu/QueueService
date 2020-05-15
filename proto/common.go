package proto

import (
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
