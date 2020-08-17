package server

import (
	"net"
)

type Handler struct {
	Version uint32
	TaskID  uint32
	Handler *func(net.Conn) error
}

// NewHandler 建立新處理
func NewHandler(version uint32, id uint32, handler func(net.Conn) error) *Handler {
	return &Handler{
		Version: version,
		TaskID:  id,
		Handler: &handler,
	}
}
