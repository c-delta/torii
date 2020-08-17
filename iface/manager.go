package iface

import "net"

type IHandlerManager interface {
	VersionHas(version uint32) (uint32, error)

	Add(handler IHandler) error
	Remove(handler IHandler) error

	Get(version uint32, taskID uint32) (*IHandler, error)
	Clear()

	AcceptTasks(conn net.Conn)
}

type ISubHandlerManager interface {
	TaskHas(version uint32) (uint32, error)

	Add(handler IHandler) error
	Remove(handler IHandler) error

	Get(taskID uint32) (*IHandler, error)
	Clear()
}
