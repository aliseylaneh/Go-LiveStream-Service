// Package repository provides an interface and implementation for interacting with the database
// to manage room-related data.
package repository

import (
	"database/sql"
	"fmt"
	"time"
	"vpeer_room/internal/models"
	"vpeer_room/internal/types"

	"github.com/lib/pq"
)

// RoomRepository is an interface defining methods for room-related database operations.
type RoomRepository interface {
	// Check if a room with the given room ID exists.
	RoomExistsByRoomId(string) (bool, *types.Error)
	// Register a new room.
	RegisterRoom(*models.Room) *types.Error
	// Add a scheduled event to a room.
	AddScheduled(string, time.Time) *types.Error
	ValidateCreatorScheduledRoom(int32) *types.Error
	CloseRoomByRoomId(string) *types.Error
	GetRoomsByUserId(string) ([]models.Room, *types.Error)
	GetRoomByRoomId(string) (*models.Room, *types.Error)
	GetOpenRoomByUserId(string) (*models.Room, *types.Error)
	GetRoomCreatorByRoomid(string) (string, *types.Error)
	IsRoomJoinable(string, string) (string, *types.Error)
	IsRoomExpired(string) (bool, *types.Error)
	GetOpenRoomsTotalCount() (int32, *types.Error)
	GetRoomsTotalCount() (int32, *types.Error)
	GetRoomLogsByRoomIdTotalCount(string) (int32, *types.Error)
	GetRoomLogsTotalCount() (int32, *types.Error)
	GetRooms(*models.Pagination) ([]models.Room, *types.Error)
	GetOpenRooms(*models.Pagination) ([]models.Room, *types.Error)
	AddRoomLog(*models.RoomLog) *types.Error
	GetRoomLogsByRoomId(string) ([]models.RoomLog, *types.Error)
	GetAllRoomLogs(string, *models.Pagination) ([]models.RoomLog, *types.Error)
	GetAllRoomResults(*models.Pagination) ([]models.RoomResult, *types.Error)
	AddRoomResult(*models.RoomResult) *types.Error
	GetRoomResultTotalCount() (int32, *types.Error)
	GetRoomResultByRoomId(string) (*models.RoomResult, *types.Error)
	GetAllRoomResultsCount() (*models.RoomResultsCount, *types.Error)
	GetAllUsers(*models.Pagination) ([]models.RoomUser, *types.Error)
	GetAllUsersCount() (int32, *types.Error)
	AddBanUser(string) *types.Error
	RemoveBanUser(string) *types.Error
	// GetAllRoomResultsCount() (*models.RoomResultsCount, *types.Error)

	// GetArchivedRoomByRoomId(string) (*models.Room, *types.Error)
}

// roomRepository is an implementation of the RoomRepository interface.
type roomRepository struct {
	db *sql.DB
}

// SQL query to check if a room with the given room ID exists.
const roomExistsByRoomIdQuery = "SELECT 1 FROM rooms WHERE room_id = $1"

// NewRoomRepository creates and returns a new RoomRepository instance.
func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

// RoomExistsByRoomId checks if a room with the given room ID exists in the database.
func (c *roomRepository) RoomExistsByRoomId(roomId string) (bool, *types.Error) {
	var result int32
	err := c.db.QueryRow(roomExistsByRoomIdQuery, roomId).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 10")
	}
	return true, nil
}

