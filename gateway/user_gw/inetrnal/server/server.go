// Package server manages the main server functionality for handling HTTP requests, WebSocket connections,
// and interfacing with other services in the application.
package server

import (
	"fmt"
	"safir/libs/appconfigs"
	"safir/libs/appstates"
	"time"
	"vpeer_usergw/inetrnal/client"
	"vpeer_usergw/inetrnal/controllers"
	"vpeer_usergw/inetrnal/global"
	"vpeer_usergw/inetrnal/helpers"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/routes"
	"vpeer_usergw/inetrnal/services"
	"vpeer_usergw/inetrnal/socket"
	pb_file "vpeer_usergw/proto/api/file"
	pb_room "vpeer_usergw/proto/api/room"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// RunServer sets up and starts the main server functionality.
func RunServer() {
	// Define variables for server configuration obtained from environment variables.
	var (
		listenAddress                  = appconfigs.String("listen-address", "Server listen address")
		roomServerAddress              = appconfigs.String("room-server-address", "Room server address")
		fileServerAddress              = appconfigs.String("file-server-address", "File server address")
		minioServerAddress             = appconfigs.String("minio-server-address", "Minio server address")
		minioAccessKey                 = appconfigs.String("minio-access-key", "Minio access key")
		minioSecretKey                 = appconfigs.String("minio-secret-key", "Minio secret key")
		fileStoragePath                = appconfigs.String("file-storage-path", "Storage Path")
		minioDownloadedFileStoragePath = appconfigs.String("minio-downloaded-file-storage-path", "Storage Path")
	)

	// Handle configuration errors and missing environment parameters.
	if err := appconfigs.Parse(); err != nil {
		appstates.PanicMissingEnvParams(err.Error())
	}

	// Establish connections to external services via gRPC.
	var (
		roomServerConnection = client.GrpcClientServerConnection(*roomServerAddress)
		fileServerConnection = client.GrpcClientServerConnection(*fileServerAddress)
	)

	global.ROOM_SERVER_CLIENT = pb_room.NewRoomServiceClient(roomServerConnection)
	global.FILE_SERVER_CLIENT = pb_file.NewFileServiceClient(fileServerConnection)
	minioServerClient := client.MinioClient(*minioServerAddress, *minioAccessKey, *minioSecretKey)

	// Initialize different services used in the application.
	var (
		fileService  services.FileService  = services.NewFileService()
		roomService  services.RoomService  = services.NewRoomService()
		minioService services.MinioService = services.NewMinioService(minioServerClient, fileService, *fileStoragePath, *minioDownloadedFileStoragePath)

		roomController   controllers.RoomController = controllers.NewRoomController(roomService)
		fileController   controllers.FileController = controllers.NewFileController(fileService, minioService, *minioDownloadedFileStoragePath)
		socketController socket.SocketController    = socket.NewSocketController(minioService, *minioDownloadedFileStoragePath)
	)

	// Create a new Fiber instance.
	app := fiber.New()

	// Use CORS middleware for handling cross-origin requests.
	app.Use(cors.New())

	v1 := app.Group("/api")
	routes.RoomGroup(v1, roomController)
	routes.FileGroup(v1, fileController)

	// Default route for a simple hello world response.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Pong!")
	})

	// Middleware for WebSocket upgrade request.
	app.Use("/room/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint handling.
	app.Get("/room/ws", websocket.New(socketController.WebsocketHandler))

	// Start a goroutine to execute functions periodically.
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			models.DispatchKeyFrames()
		}
	}()

	// go func() {
	// 	for range time.NewTicker(time.Second * 1).C {
	// 		for i, c := range global.TOKENS {
	// 			if c.ExpireAt.Before(time.Now()) {
	// 				delete(global.TOKENS, i)
	// 			}
	// 		}
	// 	}
	// }()

	go func() {
		helpers.CheckExpireRooms()
	}()

	go func() {
		socketController.UploadRecords()
	}()

	// Start the Fiber server and log any errors encountered during startup.
	err := app.Listen(*listenAddress)
	if err != nil {
		fmt.Println(err)
	}
}
