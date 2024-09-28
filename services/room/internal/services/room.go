package services

import (
	"safir/libs/idgen"
	"vpeer_room/internal/models"
	"vpeer_room/internal/repository"
	"vpeer_room/internal/types"
)

// RoomService is an interface defining methods for room-related services.
type RoomService interface {
	// RegisterRoom registers a new room and returns the room's unique ID or an error.
	RegisterRoom(*models.Room) (string, *types.Error)
	CloseRoomByRoomId(string) *types.Error
	GetRoomsByUserId(string) ([]models.Room, *types.Error)
	GetRoomByRoomId(string) (*models.Room, *types.Error)
	GetOpenRoomByUserId(string) (*models.Room, *types.Error)
	GetRoomCreatorByRoomid(string) (string, *types.Error)
	IsRoomJoinable(string, string) (string, *types.Error)
	GetRooms(*models.Pagination) (*models.Rooms, *types.Error)
	GetOpenRooms(*models.Pagination) (*models.Rooms, *types.Error)
	AddRoomLog(*models.RoomLog) *types.Error
	GetRoomLogsByRoomId(string) ([]models.RoomLog, *types.Error)
	AddRoomResult(*models.RoomResult) *types.Error
	GetAllRoomResults(*models.Pagination) (*models.RoomResults, *types.Error)
	GetRoomResultByRoomId(string) (*models.RoomResult, *types.Error)
	GetAllRoomResultsCount() (*models.RoomResultsCount, *types.Error)
	GetAllUsers(*models.Pagination) (*models.RoomUsers, *types.Error)
	AddBanUser(string) *types.Error
	RemoveBanUser(string) *types.Error
	// GetArchivedRoomByRoomId(string) (*models.Room, *types.Error)
}

// roomService is an implementation of the RoomService interface.
type roomService struct {
	repository repository.RoomRepository
}

// NewRoomService creates and returns a new RoomService instance.
func NewRoomService(repository repository.RoomRepository) RoomService {
	return &roomService{
		repository: repository,
	}
}

// RegisterRoom registers a new room and returns the room's unique ID or an error.
func (c *roomService) RegisterRoom(data *models.Room) (string, *types.Error) {
	var roomId string
	var err *types.Error

	// Attempt to generate a unique roomId multiple times
	for i := 0; i < 10; i++ {
		id, rerr := idgen.NextLowerAlphabeticString(10)
		if rerr != nil {
			return "", types.NewInternalError("خطای داخلی رخ داده است. خطا کد 13")
		}

		roomId = id[:3] + "-" + id[3:7] + "-" + id[7:]
		data.RoomId = roomId
		// Check if the generated roomId already exists
		exists, repoErr := c.repository.RoomExistsByRoomId(roomId)
		if repoErr != nil {
			return "", repoErr
		}

		if !exists {
			// If roomId is unique, attempt to register the room
			err = c.repository.RegisterRoom(data)
			if err != nil {
				return "", err
			}
			// Room registered successfully, break the loop
			break
		}
	}

	if err != nil {
		return "", err
	}

	return roomId, nil
}

func (c *roomService) CloseRoomByRoomId(roomId string) *types.Error {
	return c.repository.CloseRoomByRoomId(roomId)
}

func (c *roomService) GetRoomsByUserId(roomId string) ([]models.Room, *types.Error) {
	return c.repository.GetRoomsByUserId(roomId)
}

func (c *roomService) GetRoomByRoomId(roomId string) (*models.Room, *types.Error) {
	return c.repository.GetRoomByRoomId(roomId)
}

func (c *roomService) GetOpenRoomByUserId(userId string) (*models.Room, *types.Error) {
	return c.repository.GetOpenRoomByUserId(userId)
}

func (c *roomService) GetRoomCreatorByRoomid(roomId string) (string, *types.Error) {
	return c.repository.GetRoomCreatorByRoomid(roomId)
}

func (c *roomService) IsRoomJoinable(roomId string, userId string) (string, *types.Error) {
	return c.repository.IsRoomJoinable(roomId, userId)
}

func (c *roomService) GetRooms(data *models.Pagination) (*models.Rooms, *types.Error) {
	res, err := c.repository.GetRooms(data)
	if err != nil {
		return nil, err
	}
	var count *int32
	if data.GetTotal {
		res, err := c.repository.GetRoomsTotalCount()
		if err != nil {
			return nil, err
		}
		count = &res
	}
	return &models.Rooms{Rooms: res, Total: count}, nil
}

func (c *roomService) GetOpenRooms(data *models.Pagination) (*models.Rooms, *types.Error) {
	res, err := c.repository.GetOpenRooms(data)
	if err != nil {
		return nil, err
	}
	var count *int32
	if data.GetTotal {
		res, err := c.repository.GetOpenRoomsTotalCount()
		if err != nil {
			return nil, err
		}
		count = &res
	}
	return &models.Rooms{Rooms: res, Total: count}, nil
}

func (c *roomService) AddRoomLog(data *models.RoomLog) *types.Error {
	return c.repository.AddRoomLog(data)
}

func (c *roomService) GetRoomLogsByRoomId(roomId string) ([]models.RoomLog, *types.Error) {
	return c.repository.GetRoomLogsByRoomId(roomId)
}

func (c *roomService) AddRoomResult(data *models.RoomResult) *types.Error {
	return c.repository.AddRoomResult(data)
}

func (c *roomService) GetAllRoomResults(data *models.Pagination) (*models.RoomResults, *types.Error) {
	res, err := c.repository.GetAllRoomResults(data)
	if err != nil {
		return nil, err
	}
	var count *int32
	if data.GetTotal {
		res, err := c.repository.GetRoomResultTotalCount()
		if err != nil {
			return nil, err
		}
		count = &res
	}
	return &models.RoomResults{RoomResults: res, Total: count}, nil
}

func (c *roomService) GetRoomResultByRoomId(roomId string) (*models.RoomResult, *types.Error) {
	return c.repository.GetRoomResultByRoomId(roomId)
}

func (c *roomService) GetAllUsers(data *models.Pagination) (*models.RoomUsers, *types.Error) {
	res, err := c.repository.GetAllUsers(data)
	if err != nil {
		return nil, err
	}
	var count *int32
	if data.GetTotal {
		res, err := c.repository.GetAllUsersCount()
		if err != nil {
			return nil, err
		}
		count = &res
	}
	return &models.RoomUsers{RoomUsers: res, Total: count}, nil
}

func (c *roomService) AddBanUser(userId string) *types.Error {
	return c.repository.AddBanUser(userId)
}

func (c *roomService) RemoveBanUser(userId string) *types.Error {
	return c.repository.RemoveBanUser(userId)
}

// func (c *roomService) GetArchivedRoomByRoomId(roomId string) (*models.Room, *types.Error) {
// 	return c.repository.GetArchivedRoomByRoomId(roomId)
// }

// insertHyphens inserts hyphens in a string at specified intervals.
func insertHyphens(input string, interval int) string {
	if interval <= 0 {
		return input
	}

	// Preallocate the byte slice with an estimated capacity to reduce allocations.
	result := make([]byte, 0, len(input)+len(input)/interval)

	for i, char := range input {
		if i > 0 && i%interval == 0 {
			result = append(result, '-')
		}
		result = append(result, byte(char))
	}

	return string(result)
}

func (c *roomService) GetAllRoomResultsCount() (*models.RoomResultsCount, *types.Error) {
	return c.repository.GetAllRoomResultsCount()
}