// RegisterRoom inserts a new room into the database.
func (c *roomRepository) RegisterRoom(data *models.Room) *types.Error {
	tx, err := c.db.Begin()
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 19")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// var roomID string
	// err = tx.QueryRow("SELECT room_id FROM rooms WHERE creator = $1 AND closed = false", data.Creator).Scan(&roomID)

	// if err == sql.ErrNoRows {
	// } else if err != nil {
	// 	return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 32")
	// }

	// if roomID != "" {
	// 	return types.NewBadRequestError("شما یک جلسه باز دارید. کد خطا 33")
	// }

	_, err = tx.Exec("INSERT INTO rooms (room_id, creator, users_length, closed, closed_at, created_at) VALUES ($1,$2,$3,false,null,NOW())", data.RoomId, data.Creator, data.UsersLength)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 11")
	}
	if data.Scheduled != nil {
		// var roomID string
		// err := c.db.QueryRow(`SELECT s.room_id
		// FROM scheduled AS s
		// JOIN rooms AS r ON s.room_id = r.room_id
		// WHERE r.creator = $1 AND r.closed = false AND s.starts_at > NOW()`, data.Creator).Scan(&roomID)

		// if err == sql.ErrNoRows {
		// 	// No scheduled room found for the creator or all scheduled rooms start in the past, return nil.
		// } else if err != nil {
		// 	// An error occurred during the query.
		// 	return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 20")
		// }

		// if roomID != "" {
		// 	return types.NewBadRequestError("شما یک جلسه زمانبندی شده باز دارید. کد خطا 21")
		// }

		_, err = tx.Exec("INSERT INTO scheduled(room_id,starts_at) VALUES ($1,$2)", data.RoomId, data.Scheduled)
		if err != nil {
			return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 22")
		}
	}
	// expendTime := time.Now().Add(time.Second * time.Duration((data.UsersLength*15)+30))
	// _, err = tx.Exec("INSERT INTO room_expiry(room_id,ends_at) VALUES ($1,$2)", data.RoomId, expendTime)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 2-22")
	// }
	if data.ExpireAt != nil {
		_, err = tx.Exec("INSERT INTO room_expiry(room_id,ends_at) VALUES ($1,$2)", data.RoomId, data.ExpireAt)
		if err != nil {
			fmt.Println(err)
			return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 2-22")
		}
	}
	err = tx.Commit()
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 23")
	}
	return nil
}

// AddScheduled adds a scheduled event to a room.
func (c *roomRepository) AddScheduled(roomId string, starts_at time.Time) *types.Error {
	_, err := c.db.Exec("INSERT INTO scheduled(room_id,starts_at) VALUES ($1,$2)", roomId, starts_at)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 12")
	}
	return nil
}

// ValidateCreatorScheduledRoom checks if a creator already has a scheduled room.
// It returns an error if the creator has a scheduled room that starts in the future, and nil if not.
func (c *roomRepository) ValidateCreatorScheduledRoom(creator int32) *types.Error {

	// Define the SQL query to check if the creator has a scheduled room that starts in the future.
	// Execute the query and check if the result contains any scheduled rooms that start in the future.
	var roomID string
	err := c.db.QueryRow(`SELECT s.room_id
	FROM scheduled AS s
	JOIN rooms AS r ON s.room_id = r.room_id
	WHERE r.creator = $1 AND r.closed = false AND s.starts_at > NOW()`, creator).Scan(&roomID)

	if err == sql.ErrNoRows {
		// No scheduled room found for the creator or all scheduled rooms start in the past, return nil.
		return nil
	} else if err != nil {
		// An error occurred during the query.
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 16")
	}

	// The creator already has a scheduled room that starts in the future, return a validation error.
	return types.NewBadRequestError("شما یک جلسه زمانبندی شده دارید. کد خطا 17")
}

// func (c *roomRepository) CloseRoomByRoomId(roomId string) *types.Error {
// 	// Fetch the room details
// 	room, errc := c.GetRoomByRoomId(roomId)
// 	if errc != nil {
// 		return errc
// 	}

// 	tx, err := c.db.Begin()
// 	if err != nil {
// 		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 102")
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	fmt.Println(roomId)
// 	fmt.Println("1")
// 	// Move the room to archived_rooms table
// 	_, err = tx.Exec("INSERT INTO archived_rooms (room_id, archived_at, creator, users_length, closed, created_at) VALUES ($1, NOW(), $2, $3, true, $4)", room.RoomId, room.Creator, room.UsersLength, room.CreatedAt)
// 	if err != nil {
// 		fmt.Println(err)
// 		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 103")
// 	}
// 	fmt.Println("2")
// 	// Remove corresponding entry from room_expiry table
// 	_, err = tx.Exec("DELETE FROM room_expiry WHERE room_id = $1", roomId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 106")
// 	}
// 	fmt.Println("3")
// 	_, err = tx.Exec("DELETE FROM scheduled WHERE room_id = $1", roomId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 105")
// 	}
// 	fmt.Println("4")
// 	// Remove the room from rooms table
// 	_, err = tx.Exec("DELETE FROM rooms WHERE room_id = $1", roomId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 104")
// 	}

// 	// Remove corresponding entry from scheduled table

// 	err = tx.Commit()
// 	if err != nil {
// 		fmt.Println(err)
// 		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 107")
// 	}

// 	return nil
// }

