package pack

// Message 訊息結構
type Message struct {
	Version    uint32 // 版本 (供更新用途)
	ID         uint32 // 任務ID
	DataLength uint32 // 資料長度
	Data       []byte // 資料內容
}

// NewTaskSelector 建立任務選擇封包
func TaskSelector(version uint32, id uint32) *Message {
	return &Message{
		Version:    version,
		ID:         id,
		DataLength: uint32(len([]byte(" "))),
		Data:       []byte(" "),
	}
}

// NewMessage 建立普通新封包
func NewMessage(data []byte) *Message {
	return &Message{
		Version:    0,
		ID:         1,
		DataLength: uint32(len(data)),
		Data:       data,
	}
}

// NewMsgPacket 建立自定義新封包
func NewCustomMessage(version uint32, id uint32, data []byte) *Message {
	return &Message{
		Version:    version,
		ID:         id,
		DataLength: uint32(len(data)),
		Data:       data,
	}
}

/*
	GET
*/

// GetVersion 獲取版本號
func (msg *Message) GetVersion() uint32 {
	return msg.Version
}

// GetID 獲取任務ID
func (msg *Message) GetID() uint32 {
	return msg.ID
}

// GetDataLength 獲取資料長度
func (msg *Message) GetDataLength() uint32 {
	return msg.DataLength
}

// GetData []byte 獲取資料內容
func (msg *Message) GetData() []byte {
	return msg.Data
}

/*
	SET
*/

// SetVersion 設定封包版本
func (msg *Message) SetVersion(v uint32) {
	msg.Version = v
}

// SetID 設定任務ID
// MAX: 2147483648 * 2 - 1
func (msg *Message) SetID(id uint32) {
	msg.ID = id
}

// SetDataLength 設定資料長度
func (msg *Message) SetDataLength(len uint32) {
	msg.DataLength = len
}

// SetData []byte 設定資料內容
func (msg *Message) SetData(buff []byte) {
	msg.Data = buff
}
