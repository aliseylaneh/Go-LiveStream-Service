package global

import (
	"time"
	pb_file "vpeer_usergw/proto/api/file"
	pb_room "vpeer_usergw/proto/api/room"
)

var (
	ROOM_SERVER_CLIENT pb_room.RoomServiceClient
	FILE_SERVER_CLIENT pb_file.FileServiceClient
	TOKENS             = make(map[string]*TokenInformation)
)

type TokenInformation struct {
	RoomId   string
	UserId   string
	ExpireAt time.Time
}
