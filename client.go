package main

import (
	"net"

	"github.com/c-delta/torii/pack"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		return
	}
	defer conn.Close()
	packet := pack.NewPacket()

	packet.Write(*pack.TaskSelector(1, 20), conn)


	msg := pack.NewMessage([]byte("hello"))
	packet.Write(*msg, conn)
	msg = pack.NewMessage([]byte("quit"))
	packet.Write(*msg, conn)
	// packet.Write(*pack.TaskSelector(1, 0), conn)
}
