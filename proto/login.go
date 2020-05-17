package proto

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type (
	LoginReq struct {
		ProtoHeader
		LoginReqBody
	}

	LoginRes struct {
		ProtoHeader
		LoginResBody
	}

	LoginNotify struct {
		ProtoHeader
		LoginNotifyBody
	}

	LoginReqBody struct {
		UserName string
	}

	LoginResBody struct {
		UserName string
		Status   uint16
	}

	LoginNotifyBody struct {
		UserNameLen uint16
		UserName    string
		TokenLen    uint16
		Token       string
	}
)

func NewLoginReq(cmdNo int64, version string, userName string) *LoginReq {
	info := &LoginReq{}
	info.CmdNo = cmdNo
	info.Version = version
	info.HeaderLen = int32(unsafe.Sizeof(info.CmdNo)+unsafe.Sizeof(info.BodyLen)+unsafe.Sizeof(info.HeaderLen)) + int32(len(info.Version))
	info.BodyLen = int32(len(userName))
	info.UserName = userName

	return info
}

func ParseToLoginReq(req []byte) *LoginReq {
	info := &LoginReq{}

	info.ProtoHeader = *ParseToReqHead(req)
	info.UserName = string(req[info.HeaderLen:])

	return info
}

func (info *LoginReq) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))

	return resBuf.Bytes()
}

func NewLoginRes(cmdNo int64, version string, userName string, status uint16) *LoginRes {
	info := &LoginRes{}
	info.CmdNo = cmdNo
	info.Version = version
	info.HeaderLen = int32(unsafe.Sizeof(info.CmdNo)+unsafe.Sizeof(info.BodyLen)+unsafe.Sizeof(info.HeaderLen)) + int32(len(info.Version))
	info.BodyLen = int32(len(userName) + int(unsafe.Sizeof(info.Status)))
	info.UserName = userName
	info.Status = status

	return info
}

func ParseToLoginRes(res []byte) *LoginRes {
	info := &LoginRes{}

	info.ProtoHeader = *ParseToReqHead(res)

	statusLen := int32(unsafe.Sizeof(info.Status))
	endLen := info.HeaderLen + info.BodyLen - statusLen
	info.UserName = string(res[info.HeaderLen:endLen])

	info.Status = binary.BigEndian.Uint16(res[endLen:])

	return info
}

func (info *LoginRes) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	binary.Write(resBuf, binary.BigEndian, info.ProtoHeader.ToBytes())
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))
	binary.Write(resBuf, binary.BigEndian, info.Status)

	//log.Printf("LoginRes ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}

func NewLoginNotify(cmdNo int64, version string, userName string, token string) *LoginNotify {
	//包体
	bodyInfo := LoginNotifyBody{}
	bodyInfo.UserNameLen = uint16(len(userName))
	bodyInfo.UserName = userName
	bodyInfo.TokenLen = uint16(len(token))
	bodyInfo.Token = token

	//整个包
	info := &LoginNotify{}
	info.CmdNo = cmdNo
	info.Version = version
	info.HeaderLen = int32(unsafe.Sizeof(info.CmdNo)+unsafe.Sizeof(info.BodyLen)+unsafe.Sizeof(info.HeaderLen)) + int32(len(info.Version))
	info.BodyLen = int32(bodyInfo.UserNameLen + bodyInfo.TokenLen)
	info.LoginNotifyBody = bodyInfo

	return info
}

func ParseToLoginNotify(notify []byte) *LoginNotify {

	infoHead := ParseToReqHead(notify)
	info := &LoginNotify{}
	info.ProtoHeader = *infoHead

	curLen := info.HeaderLen
	userNameTypeLen := int32(unsafe.Sizeof(info.UserNameLen))
	endLen := info.HeaderLen + userNameTypeLen
	info.UserNameLen = binary.BigEndian.Uint16(notify[curLen:endLen])

	curLen = endLen
	endLen = curLen + int32(info.UserNameLen)
	info.UserName = string(notify[curLen:endLen])

	curLen = endLen
	tokenTypeLen := int32(unsafe.Sizeof(info.TokenLen))
	endLen += tokenTypeLen
	info.TokenLen = binary.BigEndian.Uint16(notify[curLen:endLen])

	curLen = endLen
	info.Token = string(notify[curLen:])

	//log.Printf("ParseToLoginNotify %+v\n", info)

	return info
}

func (info *LoginNotify) ToBytes() []byte {
	resBuf := &bytes.Buffer{}
	//包头
	binary.Write(resBuf, binary.BigEndian, info.CmdNo)
	binary.Write(resBuf, binary.BigEndian, info.HeaderLen)
	binary.Write(resBuf, binary.BigEndian, info.BodyLen)
	binary.Write(resBuf, binary.BigEndian, []byte(info.Version))

	//包体
	binary.Write(resBuf, binary.BigEndian, info.UserNameLen)
	binary.Write(resBuf, binary.BigEndian, []byte(info.UserName))
	binary.Write(resBuf, binary.BigEndian, info.TokenLen)
	binary.Write(resBuf, binary.BigEndian, []byte(info.Token))

	//log.Printf("LoginNotify ToBytes: %x len:%d\n", resBuf.Bytes(), resBuf.Len())

	return resBuf.Bytes()
}
