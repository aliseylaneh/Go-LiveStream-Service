package services

import (
	"context"
	"fmt"
	"vpeer_usergw/inetrnal/global"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/types"
	pb_room "vpeer_usergw/proto/api/room"
)

type (
	RoomService interface {
		RegisterRoom(*models.RegisterRoomDTO) (string, *types.Error)
		CloseRoom(string) *types.Error
		GetOpenRoomByUserId(string) (*models.RoomDTO, *types.Error)
		GetRoomsByUserId(string) ([]models.RoomDTO, *types.Error)
		GetRoomByRoomId(string) (*models.RoomDTO, *types.Error)
		GetCreatorByRoomId(string) (string, *types.Error)
		IsRoomJoinable(string, string) (string, *types.Error)
		GetRoomResultsCount() (*models.RoomResultsCount, *types.Error)
		GetRooms(*models.Pagination) (*models.RoomsDTO, *types.Error)
		GetRoomLogsByRoomId(string) ([]models.RoomLogDTO, *types.Error)
		GetRoomResultByRoomId(string) (*models.RoomResultDTO, *types.Error)
		GetOnGoingRooms() ([]models.OnGoingRoomDTO, *types.Error)
		GetAllUsers(*models.Pagination) (*models.RoomUsersDTO, *types.Error)
		AddBanUser(string) *types.Error
		RemoveBanUser(string) *types.Error
	}
	roomService struct {
	}
)

func NewRoomService() RoomService {
	return &roomService{}
}

func (c *roomService) RegisterRoom(data *models.RegisterRoomDTO) (string, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.RegisterRoom(context.TODO(), &pb_room.RegisterRoomRequest{
		Creator:     data.UserId,
		UsersLength: data.UserLength,
		Scheduled:   data.Scheduled,
	})

	if err != nil {
		return "", types.ExtractGRPCErrDetails(err)
	}
	return res.Link, nil
}

func (c *roomService) CloseRoom(roomId string) *types.Error {
	_, err := global.ROOM_SERVER_CLIENT.CloseRoom(context.TODO(), &pb_room.CloseRoomByRoomIdRequest{
		RoomId: roomId,
	})
	if err != nil {
		return types.ExtractGRPCErrDetails(err)
	}
	return nil
}

