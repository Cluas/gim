package job

const (
	OP_SEND                = int32(1) //
	OP_SINGLE_SEND         = int32(2) // 指定用户发送
	OP_ROOM_SEND           = int32(3) // 广播到房间操作
	OP_ROOM_COUNT_SEND     = int32(4) // 在线人数操作
	OP_ROOM_INFO_SEND      = int32(5) // 用户信息发送操作
	OP_ROOM_INFO_LESS_SEND = int32(6) // 用户在线列表减少用户
	OP_ROOM_INFO_ADD_SEND  = int32(6) //用户在线列表增加用户
)

type RedisMsg struct {
	Op           int32             `json:"op"`
	ServerID     int8              `json:"serverId,omitempty"`
	RoomID       int32             `json:"roomId,omitempty"`
	UserID       string            `json:"userId,omitempty"`
	Msg          []byte            `json:"msg"`
	Count        int               `json:"count"`
	RoomUserInfo map[string]string `json:"RoomUserInfo"`
}
