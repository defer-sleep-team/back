package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// UploadImage - функция для загрузки изображения в облако
func UploadImage(file multipart.File, filename string) (string, error) {
	url := "http://45.156.21.196:8765/api_s3_aether_amazon/upload"
	log.Print("Image upload started")

	// Создаем буфер для хранения тела запроса
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Создаем часть для файла с оригинальным именем
	log.Print("Z1")
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("error creating form file: %v", err)
	}

	log.Print("Z2")
	// Копируем содержимое файла в часть запроса
	_, err = io.Copy(part, file)
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("error copying file content: %v", err)
	}
	log.Print("Z3")

	// Закрываем writer
	err = writer.Close()
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	log.Print("Z4")
	// Создаем POST запрос
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("error creating request: %v", err)
	}
	log.Print("Z5")

	// Устанавливаем правильный Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	log.Print("Z6")
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		responseBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Print(readErr.Error())
			log.Printf("Failed to read response body: %v", readErr)
			return "", fmt.Errorf("unexpected status: %s", resp.Status)
		}

		log.Printf("Unexpected status: %s, Response body: %s, Request URL: %s", resp.Status, string(responseBody), resp.Request.URL.String())
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return string(responseBody), nil
}
