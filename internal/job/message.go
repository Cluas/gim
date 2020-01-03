package job

// Operation defines the type of operation.
type Operation int

const (
	// OpSend 发送
	_ Operation = iota + 1 //
	// OpSingleSend 指定用户发送
	OpSingleSend
	// OpRoomSend 广播到房间操作
	OpRoomSend
)

// RedisMsg is struct of RedisMsg
type RedisMsg struct {
	Carrier      []byte            `json:"carrier"` // 携带trace信息
	Op           Operation         `json:"op"`
	ServerID     int8              `json:"serverID,omitempty"`
	RoomID       string            `json:"roomID,omitempty"`
	UserID       string            `json:"userID,omitempty"`
	Msg          []byte            `json:"msg"`
	Count        int               `json:"count"`
	RoomUserInfo map[string]string `json:"RoomUserInfo"`
}
