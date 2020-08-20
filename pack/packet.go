package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"io"
	"net"
	"strconv"

	"github.com/c-delta/torii/utils"
)

// Packet 定義封包
type Packet struct {
	MaxPacketSize uint32
	Connection    *net.Conn
}

// NewPacket 建立新封包
func NewPacket() *Packet {
	return &Packet{
		MaxPacketSize: uint32(utils.Config().MaxPacketSize),
	}
}

// NewMaxSizePacket 建立自定義最大封包長度
func NewMaxSizePacket(size uint32) *Packet {
	return &Packet{
		MaxPacketSize: size,
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
		return nil, fmt.Errorf("msg data is too large. DataLength: %d", msg.DataLength)
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
	}
	return headOnlyMsg, errors.New("msg is empty")

}

// SmartWrite 自動辨識內容
func (p *Packet) SmartWrite(anydata interface{}, conn net.Conn) error {
	systembits := strconv.IntSize
	var code uint32 = 1
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, anydata)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	switch v := anydata.(type) {
	// String
	case string:
		code = 2
	case []byte:
		code = 3
	case bool:
		code = 4
	case float32:
		code = 5
	case float64:
		code = 6
	case int:
		if systembits == 32 {
			code = 11
		} else {
			code = 12
		}
	case int8:
		code = 16
	case int16:
		code = 17
	case int32:
		code = 18
	case int64:
		code = 19
	case uint:
		if systembits == 32 {
			code = 31
		} else {
			code = 32
		}
	case uint8:
		code = 36
	case uint16:
		code = 37
	case uint32:
		code = 38
	case uint64:
		code = 39
	default:
		return fmt.Errorf("can't process type: %s", v)
	}
	return p.Write(*NewCustomMessage(0, code, buf.Bytes()), conn)
}

// SmartRead 自動辨識接收內容
func (p *Packet) SmartRead(conn net.Conn) (interface{}, error) {

	msg, err := p.Read(conn)
	if err != nil {
		return nil, err
	}
	systembits := strconv.IntSize
	buf := bytes.NewReader(msg.Data)
	switch msg.ID {
	// String
	case 2:
		return string(msg.Data), nil
		// []Byte
	case 3:
		return msg.Data, nil
		// Bool
	case 4:
		var boolean bool
		err := binary.Read(buf, binary.LittleEndian, &boolean)
		return boolean, err
		// Float32
	case 5:
		var float float32
		err := binary.Read(buf, binary.LittleEndian, &float)
		return float, err
		// Float64
	case 6:
		var float float64
		err := binary.Read(buf, binary.LittleEndian, &float)
		return float, err
		// Int
		// Remote system x86
	case 11:
		if systembits == 32 {
			var intx32 int32
			err := binary.Read(buf, binary.LittleEndian, &intx32)
			return int(intx32), err
		}
		var intx64 int64
		err := binary.Read(buf, binary.LittleEndian, &intx64)
		return int(intx64), err

		// Remote system x64
	case 12:
		if systembits == 32 {
			return nil, errors.New("received int from x64 system cannot be converted to int32 or int")
		}
		var intx64 int64
		err := binary.Read(buf, binary.LittleEndian, &intx64)
		return int(intx64), err

		// Int
		// Remote system x86
	case 31:
		if systembits == 32 {
			var uintx32 uint32
			err := binary.Read(buf, binary.LittleEndian, &uintx32)
			return uint(uintx32), err
		}
		var uintx64 int64
		err := binary.Read(buf, binary.LittleEndian, &uintx64)
		return uint(uintx64), err
		// Remote system x64
	case 32:
		if systembits == 32 {
			return nil, errors.New("received uint from x64 system cannot be converted to uint32 or uint")
		}
		var uintx64 uint64
		err := binary.Read(buf, binary.LittleEndian, &uintx64)
		return uint(uintx64), err

		// Int Family
	case 16:
		var d int8
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	case 17:
		var d int16
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	case 18:
		var d int32
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	case 19:
		var d int64
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
		// Uint Family
	case 36:
		var d uint8
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	case 37:
		var d uint16
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	case 38:
		var d uint32
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	case 39:
		var d uint64
		err := binary.Read(buf, binary.LittleEndian, &d)
		return d, err
	default:
		return nil, errors.New("undefined id, please check server and client version")
	}

}
