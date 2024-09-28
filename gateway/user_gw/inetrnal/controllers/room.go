package controllers

import (
	"safir/libs/idgen"
	"strings"
	"time"
	"vpeer_usergw/inetrnal/global"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/services"

	"github.com/gofiber/fiber/v2"
)

type (
	RoomController interface {
		CreateRoom(*fiber.Ctx) error
		CloseRoom(*fiber.Ctx) error
		GetOpenRoomByUserId(*fiber.Ctx) error
		GetRoomsByUserId(*fiber.Ctx) error
		GetRoomByRoomId(*fiber.Ctx) error
		GetCreatorByRoomId(*fiber.Ctx) error
		JoinRoom(*fiber.Ctx) error
		GetRoomResultsCount(*fiber.Ctx) error
		GetRooms(*fiber.Ctx) error
		GetRoomLogsByRoomId(*fiber.Ctx) error
		GetRoomResultByRoomId(*fiber.Ctx) error
		GetOnGoingRooms(*fiber.Ctx) error
		GetAllUsers(*fiber.Ctx) error
		AddBanUser(*fiber.Ctx) error
		RemoveBanUser(*fiber.Ctx) error
	}
	roomController struct {
		roomService services.RoomService
	}
)

func NewRoomController(roomService services.RoomService) RoomController {
	return &roomController{
		roomService: roomService,
	}
}

func (c *roomController) CreateRoom(ctx *fiber.Ctx) error {
	registerRoomDto := new(models.RegisterRoomDTO)
	if err := ctx.BodyParser(registerRoomDto); err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا در درخواست. کد خطا 58",
			"success": false,
		})
	}
	userId := ctx.Locals("user_id")
	if userId == nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "کاربر پیدا نشد. کد خطا 59",
			"success": false,
		})
	}
	registerRoomDto.UserId = userId.(string)
	res, err := c.roomService.RegisterRoom(registerRoomDto)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})

}

func (c *roomController) CloseRoom(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")
	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 60",
			"success": false,
		})
	}
	LocalData := ctx.Locals("user_id")
	if LocalData == nil {
		return ctx.JSON(map[string]interface{}{
			"message": "کاربر پیدا نشد. کد خطا 61",
			"success": false,
		})
	}

	userId := LocalData.(string)
	res, err := c.roomService.GetCreatorByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	if userId != res {
		return ctx.JSON(map[string]interface{}{
			"message": "شما مجوز بستن جلسه را ندارید. کد خطا 73",
			"success": false,
		})
	}
	err = c.roomService.CloseRoom(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"success": true,
	})
}

