package models

import "time"

// Room represents a room entity with its attributes.
type Room struct {
	RoomId      string
	Creator     string
	UsersLength int32
	Closed      bool
	ClosedAt    *time.Time
	CreatedAt   time.Time
	Scheduled   *time.Time
	ExpireAt    *time.Time
	ArchivedAt  *time.Time
}

type Rooms struct {
	Rooms []Room
	Total *int32
}

type RoomResult struct {
	RoomId    string
	Approvers []string
	Deniers   []string
	CreatedAt time.Time
}

type RoomLog struct {
	RoomId    string
	UserId    string
	UserEvent string
	CreatedAt time.Time
}

type RoomLogs struct {
	RoomLogs []RoomLog
	Total    *int32
}

type RoomResults struct {
	RoomResults []RoomResult
	Total       *int32
}

type RoomResultsCount struct {
	Success int32
	Failed  int32
}

type RoomUser struct {
	UserId          string
	Status          string
	FirstOccurrence time.Time
}

type RoomUsers struct {
	RoomUsers []RoomUser
	Total     *int32
}
