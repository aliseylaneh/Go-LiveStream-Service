package services

import (
	"context"
	"strings"
	"vpeer_usergw/inetrnal/global"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/types"
	pb_file "vpeer_usergw/proto/api/file"
)

type (
	FileService interface {
		AddFile(*models.AddFileDTO) (string, *types.Error)
		RemoveFile(string) *types.Error
		GetFileByFileId(string) ([]models.FileDTO, *types.Error)
		GetFileByUserId(string) ([]models.FileDTO, *types.Error)
		GetFileByRoomId(string) ([]models.FileDTO, *types.Error)
		GetFiles(*models.Pagination) (*models.FilesDTO, *types.Error)
	}
	fileService struct {
	}
)

func NewFileService() FileService {
	return &fileService{}
}

func (c *fileService) AddFile(data *models.AddFileDTO) (string, *types.Error) {
	splitedFileName := strings.Split(data.FileId, ".")
	var fileType string
	if splitedFileName[1] == "ivf" || splitedFileName[1] == "mkv" {
		fileType = "video"
	} else {
		fileType = "audio"
	}
	res, err := global.FILE_SERVER_CLIENT.AddFile(context.Background(), &pb_file.AddFileRequest{
		FileId:   splitedFileName[0],
		FileType: fileType,
		RoomId:   data.RoomId,
		UserId:   data.UserId,
	})

	if err != nil {
		return "", types.ExtractGRPCErrDetails(err)
	}
	return res.FileId, nil
}

func (c *fileService) RemoveFile(fileId string) *types.Error {
	_, err := global.FILE_SERVER_CLIENT.RemoveFile(context.Background(), &pb_file.RemoveFileRequest{
		FileId: fileId,
	})
	if err != nil {
		return types.ExtractGRPCErrDetails(err)
	}
	return nil
}

func (c *fileService) GetFileByFileId(fileId string) ([]models.FileDTO, *types.Error) {
	res, err := global.FILE_SERVER_CLIENT.GetFileByFileid(context.Background(), &pb_file.GetFileByFileIdRequest{
		FileId: fileId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	files := make([]models.FileDTO, 0)
	userIds := make([]string, 0)
	for _, f := range res.Files {
		userIds = append(userIds, *f.UserId)
		file := models.FileDTO{
			FileId:    f.FileId,
			FileType:  f.FileType,
			RoomId:    f.RoomId,
			UserId:    f.UserId,
			CreatedAt: f.CreatedAt,
		}
		files = append(files, file)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	for i, r := range files {
		files[i].User = mapUsers[*r.UserId]
	}

	return files, nil
}

func (c *fileService) GetFileByUserId(userId string) ([]models.FileDTO, *types.Error) {
	res, err := global.FILE_SERVER_CLIENT.GetFileByUserid(context.Background(), &pb_file.GetFileByUserIdRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	files := make([]models.FileDTO, 0)
	userIds := make([]string, 0)
	for _, f := range res.Files {
		userIds = append(userIds, *f.UserId)
		file := models.FileDTO{
			FileId:    f.FileId,
			FileType:  f.FileType,
			RoomId:    f.RoomId,
			UserId:    f.UserId,
			CreatedAt: f.CreatedAt,
		}
		files = append(files, file)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	for i, r := range files {
		files[i].User = mapUsers[*r.UserId]
	}
	return files, nil
}

func (c *fileService) GetFileByRoomId(roomId string) ([]models.FileDTO, *types.Error) {
	res, err := global.FILE_SERVER_CLIENT.GetFileByRoomid(context.Background(), &pb_file.GetFileByRoomIdRequest{
		RoomId: roomId,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	files := make([]models.FileDTO, 0)
	userIds := make([]string, 0)
	for _, f := range res.Files {
		userIds = append(userIds, *f.UserId)
		file := models.FileDTO{
			FileId:    f.FileId,
			FileType:  f.FileType,
			RoomId:    f.RoomId,
			UserId:    f.UserId,
			CreatedAt: f.CreatedAt,
		}
		files = append(files, file)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	for i, r := range files {
		files[i].User = mapUsers[*r.UserId]
	}
	return files, nil
}

func (c *fileService) GetFiles(pagination *models.Pagination) (*models.FilesDTO, *types.Error) {
	res, err := global.FILE_SERVER_CLIENT.GetFiles(context.Background(), &pb_file.Pagination{
		Offset:   pagination.Offset,
		Limit:    pagination.Limit,
		GetTotal: pagination.GetTotal,
	})
	if err != nil {
		return nil, types.ExtractGRPCErrDetails(err)
	}
	if pagination.GetTotal && pagination.Limit == 0 {
		return &models.FilesDTO{Files: []models.FileDTO{}, Total: res.Total}, nil
	}
	files := make([]models.FileDTO, 0)
	userIds := make([]string, 0)
	for _, f := range res.Files {
		userIds = append(userIds, *f.UserId)
		file := models.FileDTO{
			FileId:    f.FileId,
			FileType:  f.FileType,
			RoomId:    f.RoomId,
			UserId:    f.UserId,
			CreatedAt: f.CreatedAt,
		}
		files = append(files, file)
	}

	mapUsers, rerr := models.GetUsersInfoByIds(userIds)
	if rerr != nil {
		return nil, rerr
	}

	for i, r := range files {
		files[i].User = mapUsers[*r.UserId]
	}

	return &models.FilesDTO{Files: files, Total: res.Total}, nil
}
