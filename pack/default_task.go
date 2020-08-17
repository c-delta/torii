package pack

func DisconnectTask(version uint32) *Message {
	return &Message{
		Version: version,
		ID: 0,
		DataLength: uint32(len([]byte("exit"))),
		Data: []byte("exit"),
	}
}