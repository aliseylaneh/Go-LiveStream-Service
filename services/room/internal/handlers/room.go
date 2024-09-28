// Package handlers provides gRPC service handlers for the "Room" service.
package handlers

import (
	"context"
	"time"
	"vpeer_room/internal/models"
	"vpeer_room/internal/services"
	pb "vpeer_room/proto/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RoomHandler is a gRPC server handler for the "Room" service.
type RoomHandler struct {
	pb.UnimplementedRoomServiceServer
	roomService services.RoomService
}

// NewRoomHandler creates a new RoomHandler instance with the provided roomService.
func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

// RegisterRoom is a gRPC handler for registering a room.
func (c *RoomHandler) RegisterRoom(ctx context.Context, request *pb.RegisterRoomRequest) (*pb.RegisterRoomResponse, error) {
	data := models.Room{
		Creator:     request.Creator,
		UsersLength: request.UsersLength,
	}

	// Parse the string to a time.Time value in UTC format
	if request.Scheduled != nil {
		if *request.Scheduled != 0 {
			scheduledTime := time.Unix(*request.Scheduled, 0)
			data.Scheduled = &scheduledTime
		} else {
			return nil, status.Error(codes.Aborted, "زمان شروع جلسه نباید خالی باشد. کد خطا 14")
		}
	}

	if request.RoomExpiry != nil {
		if *request.RoomExpiry != 0 {
			expiryTime := time.Unix(*request.RoomExpiry, 0)
			data.ExpireAt = &expiryTime
		} else {
			return nil, status.Error(codes.Aborted, "زمان پایان جلسه نباید خالی باشد. کد خطا 2-14")
		}
	}

	// Call the roomService to register the room.
	res, err := c.roomService.RegisterRoom(&data)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	// Return the response with the registration link.
	return &pb.RegisterRoomResponse{Link: res}, nil
}

func (c *RoomHandler) CloseRoom(ctx context.Context, request *pb.CloseRoomByRoomIdRequest) (*pb.Empty, error) {
	err := c.roomService.CloseRoomByRoomId(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.Empty{}, nil
}

func (c *RoomHandler) CloseRoomByRoomId(ctx context.Context, request *pb.CloseRoomByRoomIdRequest) (*pb.Empty, error) {
	err := c.roomService.CloseRoomByRoomId(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.Empty{}, nil
}

func (c *RoomHandler) GetRoomsByUserid(ctx context.Context, request *pb.GetRoomByUserId) (*pb.Rooms, error) {
	res, err := c.roomService.GetRoomsByUserId(request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return roomsToProtoList(&models.Rooms{Rooms: res}), nil
}

func (c *RoomHandler) GetRoomByRoomid(ctx context.Context, request *pb.GetRoomByRoomId) (*pb.Room, error) {
	res, err := c.roomService.GetRoomByRoomId(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return roomToProto(res), nil
}

func (c *RoomHandler) GetOpenRoomByUserid(ctx context.Context, request *pb.GetRoomByUserId) (*pb.Room, error) {
	res, err := c.roomService.GetOpenRoomByUserId(request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return roomToProto(res), nil
}

func (c *RoomHandler) GetRoomCreatorByRoomid(ctx context.Context, request *pb.GetRoomByRoomId) (*pb.GetRoomByUserId, error) {
	res, err := c.roomService.GetRoomCreatorByRoomid(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.GetRoomByUserId{UserId: res}, nil
}

func (c *RoomHandler) CheckRoomJoinable(ctx context.Context, request *pb.IsRoomJoinableRequest) (*pb.IsRoomJoinableResponse, error) {
	res, err := c.roomService.IsRoomJoinable(request.RoomId, request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.IsRoomJoinableResponse{Status: res}, nil
}

func (c *RoomHandler) GetRooms(ctx context.Context, request *pb.Pagination) (*pb.Rooms, error) {
	pagination := &models.Pagination{
		Offset:   request.Offset,
		Limit:    request.Limit,
		GetTotal: request.GetTotal,
	}

	res, err := c.roomService.GetRooms(pagination)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return roomsToProtoList(res), nil
}

func (c *RoomHandler) GetOpenRooms(ctx context.Context, request *pb.Pagination) (*pb.Rooms, error) {
	pagination := &models.Pagination{
		Offset:   request.Offset,
		Limit:    request.Limit,
		GetTotal: request.GetTotal,
	}

	res, err := c.roomService.GetOpenRooms(pagination)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return roomsToProtoList(res), nil
}

func (c *RoomHandler) AddRoomLog(ctx context.Context, request *pb.AddRoomLog) (*pb.Empty, error) {
	roomLog := &models.RoomLog{
		RoomId:    request.RoomId,
		UserId:    request.UserId,
		UserEvent: request.UserEvent,
	}
	err := c.roomService.AddRoomLog(roomLog)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.Empty{}, nil
}

func (c *RoomHandler) GetRoomLogsByRoomid(ctx context.Context, request *pb.GetRoomByRoomId) (*pb.RoomLogs, error) {
	res, err := c.roomService.GetRoomLogsByRoomId(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return roomLogsToProtoList(&models.RoomLogs{RoomLogs: res}), nil
}

func (c *RoomHandler) AddRoomResult(ctx context.Context, request *pb.AddRoomResult) (*pb.Empty, error) {
	roomResult := &models.RoomResult{
		RoomId:    request.RoomId,
		Approvers: request.Approvers,
		Deniers:   request.Deniers,
	}
	err := c.roomService.AddRoomResult(roomResult)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.Empty{}, nil
}

func (c *RoomHandler) GetRoomResults(ctx context.Context, request *pb.Pagination) (*pb.RoomResults, error) {
	pagination := &models.Pagination{
		Offset:   request.Offset,
		Limit:    request.Limit,
		GetTotal: request.GetTotal,
	}

	res, err := c.roomService.GetAllRoomResults(pagination)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return roomResultsToProtoList(res), nil
}

func (c *RoomHandler) GetRoomResultByRoomid(ctx context.Context, request *pb.GetRoomByRoomId) (*pb.RoomResult, error) {
	res, err := c.roomService.GetRoomResultByRoomId(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return roomResultToProto(res), nil
}

func (c *RoomHandler) GetRoomResultsCount(ctx context.Context, request *pb.Empty) (*pb.RoomResultsCount, error) {
	res, err := c.roomService.GetAllRoomResultsCount()
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.RoomResultsCount{Success: res.Success, Failed: res.Failed}, nil
}

func (c *RoomHandler) GetAllUsers(ctx context.Context, request *pb.Pagination) (*pb.RoomUsers, error) {
	pagination := &models.Pagination{
		Offset:   request.Offset,
		Limit:    request.Limit,
		GetTotal: request.GetTotal,
	}

	res, err := c.roomService.GetAllUsers(pagination)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	roomUsers := make([]*pb.RoomUser, 0)
	for _, u := range res.RoomUsers {
		roomUsers = append(roomUsers, &pb.RoomUser{
			UserId:          u.UserId,
			Status:          u.Status,
			FirstOccurrence: u.FirstOccurrence.Unix(),
		})
	}
	return &pb.RoomUsers{RoomUsers: roomUsers, Total: res.Total}, nil
}

func (c *RoomHandler) AddBanUser(ctx context.Context, request *pb.GetRoomByUserId) (*pb.Empty, error) {
	err := c.roomService.AddBanUser(request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.Empty{}, nil
}

func (c *RoomHandler) RemoveBanUser(ctx context.Context, request *pb.GetRoomByUserId) (*pb.Empty, error) {
	err := c.roomService.RemoveBanUser(request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.Empty{}, nil
}

// func (c *RoomHandler) GetArchivedRoomByRoomId(ctx context.Context, request *pb.GetRoomByRoomId) (*pb.Room, error) {
// 	res, err := c.roomService.GetArchivedRoomByRoomId(request.RoomId)
// 	if err != nil {
// 		return nil, err.ErrorToGRPCStatus()
// 	}
// 	return roomToProto(res), nil
// }

func roomToProto(up *models.Room) *pb.Room {
	var schaduled *int64
	if up.Scheduled != nil {
		scheduledTime := up.Scheduled.Unix()
		schaduled = &scheduledTime
	} else {
		schaduled = nil
	}

	// var archivedAt *int64
	// if up.ArchivedAt != nil {
	// 	archivedAtTime := up.ArchivedAt.Unix()
	// 	archivedAt = &archivedAtTime
	// } else {
	// 	archivedAt = nil
	// }

	var expireAt *int64
	if up.ExpireAt != nil {
		expireAtTime := up.ExpireAt.Unix()
		expireAt = &expireAtTime
	} else {
		expireAt = nil
	}

	var closedAt *int64
	if up.ClosedAt != nil {
		closedAtTime := up.ClosedAt.Unix()
		closedAt = &closedAtTime
	} else {
		closedAt = nil
	}

	return &pb.Room{
		RoomId:      up.RoomId,
		UsersLength: up.UsersLength,
		Closed:      up.Closed,
		ClosedAt:    closedAt,
		UserId:      up.Creator,
		Schaduled:   schaduled,
		CreatedAt:   up.CreatedAt.Unix(),
		// ArchivedAt:  archivedAt,
		RoomExpiry: expireAt,
	}
}

func roomsToProtoList(rooms *models.Rooms) *pb.Rooms {
	pbRooms := make([]*pb.Room, 0, len(rooms.Rooms))
	for _, room := range rooms.Rooms {
		pbRooms = append(pbRooms, roomToProto(&room))
	}
	return &pb.Rooms{Rooms: pbRooms, Total: rooms.Total}
}

func roomLogToProto(up *models.RoomLog) *pb.RoomLog {
	return &pb.RoomLog{
		RoomId:    up.RoomId,
		UserId:    up.UserId,
		UserEvent: up.UserEvent,
		CreatedAt: up.CreatedAt.Unix(),
	}
}

func roomLogsToProtoList(rooms *models.RoomLogs) *pb.RoomLogs {
	pbRooms := make([]*pb.RoomLog, 0, len(rooms.RoomLogs))
	for _, room := range rooms.RoomLogs {
		pbRooms = append(pbRooms, roomLogToProto(&room))
	}
	return &pb.RoomLogs{RoomLogs: pbRooms, Total: rooms.Total}
}

func roomResultToProto(up *models.RoomResult) *pb.RoomResult {
	return &pb.RoomResult{
		RoomId:    up.RoomId,
		Approvers: up.Approvers,
		Deniers:   up.Deniers,
		CreatedAt: up.CreatedAt.Unix(),
	}
}

func roomResultsToProtoList(rooms *models.RoomResults) *pb.RoomResults {
	pbRooms := make([]*pb.RoomResult, 0, len(rooms.RoomResults))
	for _, room := range rooms.RoomResults {
		pbRooms = append(pbRooms, roomResultToProto(&room))
	}
	return &pb.RoomResults{RoomResults: pbRooms, Total: rooms.Total}
}
