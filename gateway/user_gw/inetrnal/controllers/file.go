package controllers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/services"

	"github.com/gofiber/fiber/v2"
)

type (
	FileController interface {
		RemoveFile(*fiber.Ctx) error
		GetFilesInfoByUserId(*fiber.Ctx) error
		GetFilesInfoByRoomId(*fiber.Ctx) error
		GetFilesInfoByFileId(*fiber.Ctx) error
		DownloadFile(*fiber.Ctx) error
		DirFiles(*fiber.Ctx) error
		DirRemoveFile(*fiber.Ctx) error
		GetFiles(*fiber.Ctx) error
	}
	fileController struct {
		fileService              services.FileService
		minioService             services.MinioService
		minioDownloadedFilesPath string
	}
)

func NewFileController(fileService services.FileService, minioService services.MinioService, minioDownloadedFilesPath string) FileController {
	return &fileController{
		fileService:              fileService,
		minioService:             minioService,
		minioDownloadedFilesPath: minioDownloadedFilesPath,
	}
}

func (c *fileController) RemoveFile(ctx *fiber.Ctx) error {
	fileId := ctx.Params("file_id")
	if fileId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس فایل نباید خالی باشد. کد خطا 80",
			"success": false,
		})
	}
	err := c.fileService.RemoveFile(fileId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	err = c.minioService.RemoveFilesByMinio(fileId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    "",
		"success": true,
	})
}

func (c *fileController) GetFilesInfoByUserId(ctx *fiber.Ctx) error {
	LocalData := ctx.Locals("user_id")
	if LocalData == nil {
		return ctx.JSON(map[string]interface{}{
			"message": "کاربر پیدا نشد. کد خطا 83",
			"success": false,
		})
	}
	userId := LocalData.(string)
	res, err := c.fileService.GetFileByUserId(userId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *fileController) GetFilesInfoByRoomId(ctx *fiber.Ctx) error {
	roomId := ctx.Params("room_id")
	if roomId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 86",
			"success": false,
		})
	}
	res, err := c.fileService.GetFileByRoomId(roomId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *fileController) GetFilesInfoByFileId(ctx *fiber.Ctx) error {
	fileId := ctx.Params("file_id")
	if fileId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس جلسه نباید خالی باشد. کد خطا 85",
			"success": false,
		})
	}
	res, err := c.fileService.GetFileByFileId(fileId)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}
	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}

