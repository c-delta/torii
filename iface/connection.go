package iface

import "net"

type IConnection interface {
	RemoteAddr() net.Addr

	SendMsg(version uint32, task uint32, data string) error
	RecvMsg(conn net.Conn) (version uint32, task uint32, data string, err error)

	SendBuffer(version uint32, task uint32, data []byte) error
	RecvBuffer(conn net.Conn) (version uint32, task uint32, data []byte, err error)

	SendFileBuffer(version uint32, task uint32, pointer uint32, data []byte) error
	RecvFileBuffer(conn net.Conn) (version uint32, task uint32, pointer uint32, data []byte, err error)
}