func (c *roomRepository) CloseRoomByRoomId(roomId string) *types.Error {
	_, err := c.db.Exec("UPDATE rooms SET closed = true, closed_at = NOW() WHERE room_id = $1", roomId)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 18")
	}
	return nil
}

func (c *roomRepository) GetRoomsByUserId(userId string) ([]models.Room, *types.Error) {
	// Query to get room information.
	query := `
        SELECT r.room_id, r.creator, r.users_length, r.closed, r.closed_at, r.created_at, s.starts_at, t.ends_at
        FROM rooms r
        LEFT JOIN scheduled AS s ON r.room_id = s.room_id
		LEFT JOIN room_expiry AS t ON r.room_id = t.room_id
        WHERE r.creator = $1
    `

	rows, err := c.db.Query(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 27")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 24")
	}
	defer rows.Close()
	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(
			&room.RoomId,
			&room.Creator,
			&room.UsersLength,
			&room.Closed,
			&room.ClosedAt,
			&room.CreatedAt,
			&room.Scheduled,
			&room.ExpireAt,
		); err != nil {

			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 25")
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 26")
	}
	return rooms, nil
}

func (c *roomRepository) GetRoomByRoomId(roomId string) (*models.Room, *types.Error) {
	// Query to get room information.
	query := `
        SELECT r.room_id, r.creator, r.users_length, r.closed, r.closed_at, r.created_at, s.starts_at, t.ends_at
        FROM rooms r
        LEFT JOIN scheduled AS s ON r.room_id = s.room_id
		LEFT JOIN room_expiry AS t ON r.room_id = t.room_id
        WHERE r.room_id = $1
    `
	var room models.Room

	err := c.db.QueryRow(query, roomId).Scan(
		&room.RoomId,
		&room.Creator,
		&room.UsersLength,
		&room.Closed,
		&room.ClosedAt,
		&room.CreatedAt,
		&room.Scheduled,
		&room.ExpireAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 28")
		}

		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 29")
	}
	return &room, nil
}

func (c *roomRepository) GetOpenRoomByUserId(userId string) (*models.Room, *types.Error) {
	// Query to get room information.
	query := `
        SELECT r.room_id, r.creator, r.users_length, r.closed, r.closed_at, r.created_at, s.starts_at, t.ends_at
        FROM rooms r
        LEFT JOIN scheduled s ON r.room_id = s.room_id
		LEFT JOIN room_expiry AS t ON r.room_id = t.room_id
        WHERE r.creator = $1 AND closed = false
    `
	var room models.Room

	err := c.db.QueryRow(query, userId).Scan(
		&room.RoomId,
		&room.Creator,
		&room.UsersLength,
		&room.Closed,
		&room.ClosedAt,
		&room.CreatedAt,
		&room.Scheduled,
		&room.ExpireAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه ای پیدا نشد. کد خطا 30")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 31")
	}
	return &room, nil
}

func (c *roomRepository) GetRoomCreatorByRoomid(roomId string) (string, *types.Error) {
	query := `SELECT creator FROM rooms WHERE room_id = $1`
	var result string
	err := c.db.QueryRow(query, roomId).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", types.NewNotFoundError("جلسه پیدا نشد. کد خطا 39")
		}
		return "", types.NewInternalError("خطای داخلی رخ داده است. کد خطا 38")

	}

	return result, nil
}

// IsRoomJoinable checks if a room is joinable, considering closure, scheduled status, and expiry.
func (c *roomRepository) IsRoomJoinable(roomId string, userId string) (string, *types.Error) {
	var status string
	query := `
	SELECT CASE
	WHEN EXISTS (SELECT 1 FROM bans WHERE bans.user_id = $1) THEN 'ban'
	WHEN rm.closed = true THEN 'closed'
	WHEN sc.room_id IS NOT NULL AND sc.starts_at > NOW() THEN 'scheduled'
	WHEN re.ends_at < NOW() THEN 'expired'
	ELSE 'joinable'
	END AS status
	FROM rooms AS rm
	LEFT JOIN scheduled AS sc ON rm.room_id = sc.room_id
	LEFT JOIN room_expiry AS re ON rm.room_id = re.room_id
	WHERE rm.room_id = $2;
	`
	err := c.db.QueryRow(query, userId, roomId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "not_found", types.NewNotFoundError("جلسه پیدا نشد. کد خطا 50")
		}
		return "", types.NewInternalError("خطای داخلی رخ داده است. کد خطا 51")
	}

	return status, nil
}

