// Package repository handles file-related database operations
package repository

import (
	"database/sql"
	"fmt"
	"vpeer_file/internal/models"
	"vpeer_file/internal/types"
)

// FileRepository defines the interface for file-related operations
type FileRepository interface {
	// FileExistsByFileIdAndType checks if a file with given ID and type exists
	FileExistsByFileIdAndType(string, string) (bool, *types.Error)
	// FileExistsByFileId checks if a file with given ID exists
	FileExistsByFileId(string) (bool, *types.Error)
	// AddFile adds a new file to the repository
	AddFile(*models.File) *types.Error
	// RemoveFile removes a file by its ID
	RemoveFile(string) *types.Error
	// GetFileByFileId retrieves files by their ID
	GetFileByFileId(string) ([]models.File, *types.Error)
	// GetFileByUserId retrieves files by user ID
	GetFileByUserId(string) ([]models.File, *types.Error)
	// GetFileByRoomId retrieves files by room ID
	GetFileByRoomId(string) ([]models.File, *types.Error)

	GetFiles(*models.Pagination) ([]models.File, *types.Error)

	GetFilesCount() (int32, *types.Error)
}

// fileRepository implements the FileRepository interface
type fileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new instance of fileRepository
func NewFileRepository(db *sql.DB) FileRepository {
	return &fileRepository{
		db: db,
	}
}

// FileExistsByFileIdAndType checks if a file with given ID and type exists
func (c *fileRepository) FileExistsByFileIdAndType(fileId string, fileType string) (bool, *types.Error) {
	var result int32
	err := c.db.QueryRow("SELECT 1 FROM files WHERE file_id = $1 AND file_type = $2", fileId, fileType).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, types.NewInternalError("Internal error occurred. Error code 34")
	}
	return true, nil
}

// FileExistsByFileId checks if a file with given ID exists
func (c *fileRepository) FileExistsByFileId(fileId string) (bool, *types.Error) {
	var result int32
	err := c.db.QueryRow("SELECT 1 FROM files WHERE room_id = $1", fileId).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, types.NewInternalError("Internal error occurred. Error code 35")
	}
	return true, nil
}

// AddFile adds a new file to the repository
func (c *fileRepository) AddFile(data *models.File) *types.Error {
	_, err := c.db.Exec("INSERT INTO files (file_id, room_id, file_type, user_id, created_at) VALUES ($1,$2,$3,$4,NOW())", data.FileId, data.RoomId, data.FileType, data.UserId)
	if err != nil {
		fmt.Println(err)
		return types.NewInternalError("Internal error occurred. Error code 36")
	}
	return nil
}

// RemoveFile removes a file by its ID
func (c *fileRepository) RemoveFile(fileId string) *types.Error {
	_, err := c.db.Exec("DELETE FROM files WHERE file_id = $1", fileId)
	if err != nil {
		return types.NewInternalError("Internal error occurred. Error code 40")
	}
	return nil
}

// GetFileByFileId retrieves files by their ID
func (c *fileRepository) GetFileByFileId(fileId string) ([]models.File, *types.Error) {
	rows, err := c.db.Query("SELECT file_id, file_type, room_id, user_id, created_at FROM files WHERE file_id = $1", fileId)
	if err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 41")
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File
		if err := rows.Scan(
			&file.FileId,
			&file.FileType,
			&file.RoomId,
			&file.UserId,
			&file.CreatedAt,
		); err != nil {
			return nil, types.NewInternalError("Internal error occurred. Error code 42")
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 43")
	}
	return files, nil
}

// GetFileByUserId retrieves files by user ID
func (c *fileRepository) GetFileByUserId(userId string) ([]models.File, *types.Error) {
	rows, err := c.db.Query("SELECT file_id, file_type, room_id, user_id, created_at FROM files WHERE user_id = $1", userId)
	if err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 44")
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File
		if err := rows.Scan(
			&file.FileId,
			&file.FileType,
			&file.RoomId,
			&file.UserId,
			&file.CreatedAt,
		); err != nil {
			return nil, types.NewInternalError("Internal error occurred. Error code 45")
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 46")
	}
	return files, nil
}

// GetFileByRoomId retrieves files by room ID
func (c *fileRepository) GetFileByRoomId(roomId string) ([]models.File, *types.Error) {
	rows, err := c.db.Query("SELECT file_id, file_type, room_id, user_id, created_at FROM files WHERE room_id = $1", roomId)
	if err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 47")
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File
		if err := rows.Scan(
			&file.FileId,
			&file.FileType,
			&file.RoomId,
			&file.UserId,
			&file.CreatedAt,
		); err != nil {
			return nil, types.NewInternalError("Internal error occurred. Error code 48")
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 49")
	}
	return files, nil
}

func (c *fileRepository) GetFiles(pagination *models.Pagination) ([]models.File, *types.Error) {
	rows, err := c.db.Query("SELECT file_id, file_type, room_id, user_id, created_at FROM files ORDER BY created_at DESC OFFSET $1 LIMIT $2", pagination.Offset, pagination.Limit)
	if err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 47")
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File
		if err := rows.Scan(
			&file.FileId,
			&file.FileType,
			&file.RoomId,
			&file.UserId,
			&file.CreatedAt,
		); err != nil {
			return nil, types.NewInternalError("Internal error occurred. Error code 48")
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("Internal error occurred. Error code 49")
	}
	return files, nil
}

func (c *fileRepository) GetFilesCount() (int32, *types.Error) {
	var count int32
	err := c.db.QueryRow("SELECT COUNT(*) FROM files").Scan(&count)
	if err != nil {
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 50")
	}
	return count, nil
}
