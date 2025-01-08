package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

const DatabaseIP = "http://45.156.21.196:8003/database_zov_russ_cbo/"

// IncomingPostRequest - структура для входящего запроса на создание поста
type IncomingPostRequest struct {
	Description string   `json:"description"`
	IsPrivate   bool     `json:"is_private"`
	IsNSFW      bool     `json:"is_nsfw"`
	UserID      int      `json:"user_id"`    // Имя пользователя
	Tags        []string `json:"tags"`       // Список идентификаторов тегов
	ImageURLs   []string `json:"image_urls"` // Список URL изображений
}

// PostDetails - структура для детальной информации о посте
type PostDetails struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	IsNSFW      bool      `json:"is_nsfw"`
	RegDate     time.Time `json:"reg_date"`
	Tags        []string  `json:"tags"`
	ImageURLs   []string  `json:"image_urls"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Avatar      string    `json:"avatar"`
	Likes       int       `json:"likes"`
	Views       int       `json:"views"`
	IsLiked     bool      `json:"is_liked"`
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Unable to parse form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	postData := r.FormValue("postData")
	var request IncomingPostRequest
	if err := json.Unmarshal([]byte(postData), &request); err != nil {
		http.Error(w, "Invalid post data: "+err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
		return
	}
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Print(err.Error())
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	request.UserID = userID

	var imageNames []string
	for _, fileHeader := range r.MultipartForm.File["images"] {
		file, err := fileHeader.Open()
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Unable to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		originalFilename := fileHeader.Filename

		imageName, err := UploadImage(file, originalFilename)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		imageNames = append(imageNames, imageName)
	}

	request.ImageURLs = imageNames

	if err := savePostToDatabase(request); err != nil {
		log.Print(err)
		http.Error(w, "Failed to save post to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "successfully created a post"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	postIDStr := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Формируем URL для DELETE-запроса
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/posts/%s/%d", postIDStr, userID)

	// Выполняем DELETE-запрос
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to delete post: "+resp.Status, resp.StatusCode)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Post deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeletePostHandlerSudo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	postIDStr := chi.URLParam(r, "id")
	role, ok := r.Context().Value("privileges").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}
	if role != 2 {
		http.Error(w, "You don't have permission to delete posts", http.StatusForbidden)
		return
	}
	// Формируем URL для DELETE-запроса
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/posts/sudo/delete/post/%s", postIDStr)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to delete post: "+resp.Status, resp.StatusCode)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Post deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	postID := chi.URLParam(r, "postID")

	// Создаем запрос к Fiber-сервису
	resp, err := http.Get(DatabaseIP + "posts/" + postID)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Декодируем ответ
	var postDetails PostDetails
	if err := json.NewDecoder(resp.Body).Decode(&postDetails); err != nil {
		http.Error(w, "Failed to decode post", http.StatusInternalServerError)
		return
	}

	// Создаем WaitGroup для ожидания завершения всех запросов
	var wg sync.WaitGroup
	imageURLChan := make(chan string, len(postDetails.ImageURLs))

	// Отправляем запросы для каждого URL в image_urls
	for _, imageURL := range postDetails.ImageURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			// Формируем полный адрес для запроса
			fullURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + url
			resp, err := http.Get(fullURL)
			if err != nil {
				log.Print(err)
				imageURLChan <- url // Если ошибка, отправляем оригинальный URL
				return
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if resp.StatusCode != http.StatusOK {
				log.Print(resp.StatusCode)
				imageURLChan <- url // Если не найден, отправляем оригинальный URL
				return
			}

			// Декодируем ответ
			var result struct {
				URL string `json:"url"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				log.Print(err)
				imageURLChan <- url // Если ошибка, отправляем оригинальный URL
				return
			}

			// Отправляем полученный URL в канал
			imageURLChan <- result.URL
		}(imageURL)
	}

	// Закрываем канал после завершения всех запросов
	go func() {
		wg.Wait()
		close(imageURLChan)
	}()

	// Собираем новые URL
	var newImageURLs []string
	for newURL := range imageURLChan {
		newImageURLs = append(newImageURLs, newURL)
	}

	// Заменяем старые URL на новые
	postDetails.ImageURLs = newImageURLs

	// Возвращаем пост
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postDetails)
}

// savePostToDatabase - функция для сохранения поста в базу данных
func savePostToDatabase(request IncomingPostRequest) error {
	url := "http://45.156.21.196:8003/database_zov_russ_cbo/posts/full"

	// Преобразуем структуру в JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Print(err.Error())
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Создаем POST-запрос
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(err)
		return fmt.Errorf("error sending request to database: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status from database: %s", resp.Status)
	}

	return nil
}

