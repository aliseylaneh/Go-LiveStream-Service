// Package services manages file-related operations and interactions
package services

import (
	"vpeer_file/internal/models"
	"vpeer_file/internal/repository"
	"vpeer_file/internal/types"
)

// FileService defines the interface for file-related services
type FileService interface {
	// AddFile adds a new file and returns its ID or an error
	AddFile(*models.File) (string, *types.Error)
	// RemoveFile removes a file by its ID or returns an error
	RemoveFile(string) *types.Error
	// GetFileByFileId retrieves files by their ID or returns an error
	GetFileByFileId(string) ([]models.File, *types.Error)
	// GetFileByUserId retrieves files by user ID or returns an error
	GetFileByUserId(string) ([]models.File, *types.Error)
	// GetFileByRoomId retrieves files by room ID or returns an error
	GetFileByRoomId(string) ([]models.File, *types.Error)

	GetFiles(*models.Pagination) (*models.Files, *types.Error)
}

// fileService implements the FileService interface
type fileService struct {
	repository repository.FileRepository
}

// NewFileService creates a new instance of fileService
func NewFileService(repository repository.FileRepository) FileService {
	return &fileService{
		repository: repository,
	}
}

// AddFile adds a new file to the repository
func (c *fileService) AddFile(data *models.File) (string, *types.Error) {
	// AddFile method comments are removed for brevity
	exists, err := c.repository.FileExistsByFileIdAndType(data.FileId, data.FileType)
	if err != nil {
		return "", err
	}
	if exists {
		return "", types.NewBadRequestError("File already exists. Error code 50")
	}
	err = c.repository.AddFile(data)
	if err != nil {
		return "", err
	}
	return data.FileId, nil
}

// RemoveFile removes a file by its ID
func (c *fileService) RemoveFile(roomId string) *types.Error {
	return c.repository.RemoveFile(roomId)
}

// GetFileByFileId retrieves files by their ID
func (c *fileService) GetFileByFileId(fileId string) ([]models.File, *types.Error) {
	return c.repository.GetFileByFileId(fileId)
}

// GetFileByUserId retrieves files by user ID
func (c *fileService) GetFileByUserId(userId string) ([]models.File, *types.Error) {
	return c.repository.GetFileByUserId(userId)
}

// GetFileByRoomId retrieves files by room ID
func (c *fileService) GetFileByRoomId(roomId string) ([]models.File, *types.Error) {
	return c.repository.GetFileByRoomId(roomId)
}

func (c *fileService) GetFiles(pagination *models.Pagination) (*models.Files, *types.Error) {
	res, err := c.repository.GetFiles(pagination)
	if err != nil {
		return nil, err
	}

	var count *int32
	if pagination.GetTotal {
		res, err := c.repository.GetFilesCount()
		if err != nil {
			return nil, err
		}
		count = &res
	}
	return &models.Files{Files: res, Total: count}, nil
}
