package models

import "time"

type File struct {
	FileType  string
	FileId    string
	RoomId    string
	UserId    *string
	CreatedAt time.Time
}

type Files struct {
	Files []File
	Total *int32
}