// IsRoomExpired checks if a room has expired based on room_id.
func (c *roomRepository) IsRoomExpired(roomId string) (bool, *types.Error) {
	var expired bool
	query := `
        SELECT CASE
            WHEN ends_at IS NULL THEN false
            WHEN ends_at < NOW() THEN true
            ELSE false
        END AS is_expired
        FROM room_expiry
        WHERE room_id = $1;
    `
	err := c.db.QueryRow(query, roomId).Scan(&expired)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 50")
		}
		return false, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 101")
	}
	return expired, nil
}

// GetArchivedRoomByRoomId retrieves an archived room by its room ID.
// func (c *roomRepository) GetArchivedRoomByRoomId(roomId string) (*models.Room, *types.Error) {
// 	query := `
//         SELECT room_id, archived_at, creator, users_length, closed, created_at
//         FROM archived_rooms
//         WHERE room_id = $1
//     `
// 	var room models.Room

// 	err := c.db.QueryRow(query, roomId).Scan(
// 		&room.RoomId,
// 		&room.ArchivedAt,
// 		&room.Creator,
// 		&room.UsersLength,
// 		&room.Closed,
// 		&room.CreatedAt,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 108")
// 		}
// 		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 109")
// 	}

// 	return &room, nil
// }

func (c *roomRepository) GetRooms(data *models.Pagination) ([]models.Room, *types.Error) {
	query := `
			SELECT r.room_id, r.creator, r.users_length, r.closed, r.created_at, s.starts_at, t.ends_at
			FROM rooms r
			LEFT JOIN scheduled AS s ON r.room_id = s.room_id
			LEFT JOIN room_expiry AS t ON r.room_id = t.room_id
			ORDER BY r.created_at DESC 
			OFFSET $1 LIMIT $2
		`

	rows, err := c.db.Query(query, data.Offset, data.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 113")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 110")
	}
	defer rows.Close()
	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(
			&room.RoomId,
			&room.Creator,
			&room.UsersLength,
			&room.Closed,
			&room.CreatedAt,
			&room.Scheduled,
			&room.ExpireAt); err != nil {
			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 111")
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 112")
	}
	return rooms, nil
}

func (c *roomRepository) GetOpenRooms(data *models.Pagination) ([]models.Room, *types.Error) {
	query := `
			SELECT r.room_id, r.creator, r.users_length, r.closed, r.created_at, s.starts_at, t.ends_at
			FROM rooms AS r
			LEFT JOIN scheduled AS s ON r.room_id = s.room_id
			LEFT JOIN room_expiry AS t ON r.room_id = t.room_id
			WHERE r.closed = false OFFSET $1 LIMIT $2
		`

	rows, err := c.db.Query(query, data.Offset, data.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 113")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 110")
	}
	defer rows.Close()
	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(
			&room.RoomId,
			&room.Creator,
			&room.UsersLength,
			&room.Closed,
			&room.CreatedAt,
			&room.Scheduled,
			&room.ExpireAt); err != nil {
			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 111")
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {

		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 112")
	}
	return rooms, nil
}

func (c *roomRepository) GetOpenRoomsTotalCount() (int32, *types.Error) {
	var count int32
	err := c.db.QueryRow("SELECT COUNT(*) FROM rooms WHERE closed = false").Scan(&count)
	if err != nil {
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 161")
	}
	return count, nil
}

func (c *roomRepository) GetRoomsTotalCount() (int32, *types.Error) {
	var count int32
	err := c.db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&count)
	if err != nil {
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 167")
	}
	return count, nil
}

// func (c *roomRepository) GetArchivedRooms(data *models.Pagination) ([]models.Room, *types.Error) {
// 	query := `
// 			SELECT room_id, creator, archived_at, user_length, closed, created_at FROM archived_rooms
// 			OFFSET $1 LIMIT $2
// 		`

// 	rows, err := c.db.Query(query, data.Offset, data.Limit)
// 	if err != nil {
// 		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 114")
// 	}
// 	defer rows.Close()
// 	var rooms []models.Room
// 	for rows.Next() {
// 		var room models.Room
// 		if err := rows.Scan(
// 			&room.RoomId,
// 			&room.Creator,
// 			&room.ArchivedAt,
// 			&room.UsersLength,
// 			&room.Closed,
// 			&room.CreatedAt,
// 		); err != nil {

