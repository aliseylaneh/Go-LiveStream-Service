// Package routes manages the routing configurations for handling different API endpoints.
package routes

import (
	"vpeer_usergw/inetrnal/controllers"
	"vpeer_usergw/inetrnal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RoomGroup defines routes related to room operations.
func RoomGroup(app fiber.Router, roomController controllers.RoomController) {
	// Create a new group for room-related routes.
	roomGroup := app.Group("/rm")
	roomGroup.Use(middleware.TokenAuthentication)
	// Define various HTTP routes for room operations.
	roomGroup.Post("/room/register", roomController.CreateRoom)          // Route to create a room
	roomGroup.Put("/room/close/:room_id", roomController.CloseRoom)      // Route to close a room by ID
	roomGroup.Get("/room/open", roomController.GetOpenRoomByUserId)      // Route to get open rooms by user ID
	roomGroup.Get("/room/u/all", roomController.GetRoomsByUserId)        // Route to get all rooms by user ID
	roomGroup.Get("/room/r/:room_id", roomController.GetRoomByRoomId)    // Route to get a room by its ID
	roomGroup.Get("/room/c/:room_id", roomController.GetCreatorByRoomId) // Route to get the creator of a room by its ID
	roomGroup.Get("/room/join/:room_id", roomController.JoinRoom)        // Route to join a room by its ID

	adminRoomGroup := app.Group("/cmr/rm")
	adminRoomGroup.Get("/room/result/count", roomController.GetRoomResultsCount)        // Route to get room results count
	adminRoomGroup.Post("/rooms", roomController.GetRooms)                              // Route to get rooms
	adminRoomGroup.Get("/room/log/:room_id", roomController.GetRoomLogsByRoomId)        // Route to get room logs by room ID
	adminRoomGroup.Get("/room/result/r/:room_id", roomController.GetRoomResultByRoomId) // Route to get room result by room ID
	adminRoomGroup.Get("/room/going", roomController.GetOnGoingRooms)                   // Route to get on going rooms
	adminRoomGroup.Post("/room/user/all", roomController.GetAllUsers)                   // Route to get all users
	adminRoomGroup.Get("/room/r/:room_id", roomController.GetRoomByRoomId)              // Route to get a room by its ID
	adminRoomGroup.Post("/room/ban/add", roomController.AddBanUser)                     // Route to add to the banlist
	adminRoomGroup.Delete("/room/ban/remove", roomController.RemoveBanUser)             // Route to remove from the banlist
}
