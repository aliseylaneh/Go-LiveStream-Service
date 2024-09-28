// Package handlers manages the interaction between gRPC endpoints and services
package handlers

import (
	"context"
	"vpeer_file/internal/models"
	"vpeer_file/internal/services"
	pb "vpeer_file/proto/api"
)

// FileHandler handles gRPC requests for file-related operations
type FileHandler struct {
	pb.UnimplementedFileServiceServer
	fileService services.FileService
}

// NewFileHandler creates a new instance of FileHandler
func NewFileHandler(fileService services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// AddFile handles the gRPC request to add a file
func (c *FileHandler) AddFile(ctx context.Context, request *pb.AddFileRequest) (*pb.AddFileResponse, error) {
	// AddFile method comments are removed for brevity
	data := models.File{
		FileId:   request.FileId,
		FileType: request.FileType,
		RoomId:   request.RoomId,
		UserId:   request.UserId,
	}
	res, err := c.fileService.AddFile(&data)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return &pb.AddFileResponse{FileId: res}, nil
}

// RemoveFile handles the gRPC request to remove a file
func (c *FileHandler) RemoveFile(ctx context.Context, request *pb.RemoveFileRequest) (*pb.Empty, error) {
	err := c.fileService.RemoveFile(request.FileId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}
	return &pb.Empty{}, nil
}

// GetFileByFileid handles the gRPC request to get files by their ID
func (c *FileHandler) GetFileByFileid(ctx context.Context, request *pb.GetFileByFileIdRequest) (*pb.Files, error) {
	res, err := c.fileService.GetFileByFileId(request.FileId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return filesToProtoList(&models.Files{Files: res}), nil
}

// GetFileByUserid handles the gRPC request to get files by user ID
func (c *FileHandler) GetFileByUserid(ctx context.Context, request *pb.GetFileByUserIdRequest) (*pb.Files, error) {
	res, err := c.fileService.GetFileByUserId(request.UserId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return filesToProtoList(&models.Files{Files: res}), nil
}

// GetFileByRoomid handles the gRPC request to get files by room ID
func (c *FileHandler) GetFileByRoomid(ctx context.Context, request *pb.GetFileByRoomIdRequest) (*pb.Files, error) {
	res, err := c.fileService.GetFileByRoomId(request.RoomId)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return filesToProtoList(&models.Files{Files: res}), nil
}

func (c *FileHandler) GetFiles(ctx context.Context, request *pb.Pagination) (*pb.Files, error) {
	data := models.Pagination{
		Offset:   request.Offset,
		Limit:    request.Limit,
		GetTotal: request.GetTotal,
	}
	res, err := c.fileService.GetFiles(&data)
	if err != nil {
		return nil, err.ErrorToGRPCStatus()
	}

	return filesToProtoList(&models.Files{Files: res.Files, Total: res.Total}), nil
}

// fileToProto converts a models.File object to its protobuf representation
func fileToProto(up *models.File) *pb.File {
	// Function converts a single file to its protobuf equivalent
	return &pb.File{
		FileId:    up.FileId,
		FileType:  up.FileType,
		RoomId:    up.RoomId,
		UserId:    up.UserId,
		CreatedAt: up.CreatedAt.Unix(),
	}
}

// filesToProtoList converts a models.Files object to a protobuf list
func filesToProtoList(files *models.Files) *pb.Files {
	// Function converts a list of files to their protobuf equivalent list
	pbFiles := make([]*pb.File, 0, len(files.Files))
	for _, file := range files.Files {
		pbFiles = append(pbFiles, fileToProto(&file))
	}

	return &pb.Files{Files: pbFiles, Total: files.Total}
}
