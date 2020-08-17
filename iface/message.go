package iface

type IMessage interface {
	GetVersion() uint32
	GetID() uint32
	GetDataLength() uint32
	GetData() []byte

	SetVersion(uint32)
	SetID(uint32)
	SetDataLength(uint32)
	SetData([]byte)
}