// Mass Post-Getter 1 godoc
// @Summary Получить трендовые посты (с самым высоким коэффициентом рекомендаций)
// @Description Function available for all users, returns post by its' id. Для корректной работы требуется куки "watched", которая содержит количество уже прогруженных постов.
// @Tags posts
// @Produce json
// @Success 200 {object} PostDetails "Даёт посты, прими посты. В массиве имен картинок каждая картинка должна быть отдельно подгружена, даются готовые ссылки"
//
// @Failure 404 {object} nil "Нет такого поста"
// @Failure 500 {object} nil "Отвалилась жопа. Или картинка не грузится или прислали говна на лопате"
// @Router /api/get/trends/{n} [get]
func (h *Handler) GetTrends(w http.ResponseWriter, r *http.Request) {
	log.Print("Trends gotten")
	n := chi.URLParam(r, "n")
	userID, ok := r.Context().Value("userID").(int) // Получаем ID пользователя из контекста

	watched, ok := r.Context().Value("watched").(int)

	if !ok {
		watched = 0
	}
	numPosts, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, "Invalid number of posts", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(DatabaseIP + "posts/ratio/" + strconv.Itoa(watched) + "/" + strconv.Itoa(numPosts) + "/" + fmt.Sprintf("%d", userID))
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get recommendations", http.StatusNotFound)
		return
	}

	var posts []PostDetails
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		http.Error(w, "Failed to decode recommendations", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	imageURLChan := make(chan struct {
		PostID   int
		ImageURL string
	}, len(posts))

	for _, post := range posts {
		wg.Add(1)
		go func(post *PostDetails) {
			defer wg.Done()
			var successfulURLs []string
			for _, imageURL := range post.ImageURLs {
				fullURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + imageURL
				resp, err := http.Get(fullURL)
				if err != nil {
					log.Print(err)
					successfulURLs = append(successfulURLs, imageURL)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					log.Print(resp.StatusCode)
					successfulURLs = append(successfulURLs, imageURL)
					continue
				}

				var result struct {
					URL string `json:"url"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					log.Print(err)
					successfulURLs = append(successfulURLs, imageURL)
					continue
				}

				successfulURLs = append(successfulURLs, result.URL)
			}
			imageURLChan <- struct {
				PostID   int
				ImageURL string
			}{PostID: post.ID, ImageURL: strings.Join(successfulURLs, ",")}
		}(&post)
	}

	go func() {
		wg.Wait()
		close(imageURLChan)
	}()

	for imageInfo := range imageURLChan {
		for i := range posts {
			if posts[i].ID == imageInfo.PostID {
				posts[i].ImageURLs = strings.Split(imageInfo.ImageURL, ",")
			}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range posts {
			avatarURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + posts[i].Avatar
			resp, err := http.Get(avatarURL)
			if err != nil {
				log.Print(err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Print(resp.StatusCode)
				continue
			}

			var result struct {
				URL string `json:"url"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				log.Print("Error decoding trends json: ", err)
				continue
			}

			posts[i].Avatar = result.URL

		}
	}()

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetRecs - обработчик для получения рекомендованных постов
func (h *Handler) GetRecs(w http.ResponseWriter, r *http.Request) {
	log.Print("Z1")
	n := chi.URLParam(r, "n")
	userID := r.Context().Value("userID").(int) // Получаем ID пользователя из контекста

	numPosts, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, "Invalid number of posts", http.StatusBadRequest)
		return
	}

	// Создаем запрос к Fiber-сервису
	resp, err := http.Get(DatabaseIP + "posts/" + strconv.Itoa(userID) + "/recommendations/0/" + strconv.Itoa(numPosts))
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get recommendations", http.StatusNotFound)
		return
	}

	// Декодируем ответ
	var posts []PostDetails
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		http.Error(w, "Failed to decode recommendations", http.StatusInternalServerError)
		return
	}

	// Возвращаем рекомендованные посты
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Mass Post-Getter 2 godoc
// Posts godoc
// @Summary Получить посты одного пользователя
// @Description Здесь пока без клёвой прогрузки, но скоро намечается через watched)
// @Tags posts
// @Produce json
// @Success 200 {object} PostDetails. "Даёт посты, прими посты. В массиве имен картинок каждая картинка должна быть отдельно подгружена, даются готовые ссылки"
// @Failure 404 {object} nil "накосячил тот, кто делал запрос"
// @Failure 500 {object} nil "Отвалилась жопа, накосячил Игорь (или я)"
// @Router api/get/posts_of/{userID}/{n}" [get]
func (h *Handler) GetPostsOf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	log.Print("Z1")
	n := chi.URLParam(r, "n")
	userID := chi.URLParam(r, "userID")

	numPosts, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, "Invalid number of posts", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(DatabaseIP + "posts/posts_of/" + userID + "/" + strconv.Itoa(numPosts))
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get recommendations", http.StatusNotFound)
		return
	}

	var posts []PostDetails
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		http.Error(w, "Failed to decode recommendations", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	imageURLChan := make(chan struct {
		PostID    int
		ImageURLs []string
	}, len(posts))

	// Обработка всех ImageURLs для каждого поста
	for _, post := range posts {
		wg.Add(1)
		go func(post *PostDetails) {
			defer wg.Done()
			var successfulImageURLs []string

			// Обработка всех ImageURLs
			for _, imageURL := range post.ImageURLs {
				fullURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + imageURL
				resp, err := http.Get(fullURL)
				if err != nil {
					log.Print(err)
					successfulImageURLs = append(successfulImageURLs, imageURL) // Возвращаем оригинальный URL в случае ошибки
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					log.Print(resp.StatusCode)
					successfulImageURLs = append(successfulImageURLs, imageURL) // Возвращаем оригинальный URL в случае ошибки
					continue
				}

				var result struct {
					URL string `json:"url"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					log.Print(err)
					successfulImageURLs = append(successfulImageURLs, imageURL) // Возвращаем оригинальный URL в случае ошибки
					continue
				}

				successfulImageURLs = append(successfulImageURLs, result.URL)
			}

			imageURLChan <- struct {
				PostID    int
				ImageURLs []string
			}{PostID: post.ID, ImageURLs: successfulImageURLs}
		}(&post)
	}

	go func() {
		wg.Wait()
		close(imageURLChan)
	}()

	// Обновление ImageURLs в постах
	for imageInfo := range imageURLChan {
		for i := range posts {
			if posts[i].ID == imageInfo.PostID {
				posts[i].ImageURLs = imageInfo.ImageURLs
			}
		}
	}

	// Теперь обрабатываем аватары после обработки всех ImageURLs
	var avatarWg sync.WaitGroup
	avatarChan := make(chan struct {
		PostID    int
		AvatarURL string
	}, len(posts))

	for _, post := range posts {
		avatarWg.Add(1)
		go func(post *PostDetails) {
			defer avatarWg.Done()
			avatarURL := "http://45.156.21.196:8765/api_s3_aether_amazon/get_url/" + post.Avatar
			resp, err := http.Get(avatarURL)
			if err != nil {
				log.Print(err)
				avatarChan <- struct {
					PostID    int
					AvatarURL string
				}{PostID: post.ID, AvatarURL: post.Avatar} // Возвращаем оригинальный URL в случае ошибки
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Print(resp.StatusCode)
				avatarChan <- struct {
					PostID    int
					AvatarURL string
				}{PostID: post.ID, AvatarURL: post.Avatar} // Возвращаем оригинальный URL в случае ошибки
				return
			}
			var avatarResult struct {
				URL string `json:"url"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&avatarResult); err != nil {
				log.Print(err)
				avatarChan <- struct {
					PostID    int
					AvatarURL string
				}{PostID: post.ID, AvatarURL: post.Avatar} // Возвращаем оригинальный URL в случае ошибки
				return
			}

			// Отправляем обновленный URL аватара в канал
			avatarChan <- struct {
				PostID    int
				AvatarURL string
			}{PostID: post.ID, AvatarURL: avatarResult.URL}
		}(&post)
	}

	go func() {
		avatarWg.Wait()
		close(avatarChan)
	}()

	// Обновление аватаров в постах
	for avatarInfo := range avatarChan {
		for i := range posts {
			if posts[i].ID == avatarInfo.PostID {
				posts[i].Avatar = avatarInfo.AvatarURL
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *Handler) LikeHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Z123123")
	postID := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/posts/like/%s/%d", postID, userID)

	// Отправляем POST запрос
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		http.Error(w, "Failed to send like request: "+err.Error(), http.StatusAlreadyReported)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusAlreadyReported)
		return
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to like post: %s", body), http.StatusAlreadyReported)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "successfully liked the post"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UnlikeHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID поста из URL
	postID := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/posts/unlike/%s/%d", postID, userID)

	// Создаем DELETE запрос
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, "Failed to create unlike request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send unlike request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to unlike post: %s", body), resp.StatusCode)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "successfully unliked the post"}
	json.NewEncoder(w).Encode(response)
}
