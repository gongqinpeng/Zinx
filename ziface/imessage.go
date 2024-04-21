package ziface

type IMessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetData() []byte
	SetMsgId(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}
