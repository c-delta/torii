package main

import (
	"fmt"
	"net"

	"github.com/c-delta/torii/pack"
	"github.com/c-delta/torii/server"
)

func main() {
	serv := server.NewServer(10)

	serv.NewTask(server.NewHandler(1, 20, Handler1))

	serv.Start()
}

func Handler1(conn net.Conn) error {
	pack := pack.NewPacket()

	for {
		msg, err := pack.Read(conn)
		if err != nil {
			return err
		}
		if string(msg.Data) == "quit" {
			fmt.Println("Quit")
			return nil
		}
	}

}
