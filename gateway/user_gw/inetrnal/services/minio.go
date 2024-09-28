package services

import (
	"context"
	"fmt"
	"os"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/types"

	"github.com/minio/minio-go/v7"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

type (
	MinioService interface {
		CreateIvfFile(string) (*ivfwriter.IVFWriter, *types.Error)
		CreateOggFile(string) (*oggwriter.OggWriter, *types.Error)
		PutStorage(string, string, string) *types.Error
		OpenFile(string) (*os.File, *types.Error)
		RemoveFile(string) *types.Error
		RemoveFilesByMinio(string) *types.Error
		GetObject(string) *types.Error
		CheckEntity(string) (bool, error)
		CreateMkvFile(string) (*os.File, *types.Error)
	}
	minioService struct {
		minioClient              *minio.Client
		fileService              FileService
		storagePath              string
		minioDownloadedFilesPath string
	}
)

func NewMinioService(minioClient *minio.Client, fileService FileService, storagePath string, minioDownloadedFilesPath string) MinioService {
	return &minioService{
		minioClient:              minioClient,
		fileService:              fileService,
		storagePath:              storagePath,
		minioDownloadedFilesPath: minioDownloadedFilesPath,
	}
}
func (c *minioService) GetObject(fileName string) *types.Error {
	err := c.minioClient.FGetObject(context.Background(), "recordbucket", fileName, c.minioDownloadedFilesPath+"/"+fileName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("#1")
		return types.NewInternalError("خطا داخلی رخ داده است. کد خطا 86")
	}
	return nil
}

func (c *minioService) OpenFile(fileName string) (*os.File, *types.Error) {
	file, err := os.OpenFile(c.storagePath+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println("#2")
		return nil, types.NewInternalError("خطا داخلی رخ داده است. کد خطا 78")

	}
	return file, nil
}

func (c *minioService) CreateIvfFile(fileName string) (*ivfwriter.IVFWriter, *types.Error) {
	file, err := ivfwriter.New(c.storagePath + "/" + fileName)
	if err != nil {
		fmt.Println("#3")
		fmt.Println(err)
		return nil, types.NewInternalError("خطا داخلی رخ داده است. کد خطا 76")
	}
	return file, nil
}

func (c *minioService) CreateMkvFile(fileName string) (*os.File, *types.Error) {
	file, err := os.Create(c.storagePath + "/" + fileName)
	if err != nil {
		fmt.Println("#4")
		return nil, types.NewInternalError("خطا داخلی رخ داده است. کد خطا 77")
	}
	return file, nil
}

func (c *minioService) CreateOggFile(fileName string) (*oggwriter.OggWriter, *types.Error) {

	file, err := oggwriter.New(c.storagePath+"/"+fileName, 48000, 2)
	if err != nil {
		fmt.Println("#5")
		return nil, types.NewInternalError("خطا داخلی رخ داده است. کد خطا 75")
	}
	return file, nil
}

func (c *minioService) RemoveFile(fileName string) *types.Error {
	err := os.Remove(c.storagePath + "/" + fileName)
	if err != nil {
		return types.NewInternalError("خطا داخلی رخ داده است. کد خطا 79")

	}
	return nil
}

func (c *minioService) RemoveFilesByMinio(fileName string) *types.Error {
	err := c.minioClient.RemoveObject(context.Background(), "recordbucket", fileName+".ivf", minio.RemoveObjectOptions{})
	if err != nil {
		return types.NewInternalError("خطا داخلی رخ داده است. کد خطا 81")
	}
	err = c.minioClient.RemoveObject(context.Background(), "recordbucket", fileName+".ogg", minio.RemoveObjectOptions{})
	if err != nil {
		return types.NewInternalError("خطا داخلی رخ داده است. کد خطا 82")

	}
	return nil
}

func (c *minioService) PutStorage(fileName string, roomId string, userId string) *types.Error {
	_, err := c.minioClient.FPutObject(context.Background(), "recordbucket", fileName, c.storagePath+"/"+fileName, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("#6")
		fmt.Println(err)
		return types.NewInternalError("خطا داخلی رخ داده است. کد خطا 73")
	}
	data := models.AddFileDTO{
		FileId: fileName,
		RoomId: roomId,
		UserId: &userId,
	}
	_, cerr := c.fileService.AddFile(&data)
	if cerr != nil {
		return cerr
	}
	err = os.Remove(c.storagePath + "/" + fileName)
	if err != nil {
		fmt.Println("#7")
		return types.NewInternalError("خطا داخلی رخ داده است. کد خطا 74")

	}
	return nil
}

func (c *minioService) CheckEntity(fileId string) (bool, error) {
	_, err := c.minioClient.GetObjectACL(context.Background(), "recordbucket", fileId)
	if err != nil {
		fmt.Println("#8")
		if err.Error() != "The specified key does not exist. " {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
