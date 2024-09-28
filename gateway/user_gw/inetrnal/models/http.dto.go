package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"vpeer_usergw/inetrnal/types"
)

type RegisterRoomDTO struct {
	UserId     string
	UserLength int32  `json:"user_length"`
	Scheduled  *int64 `json:"scheduled"`
}

type RoomId struct {
	RoomId string `json:"room_id"`
}
type UserIdDTO struct {
	UserId string `json:"user_id"`
}
type RoomDTO struct {
	RoomId      string      `json:"room_id"`
	Creator     string      `json:"creator"`
	CreatorInfo UserInfoDTO `json:"creator_info"`
	UsersLength int32       `json:"users_length"`
	Closed      bool        `json:"closed"`
	CreatedAt   int64       `json:"created_at"`
	Scheduled   *int64      `json:"scheduled_at,omitempty"`
	ExpireAt    *int64      `json:"expire_at,omitempty"`
	ArchivedAt  *int64      `json:"archived_at,omitempty"`
	ClosedAt    *int64      `json:"closed_at,omitempty"`
}

type OnGoingRoomDTO struct {
	RoomId      string        `json:"room_id"`
	Creator     string        `json:"creator"`
	UsersLength int32         `json:"users_length"`
	Users       []UserInfoDTO `json:"users"`
	CreatedAt   int64         `json:"created_at"`
}

type AddFileDTO struct {
	FileId string
	RoomId string
	UserId *string
}

type FileDTO struct {
	FileId    string      `json:"file_id"`
	FileType  string      `json:"file_type"`
	RoomId    string      `json:"room_id"`
	UserId    *string     `json:"user_id"`
	User      UserInfoDTO `json:"user"`
	CreatedAt int64       `json:"created_at"`
}

type RoomResultsCount struct {
	Success int32 `json:"success"`
	Failed  int32 `json:"failed"`
	OnGoing int32 `json:"on_going"`
}

type Pagination struct {
	Offset   int32 `json:"offset"`
	Limit    int32 `json:"limit"`
	GetTotal bool  `json:"get_total"`
}

// type File struct {
// 	FileType  string    `json:"file_type"`
// 	FileId    string    `json:"file_id"`
// 	RoomId    string    `json:"room_id"`
// 	UserId    *string   `json:"user_id"`
// 	CreatedAt time.Time `json:"created_at"`
// }

type FilesDTO struct {
	Files []FileDTO `json:"files"`
	Total *int32    `json:"total,omitempty"`
}
type RoomsDTO struct {
	Rooms []RoomDTO `json:"rooms"`
	Total *int32    `json:"total,omitempty"`
}

type RoomLogDTO struct {
	RoomId    string      `json:"room_id"`
	UserId    string      `json:"user_id"`
	User      UserInfoDTO `json:"user"`
	UserEvent string      `json:"user_event"`
	CreatedAt int64       `json:"created_at"`
}

type RoomLogsDTO struct {
	RoomLogs []RoomLogDTO `json:"room_logs"`
	Total    *int32       `json:"total,omitempty"`
}

type RoomResultDTO struct {
	Approvers []UserInfoDTO `json:"approvers"`
	Deniers   []UserInfoDTO `json:"deniers"`
	CreatedAt int64         `json:"created_at"`
}

type RoomUserDTO struct {
	User            UserInfoDTO `json:"user"`
	Status          string      `json:"status"`
	FirstOccurrence int64       `json:"first_occurrence"`
}

type RoomUsersDTO struct {
	RoomUsers []RoomUserDTO `json:"users"`
	Total     *int32        `json:"total,omitempty"`
}

type UserInfoDTO struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NationalCode string `json:"national_code"`
	Username     string `json:"username"`
}

// deepCopyUserInfoDTO creates a deep copy of a UserInfoDTO object.
func DeepCopyUserInfoDTO(user *UserInfoDTO) UserInfoDTO {
	return UserInfoDTO{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		NationalCode: user.NationalCode,
		Username:     user.Username,
	}
}

func GetUsersInfoByIds(userIds []string) (map[string]UserInfoDTO, *types.Error) {
	// Marshal the user IDs to JSON
	jsonData, err := json.Marshal(userIds)
	if err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 401")
	}
	fmt.Println("HTTP => ", string(jsonData))
	// Create a new HTTP POST request with the JSON payload
	req, err := http.NewRequest("POST", "https://estate.sedrehgroup.ir/api/meets/users/detail/", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 405")
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 402")
	}

	// Close the response body when done reading from it
	defer resp.Body.Close()

	// Check if the response status code is okay
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, types.NewBadRequestError("خطای داخلی رخ داده است. کد خطا 403")
	}

	// Decode the JSON response into a map of user ID to UserInfoDTO
	userInfoMap := make(map[string]UserInfoDTO)
	users := make([]UserInfoDTO, 0)
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 404")
	}
	for _, u := range users {
		userInfoMap[u.Id] = u
	}

	return userInfoMap, nil
}
