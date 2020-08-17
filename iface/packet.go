package iface

import "net"

type IPacket interface {
	GetHeaderLength() uint32

	Pack() ([]byte, error)
	Unpack() (IMessage, error)

	Write(msg IMessage, conn net.Conn) error
	Read(conn net.Conn) (IMessage, error)
}