func (c *fileController) DownloadFile(ctx *fiber.Ctx) error {
	fileId := ctx.Params("file_id")
	if fileId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس فایل نباید خالی باشد. کد خطا 87",
			"success": false,
		})
	}
	files, err := os.ReadDir(c.minioDownloadedFilesPath)
	if err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطای داخلی رخ داده است. کد خطا 98",
			"success": false,
		})
	}
	fullFilePath := fmt.Sprintf(c.minioDownloadedFilesPath+"/"+"%s.mkv", fileId)
	for _, f := range files {
		if f.Name() == fileId+".mkv" {
			file, err := os.Open(fullFilePath)
			if err != nil {
				return ctx.Status(500).JSON(map[string]interface{}{
					"message": "خطا در خواندن فایل. کد خطا 89",
					"success": false,
				})
			}
			fileInfo, err := file.Stat()
			if err != nil {
				return ctx.Status(500).JSON(map[string]interface{}{
					"message": "خطا در خواندن فایل. کد خطا 90",
					"success": false,
				})
			}
			// Create a reader from fileBytes
			fileReader := bufio.NewReader(file)

			// Use the reader in ctx.SendStream
			ctx.SendStream(fileReader, int(fileInfo.Size()))
			formats := []string{".ivf", ".ogg"}
			for _, f := range formats {
				os.Remove(c.minioDownloadedFilesPath + "/" + fileId + f)
			}
			return nil
		}

	}

	exists, err := c.minioService.CheckEntity(fileId + ".mkv")
	if err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا داخلی رخ داده است. کد خطا 2-92",
			"success": false,
		})
	}
	if exists {
		cerr := c.minioService.GetObject(fileId + ".mkv")
		if cerr != nil {
			return ctx.Status(cerr.ErrorToHttpStatus()).JSON(cerr.ErrorToJsonMessage())
		}
		file, err := os.Open(fullFilePath)
		if err != nil {
			return ctx.Status(500).JSON(map[string]interface{}{
				"message": "خطا در خواندن فایل. کد خطا 89",
				"success": false,
			})
		}
		fileInfo, err := file.Stat()
		if err != nil {
			return ctx.Status(500).JSON(map[string]interface{}{
				"message": "خطا در خواندن فایل. کد خطا 90",
				"success": false,
			})
		}
		// Create a reader from fileBytes
		fileReader := bufio.NewReader(file)

		// Use the reader in ctx.SendStream
		ctx.SendStream(fileReader, int(fileInfo.Size()))
		return nil

	}

	exists, err = c.minioService.CheckEntity(fileId + ".ivf")
	if err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا داخلی رخ داده است. کد خطا 92",
			"success": false,
		})
	}
	if !exists {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "فایل پیدا نشد. کد خطا 93",
			"success": false,
		})
	}

	exists, err = c.minioService.CheckEntity(fileId + ".ogg")
	if err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا داخلی رخ داده است. کد خطا 94",
			"success": false,
		})
	}
	if !exists {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "فایل پیدا نشد. کد خطا 95",
			"success": false,
		})
	}
	cerr := c.minioService.GetObject(fileId + ".ivf")
	if cerr != nil {
		return ctx.Status(cerr.ErrorToHttpStatus()).JSON(cerr.ErrorToJsonMessage())
	}
	cerr = c.minioService.GetObject(fileId + ".ogg")
	if cerr != nil {
		return ctx.Status(cerr.ErrorToHttpStatus()).JSON(cerr.ErrorToJsonMessage())
	}

	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-y",
		// "-f", "video",
		// "-pix_fmt", "rgb24",
		// "-pixel_format", "yuv420p",
		// "-video_size", "200x200",
		// "-framerate", fmt.Sprint(60),
		"-i", fmt.Sprintf(c.minioDownloadedFilesPath+"/"+fileId+".ivf"),
		"-i", fmt.Sprintf(c.minioDownloadedFilesPath+"/"+fileId+".ogg"),
		"-c:v", "copy",
		"-c:a", "copy",
		// "-b:v", "1M",
		// "-vf", "hqdn3d=3:3:6:6",
		// "-preset", "medium",
		// "-crf", "23",
		// "-strict", "-2",
		// "-profile:v", "baseline", // baseline profile is compatible with most devices
		// "-level", "3.0",
		// "-start_number", "0", // start numbering segments from 0
		// "-hls_time", fmt.Sprint(10), // duration of each segment in seconds
		// "-hls_list_size", "0", // keep all segments in the playlist
		// "-f", "hls",
		fmt.Sprintf(c.minioDownloadedFilesPath+"/"+"%s.mkv", fileId),
	)
	// ffmpegCmd := exec.Command(
	// 	"ffmpeg",
	// 	"-y", // Overwrite output files without asking
	// 	"-i", fmt.Sprintf("%s/%s.ivf", c.minioDownloadedFilesPath, fileId),
	// 	"-i", fmt.Sprintf("%s/%s.ogg", c.minioDownloadedFilesPath, fileId),
	// 	"-c:v", "libx264", // Video codec for H.264 encoding
	// 	"-c:a", "aac", // Audio codec for AAC encoding
	// 	"-crf", "18", // Constant Rate Factor for video quality (lower is better quality)
	// 	// "-preset", "fast", // Preset for encoding speed vs. compression ratio
	// 	"-strict", "-2", // Enable experimental codecs
	// 	fmt.Sprintf("%s/%s.mkv", c.minioDownloadedFilesPath, fileId),
	// )

	_, err = ffmpegCmd.CombinedOutput()
	if err != nil {
		return ctx.Status(500).JSON(map[string]interface{}{
			"message": "خطا در خواندن فایل. کد خطا 91",
			"success": false,
		})
	}
	file, err := os.Open(fullFilePath)
	if err != nil {
		return ctx.Status(500).JSON(map[string]interface{}{
			"message": "خطا در خواندن فایل. کد خطا 89",
			"success": false,
		})
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return ctx.Status(500).JSON(map[string]interface{}{
			"message": "خطا در خواندن فایل. کد خطا 90",
			"success": false,
		})
	}
	// Create a reader from fileBytes
	fileReader := bufio.NewReader(file)

	// Use the reader in ctx.SendStream
	ctx.SendStream(fileReader, int(fileInfo.Size()))
	formats := []string{".ivf", ".ogg"}
	for _, f := range formats {
		os.Remove(c.minioDownloadedFilesPath + "/" + fileId + f)
	}
	return nil
}

func (c *fileController) DirFiles(ctx *fiber.Ctx) error {
	files, err := os.ReadDir(c.minioDownloadedFilesPath)
	if err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطای داخلی رخ داده است. کد خطا 98",
			"success": false,
		})
	}
	filesNames := make([]string, 0)
	for _, f := range files {
		filesNames = append(filesNames, f.Name())
	}

	return ctx.Status(200).JSON(map[string]interface{}{
		"data":    filesNames,
		"success": true,
	})

}

func (c *fileController) DirRemoveFile(ctx *fiber.Ctx) error {
	fileId := ctx.Params("file_id")
	if fileId == "" {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "آدرس فایل نباید خالی باشد. کد خطا 99",
			"success": true,
		})
	}
	os.Remove(c.minioDownloadedFilesPath + "/" + fileId + ".mkv")

	return ctx.Status(200).JSON(map[string]interface{}{
		"success": true,
	})
}

func (c *fileController) GetFiles(ctx *fiber.Ctx) error {
	paginationDto := new(models.Pagination)
	if err := ctx.BodyParser(paginationDto); err != nil {
		return ctx.Status(400).JSON(map[string]interface{}{
			"message": "خطا در درخواست. کد خطا 217",
			"success": false,
		})
	}

	res, err := c.fileService.GetFiles(paginationDto)
	if err != nil {
		return ctx.Status(err.ErrorToHttpStatus()).JSON(err.ErrorToJsonMessage())
	}

	return ctx.JSON(map[string]interface{}{
		"data":    res,
		"success": true,
	})
}