// 			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 115")
// 		}
// 		rooms = append(rooms, room)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 116")
// 	}
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 117")
// 		}
// 	}
// 	return rooms, nil
// }

func (c *roomRepository) AddRoomLog(data *models.RoomLog) *types.Error {
	_, err := c.db.Exec("INSERT INTO room_log(room_id,user_id,user_event,created_at) VALUES ($1,$2,$3,NOW())", data.RoomId, data.UserId, data.UserEvent)
	if err != nil {
		fmt.Println(err)
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 151")
	}
	return nil
}

func (c *roomRepository) GetRoomLogsByRoomId(roomId string) ([]models.RoomLog, *types.Error) {
	query := `
			SELECT room_id, user_id, user_event, created_at FROM room_log WHERE room_id = $1
		`

	rows, err := c.db.Query(query, roomId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 155")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 152")
	}
	defer rows.Close()
	var roomLogs []models.RoomLog
	for rows.Next() {
		var roomLog models.RoomLog
		if err := rows.Scan(
			&roomLog.RoomId,
			&roomLog.UserId,
			&roomLog.UserEvent,
			&roomLog.CreatedAt,
		); err != nil {

			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 153")
		}
		roomLogs = append(roomLogs, roomLog)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 154")
	}
	return roomLogs, nil
}

func (c *roomRepository) GetAllRoomLogs(roomId string, data *models.Pagination) ([]models.RoomLog, *types.Error) {
	query := `
			SELECT room_id, user_id, user_event, created_at FROM room_log WHERE room_id = $1 OFFSET $2 LIMIT $3
		`

	rows, err := c.db.Query(query, roomId, data.Offset, data.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 159")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 156")
	}
	defer rows.Close()
	var roomLogs []models.RoomLog
	for rows.Next() {
		var roomLog models.RoomLog
		if err := rows.Scan(
			&roomLog.RoomId,
			&roomLog.UserId,
			&roomLog.UserEvent,
			&roomLog.CreatedAt,
		); err != nil {

			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 157")
		}
		roomLogs = append(roomLogs, roomLog)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 158")
	}
	return roomLogs, nil
}

func (c *roomRepository) GetRoomLogsByRoomIdTotalCount(roomId string) (int32, *types.Error) {
	var count int32
	err := c.db.QueryRow("SELECT COUNT(*) FROM room_log WHERE room_id = $1", roomId).Scan(&count)
	if err != nil {
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 163")
	}
	return count, nil
}

func (c *roomRepository) GetRoomLogsTotalCount() (int32, *types.Error) {
	var count int32
	err := c.db.QueryRow("SELECT COUNT(*) FROM room_log").Scan(&count)
	if err != nil {
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 165")
	}
	return count, nil
}

func (c *roomRepository) AddRoomResult(data *models.RoomResult) *types.Error {
	_, err := c.db.Exec("INSERT INTO room_result(room_id,approvers,deniers,created_at) VALUES ($1,$2,$3,NOW())", data.RoomId, pq.Array(data.Approvers), pq.Array(data.Deniers))
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 171")
	}
	return nil
}

func (c *roomRepository) GetAllRoomResults(data *models.Pagination) ([]models.RoomResult, *types.Error) {
	query := `
	SELECT room_id, approvers, deniers, created_at 
	FROM room_result 
	ORDER BY created_at DESC 
	OFFSET $1 LIMIT $2;
		`

	rows, err := c.db.Query(query, data.Offset, data.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 175")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 172")
	}
	defer rows.Close()
	var roomResults []models.RoomResult
	for rows.Next() {
		var roomResult models.RoomResult
		if err := rows.Scan(
			&roomResult.RoomId,
			pq.Array(&roomResult.Approvers),
			pq.Array(&roomResult.Deniers),
			&roomResult.CreatedAt,
		); err != nil {
			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 173")
		}
		roomResults = append(roomResults, roomResult)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 174")
	}
	return roomResults, nil
}

func (c *roomRepository) GetRoomResultByRoomId(roomId string) (*models.RoomResult, *types.Error) {
	query := `
			SELECT room_id, approvers, deniers, created_at FROM room_result WHERE room_id = $1
		`

	var room models.RoomResult

	err := c.db.QueryRow(query, roomId).Scan(
		&room.RoomId,
		pq.Array(&room.Approvers),
		pq.Array(&room.Deniers),
		&room.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه ای پیدا نشد. کد خطا 187")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 188")
	}
	return &room, nil
}

