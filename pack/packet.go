package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"io"
	"net"

	"github.com/c-delta/torii/utils"
)

// Packet 定義封包
type Packet struct {
	MaxPacketSize uint32
	Connection    net.Conn
}

// NewPacket 建立新封包
func NewPacket() *Packet {
	return &Packet{
		MaxPacketSize: uint32(utils.Config().MaxPacketSize),
	}
}

// GetHeaderLength 獲取封包標頭長度
func (p *Packet) GetHeaderLength() uint32 {
	// uint32 = 4 bytes
	// 4 * 3 = 12
	return 12
}

// Pack 封包封裝
func (p *Packet) Pack(msg Message) ([]byte, error) {
	Buffer := bytes.NewBuffer([]byte{})

	// 在緩衝區寫入版本資訊
	if err := binary.Write(Buffer, binary.BigEndian, msg.GetVersion()); err != nil {
		return nil, err
	}
	// 在緩衝區寫入任務ID
	if err := binary.Write(Buffer, binary.BigEndian, msg.GetID()); err != nil {
		return nil, err
	}
	// 在緩衝區寫入資料長度
	if err := binary.Write(Buffer, binary.BigEndian, msg.GetDataLength()); err != nil {
		return nil, err
	}
	// fmt.Println(msg.GetData())
	// 在緩衝區寫入資料
	if err := binary.Write(Buffer, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return Buffer.Bytes(), nil
}

// Unpack 解開封包
func (p *Packet) Unpack(buffer []byte) (*Message, error) {
	// 建立一個 ioReader 讀取緩衝區內容
	pBuffer := bytes.NewReader(buffer)

	msg := &Message{}

	// 讀取版本
	if err := binary.Read(pBuffer, binary.BigEndian, &msg.Version); err != nil {
		return nil, err
	}
	// 讀取版本
	if err := binary.Read(pBuffer, binary.BigEndian, &msg.ID); err != nil {
		return nil, err
	}
	// 讀取版本
	if err := binary.Read(pBuffer, binary.BigEndian, &msg.DataLength); err != nil {
		return nil, err
	}

	if p.MaxPacketSize > 0 && msg.DataLength > p.MaxPacketSize {
		return nil, errors.New(fmt.Sprintf("msg data is too large. DataLength: %d", msg.DataLength))
	}

	return msg, nil
}

// Write 發送封包
func (p *Packet) Write(msg Message, conn net.Conn) error {
	buf, err := p.Pack(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	if err != nil {
		return err
	}
	return nil

}

// Read 讀取緩衝區
func (p *Packet) Read(conn net.Conn) (*Message, error) {
	headerBuffer := make([]byte, p.GetHeaderLength())
	_, err := io.ReadFull(conn, headerBuffer)
	if err != nil {
		return nil, err
		panic(err)
	}

	headOnlyMsg, err := p.Unpack(headerBuffer)
	if err != nil {
		return nil, err
	}

	if headOnlyMsg.GetDataLength() > 0 {
		msg := headOnlyMsg
		msg.Data = make([]byte, msg.GetDataLength())

		_, err = io.ReadFull(conn, msg.Data)
		// fmt.Println(msg.Data)
		if err != nil {
			return nil, err
		}
		return msg, nil
	} else {
		return headOnlyMsg, errors.New("msg is empty")
	}

}
