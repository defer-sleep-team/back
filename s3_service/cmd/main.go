package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	s3_service "github.com/defer-sleep-team/Aether_backend/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func main() {
	// Определяем флаги командной строки
	bucketName := flag.String("bucket", "aether", "S3 bucket name")
	region := flag.String("region", "ru-central1", "Region")
	flag.Parse()

	// Инициализируем клиент S3
	client, err := s3_service.NewS3Client(*region)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	cloud := s3_service.Cloud{S3: client}

	// Создаем новый Fiber приложение
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // Установите лимит на 10 МБ
	})
	// Эндпоинт для загрузки файла

	app.Post("/api_s3_aether_amazon/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			log.Print(err)
			return c.Status(http.StatusBadRequest).SendString("Failed to get pic")
		}
		log.Print("Started upload")
		// Generate a new UUID
		newUUID := uuid.New().String()

		// Get the file extension
		fileExt := filepath.Ext(file.Filename)
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" && fileExt != ".webp" {
			return c.Status(http.StatusBadRequest).SendString("Invalid file extension")
		}
		// Create a new filename with UUID and original extension
		newFilename := newUUID + fileExt

		tempFile, err := os.CreateTemp("", newFilename)
		if err != nil {
			log.Print(err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to create temp file")
		}
		defer os.Remove(tempFile.Name())

		if err := c.SaveFile(file, tempFile.Name()); err != nil {
			log.Print(err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to save file")
		}

		// Upload file to S3 with the new filename
		err = cloud.UploadFile(*bucketName, newFilename, tempFile.Name())
		if err != nil {
			log.Print(err)
			return c.Status(http.StatusInternalServerError).SendString("Failed to upload file")
		}

		return c.SendString(newFilename)
	})
	// Эндпоинт для удаления файла
	app.Delete("/api_s3_aether_amazon/delete/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		deleted, err := cloud.DeleteObject(c.Context(), *bucketName, key, "", false)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to delete file")
		}
		if !deleted {
			return c.Status(http.StatusNotFound).SendString("File not found")
		}
		return c.SendString("File deleted successfully")
	})

	// Эндпоинт для получения ссылки на файл
	app.Get("/api_s3_aether_amazon/get_url/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		signedURL, err := cloud.GetSignedObjectURL(*bucketName, key)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to get signed URL")
		}
		return c.JSON(fiber.Map{"url": signedURL})
	})
	// Запускаем сервер
	log.Fatal(app.Listen(":8765"))
}
