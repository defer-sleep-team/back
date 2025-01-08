package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (h *Handler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	postID := chi.URLParam(r, "postID")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Читаем тело запроса
	var requestBody struct {
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что контент комментария не пустой
	if requestBody.Content == "" {

		log.Print("Error body: ", requestBody)
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Формируем URL для запроса к базе данных
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/posts/addcomment/%s/%d", postID, userID)

	// Создаем JSON для отправки
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Failed to marshal request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем POST запрос
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to send comment request: "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Failed to add comment: %s", body), resp.StatusCode)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Comment added successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	// Получаем postID из параметров URL
	postID := chi.URLParam(r, "postID")

	// Формируем URL для запроса к базе данных
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/comments/comments/%s/100", postID)

	// Отправляем GET запрос
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to send request to get comments: "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Failed to get comments: %s", body), resp.StatusCode)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Возвращаем комментарии клиенту
	w.Write(body)
}

func (h *Handler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	commID := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Формируем URL для DELETE-запроса
	url := fmt.Sprintf("http://45.156.21.196:8003/database_zov_russ_cbo/comments/%s/%d", commID, userID)

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