func (c *roomService) GetOpenRoomByUserId(userId string) (*models.RoomDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetOpenRoomByUserid(context.TODO(), &pb_room.GetRoomByUserId{
		UserId: userId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}

	mapUsers, rerr := models.GetUsersInfoByIds([]string{res.UserId})
	if rerr != nil {
		return nil, rerr
	}

	return &models.RoomDTO{
			RoomId:      res.RoomId,
			Creator:     res.UserId,
			CreatorInfo: mapUsers[res.UserId],
			UsersLength: res.UsersLength,
			Closed:      res.Closed,
			ClosedAt:    res.ClosedAt,
			CreatedAt:   res.CreatedAt,
			Scheduled:   res.Schaduled,
			ExpireAt:    res.RoomExpiry,
		},
		nil
}

func (c *roomService) GetRoomByRoomId(roomId string) (*models.RoomDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRoomByRoomid(context.TODO(), &pb_room.GetRoomByRoomId{
		RoomId: roomId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	mapUsers, rerr := models.GetUsersInfoByIds([]string{res.UserId})
	if rerr != nil {
		return nil, rerr
	}

	return &models.RoomDTO{
			RoomId:      res.RoomId,
			Creator:     res.UserId,
			CreatorInfo: mapUsers[res.UserId],
			UsersLength: res.UsersLength,
			Closed:      res.Closed,
			ClosedAt:    res.ClosedAt,
			CreatedAt:   res.CreatedAt,
			Scheduled:   res.Schaduled,
			ExpireAt:    res.RoomExpiry,
		},
		nil
}

func (c *roomService) GetRoomsByUserId(userId string) ([]models.RoomDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRoomsByUserid(context.TODO(), &pb_room.GetRoomByUserId{
		UserId: userId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}

	rooms := make([]models.RoomDTO, 0)
	userIds := make([]string, 0)
	for _, pbRoom := range res.Rooms {
		userIds = append(userIds, pbRoom.UserId)
		room := models.RoomDTO{
			RoomId:      pbRoom.RoomId,
			Creator:     pbRoom.UserId,
			UsersLength: pbRoom.UsersLength,
			Closed:      pbRoom.Closed,
			ClosedAt:    pbRoom.ClosedAt,
			CreatedAt:   pbRoom.CreatedAt,
			Scheduled:   pbRoom.Schaduled,
			ExpireAt:    pbRoom.RoomExpiry,
		}
		rooms = append(rooms, room)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}
	for i, r := range rooms {
		rooms[i].CreatorInfo = mapUsers[r.Creator]
	}
	return rooms, nil
}

func (c *roomService) GetCreatorByRoomId(roomId string) (string, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRoomCreatorByRoomid(context.TODO(), &pb_room.GetRoomByRoomId{
		RoomId: roomId,
	})
	if err != nil {
		return "", types.ExtractGRPCErrDetails(err)
	}
	return res.UserId, nil
}

func (c *roomService) IsRoomJoinable(roomId string, userId string) (string, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.CheckRoomJoinable(context.TODO(), &pb_room.IsRoomJoinableRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return "", types.ExtractGRPCErrDetails(err)
	}

	return res.Status, nil
}

func (c *roomService) GetRoomResultsCount() (*models.RoomResultsCount, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRoomResultsCount(context.Background(), &pb_room.Empty{})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	return &models.RoomResultsCount{
		Success: res.Success,
		Failed:  res.Failed,
		OnGoing: int32(len(models.Rooms)),
	}, nil
}

func (c *roomService) GetRooms(pagination *models.Pagination) (*models.RoomsDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRooms(context.Background(), &pb_room.Pagination{
		Offset:   pagination.Offset,
		Limit:    pagination.Limit,
		GetTotal: pagination.GetTotal,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	if pagination.GetTotal && pagination.Limit == 0 {
		return &models.RoomsDTO{Rooms: []models.RoomDTO{}, Total: res.Total}, nil
	}
	rooms := make([]models.RoomDTO, 0)
	userIds := make([]string, 0)
	for _, pbRoom := range res.Rooms {
		userIds = append(userIds, pbRoom.UserId)
		room := models.RoomDTO{
			RoomId:      pbRoom.RoomId,
			Creator:     pbRoom.UserId,
			UsersLength: pbRoom.UsersLength,
			Closed:      pbRoom.Closed,
			ClosedAt:    pbRoom.ClosedAt,
			CreatedAt:   pbRoom.CreatedAt,
			Scheduled:   pbRoom.Schaduled,
			ExpireAt:    pbRoom.RoomExpiry,
		}
		rooms = append(rooms, room)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	for i, r := range rooms {
		rooms[i].CreatorInfo = mapUsers[r.Creator]
	}
	return &models.RoomsDTO{Rooms: rooms, Total: res.Total}, nil
}

func (c *roomService) GetRoomLogsByRoomId(roomId string) ([]models.RoomLogDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRoomLogsByRoomid(context.Background(), &pb_room.GetRoomByRoomId{
		RoomId: roomId,
	})

	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}

	roomLogs := make([]models.RoomLogDTO, 0)
	userIds := make([]string, 0)
	for _, pbRoomLog := range res.RoomLogs {
		userIds = append(userIds, pbRoomLog.UserId)
		roomLog := models.RoomLogDTO{
			RoomId:    pbRoomLog.RoomId,
			UserId:    pbRoomLog.UserId,
			UserEvent: pbRoomLog.UserEvent,
			CreatedAt: pbRoomLog.CreatedAt,
		}
		roomLogs = append(roomLogs, roomLog)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	for i, r := range roomLogs {
		roomLogs[i].User = mapUsers[r.UserId]
	}

	return roomLogs, nil
}

func (c *roomService) GetRoomResultByRoomId(roomId string) (*models.RoomResultDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetRoomResultByRoomid(context.Background(), &pb_room.GetRoomByRoomId{
		RoomId: roomId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	fmt.Println(res.Approvers)
	fmt.Println(res.Deniers)
	mapApprovers, rerr := models.GetUsersInfoByIds(res.Approvers)
	if rerr != nil {
		return nil, rerr
	}
	approvers := make([]models.UserInfoDTO, 0)
	for i := range mapApprovers {
		approvers = append(approvers, mapApprovers[i])
	}

	mapDeniers, rerr := models.GetUsersInfoByIds(res.Deniers)
	if rerr != nil {
		return nil, rerr
	}
	deniers := make([]models.UserInfoDTO, 0)

	for i := range mapDeniers {
		deniers = append(deniers, mapDeniers[i])
	}

	return &models.RoomResultDTO{
		Approvers: approvers,
		Deniers:   deniers,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (c *roomService) GetOnGoingRooms() ([]models.OnGoingRoomDTO, *types.Error) {
	rooms := make([]models.OnGoingRoomDTO, 0)
	copy_rooms := models.GetRoomCopy()
	for roomId, room := range copy_rooms {
		userIds := make([]string, 0)
		for _, con := range room.Peers.Connections {
			userIds = append(userIds, con.UserId)
		}

		mapUsers, rerr := models.GetUsersInfoByIds(userIds)
		if rerr != nil {
			return nil, rerr
		}
		users := make([]models.UserInfoDTO, 0)
		for _, u := range mapUsers {
			users = append(users, u)
		}

		rooms = append(rooms, models.OnGoingRoomDTO{
			RoomId:      roomId,
			Creator:     room.OwnerId,
			UsersLength: room.UserLength,
			CreatedAt:   room.CreatedAt.Unix(),
			Users:       users,
		})
	}
	return rooms, nil
}

func (c *roomService) GetAllUsers(pagination *models.Pagination) (*models.RoomUsersDTO, *types.Error) {
	res, err := global.ROOM_SERVER_CLIENT.GetAllUsers(context.Background(), &pb_room.Pagination{
		Offset:   pagination.Offset,
		Limit:    pagination.Limit,
		GetTotal: pagination.GetTotal,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	if pagination.GetTotal && pagination.Limit == 0 {
		return &models.RoomUsersDTO{RoomUsers: []models.RoomUserDTO{}, Total: res.Total}, nil
	}
	userIds := make([]string, 0)
	for _, u := range res.RoomUsers {
		userIds = append(userIds, u.UserId)
	}
	users, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	roomUsers := make([]models.RoomUserDTO, 0)
	for _, u := range res.RoomUsers {
		roomUsers = append(roomUsers, models.RoomUserDTO{
			User:            users[u.UserId],
			Status:          u.Status,
			FirstOccurrence: u.FirstOccurrence,
		})
	}
	return &models.RoomUsersDTO{RoomUsers: roomUsers, Total: res.Total}, nil
}

func (c *roomService) AddBanUser(userId string) *types.Error {
	_, err := global.ROOM_SERVER_CLIENT.AddBanUser(context.Background(), &pb_room.GetRoomByUserId{
		UserId: userId,
	})
	if err != nil {
		return types.ExtractGRPCErrDetails(err)
	}
	return nil
}

func (c *roomService) RemoveBanUser(userId string) *types.Error {
	_, err := global.ROOM_SERVER_CLIENT.RemoveBanUser(context.Background(), &pb_room.GetRoomByUserId{
		UserId: userId,
	})
	if err != nil {
		return types.ExtractGRPCErrDetails(err)
	}
	return nil
}
