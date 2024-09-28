package socket

import (
	"time"
	"vpeer_usergw/inetrnal/models"
)

func (c *socketController) UploadRecords() {
	for range time.NewTicker(time.Second * 2).C {
		for _, record := range models.Recorded {
			if record.MkvFileName != "" {
				err := c.minioService.PutStorage(record.MkvFileName, record.RoomId, record.UserId)
				if err != nil {
					record.WebSocket.WriteJSON(&models.WebSocketMessage{
						Event: "error",
						State: "error",
						Data:  err.Message,
					})
					c.minioService.RemoveFile(record.MkvFileName)
					break
				}
			} else {

				err := c.minioService.PutStorage(record.IvfFileName, record.RoomId, record.UserId)
				if err != nil {
					record.WebSocket.WriteJSON(&models.WebSocketMessage{
						Event: "error",
						State: "error",
						Data:  err.Message,
					})
					c.minioService.RemoveFile(record.IvfFileName)
					c.minioService.RemoveFile(record.OggFileName)
					delete(models.Recorded, record.UserId)
					break
				}

				err = c.minioService.PutStorage(record.OggFileName, record.RoomId, record.UserId)
				if err != nil {
					record.WebSocket.WriteJSON(&models.WebSocketMessage{
						Event: "error",
						State: "error",
						Data:  err.Message,
					})
					c.minioService.RemoveFile(record.OggFileName)
					delete(models.Recorded, record.UserId)
					break
				}
			}
			delete(models.Recorded, record.UserId)
		}
	}
}
