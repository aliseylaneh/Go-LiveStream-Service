// Package routes manages the routing configurations for handling different API endpoints.
package routes

import (
	"vpeer_usergw/inetrnal/controllers"
	"vpeer_usergw/inetrnal/middleware"

	"github.com/gofiber/fiber/v2"
)

// FileGroup defines routes related to file operations.
func FileGroup(app fiber.Router, fileController controllers.FileController) {
	// Create a new group for file-related routes.
	fileGroup := app.Group("/fl")
	fileGroup.Use(middleware.TokenAuthentication)
	// Define various HTTP routes for file operations.
	fileGroup.Delete("/file/remove/:file_id", fileController.RemoveFile)   // Route to remove a file by ID
	fileGroup.Get("/file/u", fileController.GetFilesInfoByUserId)          // Route to get file information by user ID
	fileGroup.Get("/file/r/:room_id", fileController.GetFilesInfoByRoomId) // Route to get file information by room ID
	fileGroup.Get("/file/f/:file_id", fileController.GetFilesInfoByFileId) // Route to get file information by file ID
	fileGroup.Get("/file/download/:file_id", fileController.DownloadFile)  // Route to download a file by its ID

	adminFileGroup := app.Group("/cmf/fl")

	adminFileGroup.Get("/file/dir/files", fileController.DirFiles)                   // Route to get files from a directory
	adminFileGroup.Get("/file/r/:room_id", fileController.GetFilesInfoByRoomId)      // Route to get file information by room ID
	adminFileGroup.Delete("/file/dir/remove/:file_id", fileController.DirRemoveFile) // Route to remove a file from a directory by ID
	adminFileGroup.Delete("/file/remove/:file_id", fileController.RemoveFile)        // Route to remove a file by ID
	adminFileGroup.Get("/file/download/:file_id", fileController.DownloadFile)       // Route to download a file by its ID
	adminFileGroup.Post("/files", fileController.GetFiles)                           // Route to get files
}