func (c *roomController) GetOpenRoomByUserId(ctx *fiber.Ctx) error {
	LocalData := ctx.Locals("user_id")
	if LocalData == nil {
		return ctx.JSON(map[string]interface{}{
			"message": "کاربر پیدا نشد. کد خطا 62",
			"success": false,
		})
	}

	userId := LocalData.(string)
	res, err := c.roomService.GetOpenRoomByUserId(userId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetRoomsByUserId(ctx *fiber.Ctx) error {
	LocalData := ctx.Locals("user_id")
	if LocalData == nil {
		return ctx.JSON(map[string]interface{}{
			"message": "کاربر پیدا نشد. کد خطا 63",
			"success": false,
		})
	}

	userId := LocalData.(string)
	res, err := c.roomService.GetRoomsByUserId(userId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetRoomByRoomId(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")
	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 64",
			"success": false,
		})
	}
	// LocalData := ctx.Locals("user_id")
	// if LocalData == nil {
	// 	return ctx.JSON(map[string]interface{}{
	// 		"message": "کاربر پیدا نشد. کد خطا 65",
	// 		"success": false,
	// 	})
	// }

	// userId := LocalData.(int32)
	res, err := c.roomService.GetRoomByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetCreatorByRoomId(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")
	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 66",
			"success": false,
		})
	}
	// LocalData := ctx.Locals("user_id")
	// if LocalData == nil {
	// 	return ctx.JSON(map[string]interface{}{
	// 		"message": "کاربر پیدا نشد. کد خطا 67",
	// 		"success": false,
	// 	})
	// }

	// userId := LocalData.(int32)
	res, err := c.roomService.GetCreatorByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) JoinRoom(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")

	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 68",
			"success": false,
		})
	}
	roomId = strings.Trim(roomId, "")
	LocalData := ctx.Locals("user_id")
	if LocalData == nil {
		return ctx.JSON(map[string]interface{}{
			"message": "کاربر پیدا نشد. کد خطا 72",
			"success": false,
		})
	}
	for i, c := range global.TOKENS {
		if c.ExpireAt.Before(time.Now()) {
			delete(global.TOKENS, i)
		}
	}

	userIdLocalData := LocalData.(string)
	userId := copyString(userIdLocalData)
	res, err := c.roomService.IsRoomJoinable(roomId, userId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	roomResult, err := c.roomService.GetRoomByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	switch res {
	case "closed":
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "جلسه بسته شده است. کد خطا 69",
			"success": false,
		})
	case "scheduled":
		if roomResult.Creator != userId {
			return ctx.Status(400).JSON(map[string]interface{}{
				"message": "جلسه شروع نشده است. کد خطا 70",
				"success": false,
			})
		}
	case "expired":
		// if roomResult.Creator != userId {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "جلسه به پایان رسیده است. کد خطا 1-70",
			"success": false,
		})
		// }
	case "ban":
		// if roomResult.Creator != userId {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "حساب شما مسدود شده است. کد خطا 2-70",
			"success": false,
		})
		// }
	}
	token, cerr := idgen.NextAlphabeticString(30)
	if cerr != nil {
		return ctx.Status(500).JSON(map[string]interface{}{
			"message": "خطای داخلی رخ داده است. کد خطا 71",
			"success": false,
		})
	}
	room_copy := models.GetRoomCopy()[roomId]
	if room_copy != nil {

		found := false
		for _, conn := range room_copy.Peers.Connections {
			if conn.UserId == roomResult.Creator {
				found = true
				break
			}
		}
		if int32(len(room_copy.Peers.Connections)) == roomResult.UsersLength-1 && !found {
			if userId != roomResult.Creator {
				return ctx.Status(400).JSON(map[string]interface{}{
					"message": "اجازه ورود ندارید",
					"success": false,
				})
			}
		} else if int32(len(room_copy.Peers.Connections)) >= roomResult.UsersLength {
			return ctx.Status(400).JSON(map[string]interface{}{
				"message": "ظرفیت جلسه پر شده است",
				"success": false,
			})
		}
	}

	global.TOKENS[token] = &global.TokenInformation{RoomId: roomId, UserId: userId, ExpireAt: time.Now().Add(time.Hour * 2)}
	return ctx.JSON(map[string]interface{}{
		"data":    token,
		"success": true,
	})
}

func copyString(original string) string {
	// Convert the original string to a byte slice and back to a string.
	copiedString := string([]byte(original))
	return copiedString
}

func (c *roomController) GetRoomResultsCount(ctx *fiber.Ctx) error {
	res, err := c.roomService.GetRoomResultsCount()
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetRooms(ctx *fiber.Ctx) error {
	paginationDto := new(models.Pagination)
	if err := ctx.BodyParser(paginationDto); err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا در درخواست. کد خطا 218",
			"success": false,
		})
	}

	res, err := c.roomService.GetRooms(paginationDto)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetRoomLogsByRoomId(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")
	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 219",
			"success": false,
		})
	}

	res, err := c.roomService.GetRoomLogsByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetRoomResultByRoomId(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")
	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 220",
			"success": false,
		})
	}

	res, err := c.roomService.GetRoomResultByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetOnGoingRooms(ctx *fiber.Ctx) error {
	res, err := c.roomService.GetOnGoingRooms()
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) GetAllUsers(ctx *fiber.Ctx) error {
	paginationDto := new(models.Pagination)
	if err := ctx.BodyParser(paginationDto); err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا در درخواست. کد خطا 221",
			"success": false,
		})
	}

	res, err := c.roomService.GetAllUsers(paginationDto)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *roomController) AddBanUser(ctx *fiber.Ctx) error {
	userIdDTO := new(models.UserIdDTO)
	if err := ctx.BodyParser(userIdDTO); err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا در درخواست. کد خطا 222",
			"success": false,
		})
	}

	err := c.roomService.AddBanUser(userIdDTO.UserId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"message": "کاربر به لیست مسدودی ها اضافه شد",
		"success": true,
	})
}

func (c *roomController) RemoveBanUser(ctx *fiber.Ctx) error {
	userIdDTO := new(models.UserIdDTO)
	if err := ctx.BodyParser(userIdDTO); err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا در درخواست. کد خطا 223",
			"success": false,
		})
	}

	err := c.roomService.RemoveBanUser(userIdDTO.UserId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"message": "کاربر از لیست مسدودی ها خارج شد",
		"success": true,
	})
}
