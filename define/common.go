package define

type (
	ProtoHeader struct {
		CmdNo     int64
		HeaderLen int32
		BodyLen   int32
		Version   string //固定长度
	}
)
