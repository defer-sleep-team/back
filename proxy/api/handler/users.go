package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

type User struct {
	ID               int    `json:"id"`
	Yid              string `json:"yid"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Avatar           string `json:"avatar"`
	Bio              string `json:"bio"`
	PrivilegeLevel   int    `json:"privilege_level"`
	Payments         string `json:"payments"`
	IsBlock          bool   `json:"block"`
	Background       string `json:"background"`
	SubscribersCount int    `json:"subscribers_count"`
	IsSubscribed     bool   `json:"is_subscribed"`
}
type UserFollower struct {
	FollowerID int `json:"follower_id"`
	FolloweeID int `json:"followee_id"`
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	fromUserID, ok := r.Context().Value("userID").(int)
	if !ok {
		fromUserID = 0
	}
	userID := chi.URLParam(r, "userID")
	url := fmt.Sprintf("http://%s:%s/database_zov_russ_cbo/users/get/%d/%s", IP, UsersPort, fromUserID, userID)
	log.Print(url)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	// Создаем WaitGroup для ожидания завершения всех запросов
	var wg sync.WaitGroup
	imageURLChan := make(chan string, 1)      // Для одного URL аватара
	backgroundURLChan := make(chan string, 1) // Для одного URL фона

	wg.Add(2) // Увеличиваем счетчик на 2

	// Получение URL аватара
	go func(url string) {
		defer wg.Done()
		fullURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + user.Avatar
		resp, err := http.Get(fullURL)
		if err != nil {
			log.Print(err)
			imageURLChan <- user.Avatar // Если ошибка, отправляем оригинальный URL
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Print(resp.StatusCode)
			imageURLChan <- user.Avatar // Если не найден, отправляем оригинальный URL
			return
		}

		var result struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Print(err)
			imageURLChan <- user.Avatar // Если ошибка, отправляем оригинальный URL
			return
		}

		imageURLChan <- result.URL
	}(user.Avatar)

	// Получение URL фона
	go func(url string) {
		defer wg.Done()
		fullURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + user.Background
		resp, err := http.Get(fullURL)
		if err != nil {
			log.Print(err)
			backgroundURLChan <- user.Background // Если ошибка, отправляем оригинальный URL
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Print(resp.StatusCode)
			backgroundURLChan <- user.Background // Если не найден, отправляем оригинальный URL
			return
		}

		var result struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Print(err)
			backgroundURLChan <- user.Background // Если ошибка, отправляем оригинальный URL
			return
		}

		backgroundURLChan <- result.URL
	}(user.Background)

	// Закрываем каналы после завершения всех запросов
	go func() {
		wg.Wait()
		close(imageURLChan)
		close(backgroundURLChan)
	}()

	// Получаем новый URL аватара
	var newAvatarURL string
	for newURL := range imageURLChan {
		newAvatarURL = newURL
	}

	// Получаем новый URL фона
	var newBackgroundURL string
	for newURL := range backgroundURLChan {
		newBackgroundURL = newURL
	}

	// Заменяем старые URL на новые
	user.Avatar = newAvatarURL
	user.Background = newBackgroundURL

	// Возвращаем пользователя
	// Возвращаем пользователя
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ChangeUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// Получаем файлы из formdata
	var avatarURL, backgroundURL string
	var username, bio string

	// Получаем данные из формы
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		log.Print(err.Error())
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем имя пользователя и биографию
	username = r.FormValue("username")
	bio = r.FormValue("bio")

	// Проверяем наличие файла аватара
	if len(r.MultipartForm.File["avatar"]) > 0 {
		avatarFile, err := r.MultipartForm.File["avatar"][0].Open()
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Unable to open avatar file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer avatarFile.Close()

		// Загружаем аватар
		avatarYid := r.MultipartForm.File["avatar"][0].Filename
		avatarURL, err = UploadImage(avatarFile, avatarYid)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to upload avatar: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Проверяем наличие файла обложки
	if len(r.MultipartForm.File["background"]) > 0 {
		backgroundFile, err := r.MultipartForm.File["background"][0].Open()
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Unable to open background file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer backgroundFile.Close()

		// Загружаем обложку
		backgroundYid := r.MultipartForm.File["background"][0].Filename
		backgroundURL, err = UploadImage(backgroundFile, backgroundYid)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to upload background: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Создаем объект пользователя с обновленными данными
	user := User{
		ID:         userID,
		Username:   username,
		Avatar:     avatarURL,
		Bio:        bio,
		Background: backgroundURL,
	}

	// Отправляем данные пользователя на другой сервер
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/users/changeuser/%d", userID)
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to marshal user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to send user data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status: %s", resp.Status)
		http.Error(w, "Failed to change user data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User data updated successfully"))
}

func (h *Handler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Получаем ID подписчика из URL
	followeeID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	// Формируем JSON для запроса
	follower := UserFollower{
		FollowerID: userID,
		FolloweeID: followeeID,
	}

	jsonData, err := json.Marshal(follower)
	if err != nil {
		http.Error(w, "Failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем POST запрос
	url := "http://45.156.21.196:8003/database_zov_russ_cbo/user_followers"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to send add follower request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to add follower: %s", body), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "successfully added follower"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UnSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Получаем ID подписчика из URL
	followeeID := chi.URLParam(r, "id")

	// Формируем URL для запроса
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/user_followers/%d/%s", userID, followeeID)

	// Создаем DELETE запрос
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, "Failed to create remove follower request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send remove follower request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to remove follower: %s", body), resp.StatusCode)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "successfully removed follower"}
	json.NewEncoder(w).Encode(response)
}
func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Print("Fetching subscriptions")
	userID := r.Context().Value("userID").(int) // Получаем ID пользователя из контекста

	// Создаем запрос к сервису для получения подписок
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/user_followers/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to get subscriptions", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get subscriptions", http.StatusNotFound)
		return
	}

	// Декодируем ответ
	var followers []UserFollower
	if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
		http.Error(w, "Failed to decode subscriptions", http.StatusInternalServerError)
		return
	}

	// Возвращаем подписки
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func (h *Handler) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	log.Print("Fetching subscriptions")
	userID := r.Context().Value("userID").(int) // Получаем ID пользователя из контекста

	// Создаем запрос к сервису для получения подписок
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/user_followers/follower/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to get subscriptions", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get subscriptions", http.StatusNotFound)
		return
	}

	// Декодируем ответ
	var followers []UserFollower
	if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
		http.Error(w, "Failed to decode subscriptions", http.StatusInternalServerError)
		return
	}

	// Возвращаем подписки
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