func (c *roomRepository) GetRoomResultTotalCount() (int32, *types.Error) {
	var count int32
	err := c.db.QueryRow("SELECT COUNT(*) FROM room_result").Scan(&count)
	if err != nil {
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 165")
	}
	return count, nil
}

func (c *roomRepository) GetAllRoomResultsCount() (*models.RoomResultsCount, *types.Error) {
	query := `SELECT 
    SUM(CASE WHEN cardinality(approvers) > cardinality(deniers) THEN 1 ELSE 0 END) AS success,
    COUNT(*) - SUM(CASE WHEN cardinality(approvers) > cardinality(deniers) THEN 1 ELSE 0 END) AS failed
	FROM 
    room_result;
 	`
	var success int32
	var failed int32
	row := c.db.QueryRow(query)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, types.NewNotFoundError("هیچ نتیجه ای پیدا نشد. کد خطا 210")
		}
		return nil, types.NewInternalError("خطا داخلی رخ داده است. کد خطا 211")
	}
	if err := row.Scan(&success, &failed); err != nil {
		return nil, types.NewInternalError("خطا داخلی رخ داده است. کد خطا 212")
	}
	return &models.RoomResultsCount{Success: success, Failed: failed}, nil
}

func (c *roomRepository) GetAllUsers(data *models.Pagination) ([]models.RoomUser, *types.Error) {
	query := `
	SELECT user_id,
	CASE
		WHEN EXISTS (SELECT 1 FROM bans WHERE bans.user_id = combined_users.user_id) THEN 'ban'
		ELSE 'on_going'
	END AS status,
	MIN(created_at) AS first_occurrence
	FROM (
	 SELECT user_id, created_at FROM room_log
	 UNION ALL
	 SELECT creator AS user_id, created_at FROM rooms
	) AS combined_users
	GROUP BY user_id
	ORDER BY first_occurrence OFFSET $1 LIMIT $2
	`
	rows, err := c.db.Query(query, data.Offset, data.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("جلسه پیدا نشد. کد خطا 301")
		}
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 302")
	}
	defer rows.Close()
	var roomUsers []models.RoomUser
	for rows.Next() {
		var roomUser models.RoomUser
		if err := rows.Scan(
			&roomUser.UserId,
			&roomUser.Status,
			&roomUser.FirstOccurrence,
		); err != nil {
			return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 303")
		}
		roomUsers = append(roomUsers, roomUser)
	}
	if err := rows.Err(); err != nil {
		return nil, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 304")
	}
	return roomUsers, nil

}

func (c *roomRepository) GetAllUsersCount() (int32, *types.Error) {
	var count int32
	query := `
        SELECT COUNT(DISTINCT user_id) AS user_count
        FROM (
            SELECT user_id FROM room_log
            UNION ALL
            SELECT creator AS user_id FROM rooms
        ) AS user_ids` // Give the subquery an alias (user_ids)
	err := c.db.QueryRow(query).Scan(&count)
	fmt.Sprintln("hello")
	if err != nil {
		fmt.Println(err)
		return 0, types.NewInternalError("خطای داخلی رخ داده است. کد خطا 305")
	}
	return count, nil
}

func (c *roomRepository) AddBanUser(userId string) *types.Error {
	// Check if the user is already in the ban list
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS (SELECT 1 FROM bans WHERE user_id = $1)", userId).Scan(&exists)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 306")
	}
	if exists {
		return types.NewBadRequestError("کاربر قبلا به لیست مسدود شده‌ها اضافه شده است. کد خطا 307")
	}

	// Insert the user into the ban list
	_, err = c.db.Exec("INSERT INTO bans(user_id,created_at) VALUES ($1,NOW())", userId)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 308")
	}
	return nil
}
func (c *roomRepository) RemoveBanUser(userId string) *types.Error {
	// Check if the user is in the ban list
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS (SELECT 1 FROM bans WHERE user_id = $1)", userId).Scan(&exists)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 309")
	}
	if !exists {
		return types.NewBadRequestError("کاربر در لیست مسدود شده‌ها وجود ندارد. کد خطا 310")

	}

	// Remove the user from the ban list
	_, err = c.db.Exec("DELETE FROM bans WHERE user_id = $1", userId)
	if err != nil {
		return types.NewInternalError("خطای داخلی رخ داده است. کد خطا 311")
	}
	return nil
}
