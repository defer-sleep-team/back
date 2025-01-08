package post

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	IsNSFW      bool      `json:"is_nsfw"`
	RegDate     time.Time `json:"reg_date"`
}
type PostDetailsSecond struct {
	ID               int       `json:"id"`
	Description      string    `json:"description"`
	IsPrivate        bool      `json:"is_private"`
	IsNSFW           bool      `json:"is_nsfw"`
	RegDate          time.Time `json:"reg_date"`
	Tags             []string  `json:"tags"`
	ImageURLs        []string  `json:"image_urls"`
	UserID           int       `json:"user_id"`
	Username         string    `json:"username"`
	Avatar           string    `json:"avatar"`
	Likes            int       `json:"likes"`
	Views            int       `json:"views"`
	IsLiked          bool      `json:"is_liked"`
	SubscribersCount int       `json:"subscribers_count"`
}

// IncomingPostRequest - структура для входящего запроса на создание поста
type IncomingPostRequest struct {
	Description string   `json:"description"`
	IsPrivate   bool     `json:"is_private"`
	IsNSFW      bool     `json:"is_nsfw"`
	UserID      int      `json:"user_id"`    // Имя пользователя
	Tags        []string `json:"tags"`       // Список тегов
	ImageURLs   []string `json:"image_urls"` // Список URL изображений
}

// OutgoingPostResponse - структура для исходящего ответа после создания поста
type OutgoingPostResponse struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	IsNSFW      bool      `json:"is_nsfw"`
	RegDate     time.Time `json:"reg_date"`
}

// PostDetails - структура для полного представления поста
type PostDetails struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	IsNSFW      bool      `json:"is_nsfw"`
	RegDate     time.Time `json:"reg_date"`
	Tags        []string  `json:"tags"`       // Список тегов
	ImageURLs   []string  `json:"image_urls"` // Список URL изображений
}

// CreateFullPost функция для создания поста с учетом всех связей
func CreateFullPost(db *sql.DB, request IncomingPostRequest) (OutgoingPostResponse, error) {
	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return OutgoingPostResponse{}, err
	}

	// Создание поста
	post := Post{
		Description: request.Description,
		IsPrivate:   request.IsPrivate,
		IsNSFW:      request.IsNSFW,
		RegDate:     time.Now(),
	}

	postID, err := CreatePost(db, post)
	if err != nil {
		tx.Rollback()
		return OutgoingPostResponse{}, err
	}

	// Получение user_id по имени пользователя
	var userID int
	err = tx.QueryRow(`SELECT id FROM users WHERE id=$1`, request.UserID).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return OutgoingPostResponse{}, err
	}

	// Создание связи с пользователем
	_, err = tx.Exec(`INSERT INTO user_posts (user_id, post_id) VALUES ($1, $2)`, userID, postID)
	if err != nil {
		tx.Rollback()
		return OutgoingPostResponse{}, err
	}

	// Создание связей с тегами
	for _, tag := range request.Tags {
		var tagID int
		err = tx.QueryRow(`SELECT id FROM tags WHERE name=$1`, tag).Scan(&tagID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Если тег не найден, вставляем новый тег
				insertQuery := `INSERT INTO tags (name) VALUES ($1) RETURNING id`
				err = tx.QueryRow(insertQuery, tag).Scan(&tagID)
				if err != nil {
					tx.Rollback()
					return OutgoingPostResponse{}, err
				}
			} else {
				tx.Rollback()
				return OutgoingPostResponse{}, err
			}
		}
		_, err = tx.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`, postID, tagID)
		if err != nil {
			tx.Rollback()
			return OutgoingPostResponse{}, err
		}
	}

	// Создание связей с изображениями
	for _, imageURL := range request.ImageURLs {
		_, err = tx.Exec(`INSERT INTO post_images (post_id, image_url) VALUES ($1, $2)`, postID, imageURL)
		if err != nil {
			tx.Rollback()
			return OutgoingPostResponse{}, err
		}
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return OutgoingPostResponse{}, err
	}

	// Формирование ответа
	response := OutgoingPostResponse{
		ID:          postID,
		Description: post.Description,
		IsPrivate:   post.IsPrivate,
		IsNSFW:      post.IsNSFW,
		RegDate:     post.RegDate,
	}

	return response, nil

}

// UpdateRatio - функция для изменения ratio на заданное значение
func UpdateRatio(db *sql.DB, postID int, delta int) error {
	query := `UPDATE ratios SET ratio = ratio + $1 WHERE post_id = $2`
	_, err := db.Exec(query, delta, postID)
	return err
}

// всё добавил и статитсику и всё заеб
func GetNPostsByRatio(db *sql.DB, offset int, n int, userID int) ([]PostDetailsSecond, error) {
	query := `
		SELECT p.id, p.description, p.is_private, p.is_nsfw, p.reg_date, u.id AS user_id, u.username, u.avatar, ps.likes, ps.views,
		       EXISTS (SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $2) AS is_liked,
		       (SELECT COUNT(*) FROM user_followers uf WHERE uf.followee_id = u.id) AS subscribers_count
		FROM posts p
		JOIN ratios r ON p.id = r.post_id
		JOIN user_posts up ON p.id = up.post_id
		JOIN users u ON up.user_id = u.id
		JOIN post_stats ps ON p.id = ps.post_id
		ORDER BY r.ratio DESC
		LIMIT $1 OFFSET $3`

	rows, err := db.Query(query, n, userID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostDetailsSecond
	for rows.Next() {
		var post PostDetailsSecond
		if err := rows.Scan(&post.ID, &post.Description, &post.IsPrivate, &post.IsNSFW, &post.RegDate,
			&post.UserID, &post.Username, &post.Avatar, &post.Likes, &post.Views, &post.IsLiked, &post.SubscribersCount); err != nil {
			log.Print(0)
			return nil, err
		}

		// Add view to post
		AddViewToPost(db, post.ID)

		// Получение тегов для поста
		tagQuery := `SELECT t.name
                     FROM tags t
                     JOIN post_tags pt ON t.id = pt.tag_id
                     WHERE pt.post_id = $1`
		tagRows, err := db.Query(tagQuery, post.ID)
		if err != nil {
			log.Print(1)
			return nil, err
		}
		defer tagRows.Close()

		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				log.Print(2)
				return nil, err
			}
			post.Tags = append(post.Tags, tag)
		}

		// Получение изображений для поста
		imageQuery := `SELECT pi.image_url
                         FROM post_images pi
                         WHERE pi.post_id = $1`
		imageRows, err := db.Query(imageQuery, post.ID)
		if err != nil {
			log.Print(3)
			return nil, err
		}
		defer imageRows.Close()

		for imageRows.Next() {
			var imageURL string
			if err := imageRows.Scan(&imageURL); err != nil {
				log.Print(4)
				return nil, err
			}
			post.ImageURLs = append(post.ImageURLs, imageURL)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// Функция для получения постов пользователя
func GetNPostsOfUser(db *sql.DB, userID int, n int) ([]PostDetailsSecond, error) {
	query := `
		SELECT p.id, p.description, p.is_private, p.is_nsfw, p.reg_date,
		       u.id AS user_id, u.username, u.avatar,
		       COALESCE(ps.likes, 0) AS likes, COALESCE(ps.views, 0) AS views,
		       (SELECT COUNT(*) FROM user_followers WHERE followee_id = u.id) AS followers
		FROM posts p
		JOIN user_posts up ON p.id = up.post_id
		JOIN users u ON up.user_id = u.id
		LEFT JOIN post_stats ps ON p.id = ps.post_id
		WHERE up.user_id = $1
		ORDER BY p.reg_date DESC
		LIMIT $2`

	rows, err := db.Query(query, userID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostDetailsSecond
	for rows.Next() {
		var post PostDetailsSecond
		if err := rows.Scan(&post.ID, &post.Description, &post.IsPrivate, &post.IsNSFW, &post.RegDate,
			&post.UserID, &post.Username, &post.Avatar, &post.Likes, &post.Views, &post.SubscribersCount); err != nil {
			return nil, err
		}

		// Получение тегов для поста
		tagQuery := `SELECT t.name
                     FROM tags t
                     JOIN post_tags pt ON t.id = pt.tag_id
                     WHERE pt.post_id = $1`
		tagRows, err := db.Query(tagQuery, post.ID)
		if err != nil {
			return nil, err
		}
		defer tagRows.Close()

		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				return nil, err
			}
			post.Tags = append(post.Tags, tag)
		}

		// Получение изображений для поста
		imageQuery := `SELECT pi.image_url
                         FROM post_images pi
                         WHERE pi.post_id = $1`
		imageRows, err := db.Query(imageQuery, post.ID)
		if err != nil {
			return nil, err
		}
		defer imageRows.Close()

		for imageRows.Next() {
			var imageURL string
			if err := imageRows.Scan(&imageURL); err != nil {
				return nil, err
			}
			post.ImageURLs = append(post.ImageURLs, imageURL)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetNPostsBySubscription(db *sql.DB, userID int, n int) ([]PostDetails, error) {
	query := `
		SELECT p.id, p.description, p.is_private, p.is_nsfw, p.reg_date
		FROM posts p
		JOIN user_posts up ON p.id = up.post_id
		JOIN user_followers uf ON up.user_id = uf.followee_id
		JOIN ratios r ON p.id = r.post_id
		WHERE uf.follower_id = $1
		ORDER BY r.ratio DESC, p.reg_date DESC
		LIMIT $2`

	rows, err := db.Query(query, userID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostDetails
	for rows.Next() {
		var post PostDetails
		if err := rows.Scan(&post.ID, &post.Description, &post.IsPrivate, &post.IsNSFW, &post.RegDate); err != nil {
			return nil, err
		}

		// Получение тегов для поста
		tagQuery := `SELECT t.name
                     FROM tags t
                     JOIN post_tags pt ON t.id = pt.tag_id
                     WHERE pt.post_id = $1`
		tagRows, err := db.Query(tagQuery, post.ID)
		if err != nil {
			return nil, err
		}
		defer tagRows.Close()

		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				return nil, err
			}
			post.Tags = append(post.Tags, tag)
		}

		// Получение изображений для поста
		imageQuery := `SELECT pi.image_url
                         FROM post_images pi
                         WHERE pi.post_id = $1`
		imageRows, err := db.Query(imageQuery, post.ID)
		if err != nil {
			return nil, err
		}
		defer imageRows.Close()

		for imageRows.Next() {
			var imageURL string
			if err := imageRows.Scan(&imageURL); err != nil {
				return nil, err
			}
			post.ImageURLs = append(post.ImageURLs, imageURL)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(db *sql.DB, post Post) (int, error) {
	query := `INSERT INTO posts (description, is_private, is_nsfw, reg_date)
              VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	err := db.QueryRow(query, post.Description, post.IsPrivate, post.IsNSFW, post.RegDate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// описания поста
func UpdatePostDescription(db *sql.DB, postID int, description string) error {
	query := `UPDATE posts SET description = $1 WHERE id = $2`
	_, err := db.Exec(query, description, postID)
	return err
}

// UpdatePostTags обновляет теги поста
func UpdatePostTags(db *sql.DB, postID int, tagNames []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Удаление всех текущих тегов поста
	_, err = tx.Exec(`DELETE FROM post_tags WHERE post_id = $1`, postID)
	if err != nil {
		return err
	}

	// Добавление новых тегов
	for _, tagName := range tagNames {
		var tagID int
		query := `SELECT id FROM tags WHERE name = $1`
		err := tx.QueryRow(query, tagName).Scan(&tagID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Если тег не найден, вставляем новый тег
				insertQuery := `INSERT INTO tags (name) VALUES ($1) RETURNING id`
				err = tx.QueryRow(insertQuery, tagName).Scan(&tagID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Добавление связи между постом и тегом
		_, err = tx.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`, postID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPost - функция для получения одного поста с изменением ratio и всеми его данными
func GetPost(db *sql.DB, id int) (PostDetails, error) {

	// Получение поста
	postQuery := `SELECT p.id, p.description, p.is_private, p.is_nsfw, p.reg_date 
                  FROM posts p 
                  WHERE p.id = $1`
	var post PostDetails
	err := db.QueryRow(postQuery, id).Scan(&post.ID, &post.Description, &post.IsPrivate, &post.IsNSFW, &post.RegDate)
	if err != nil {
		return PostDetails{}, err
	}

	// Получение тегов
	tagQuery := `SELECT t.name 
                 FROM tags t 
                 JOIN post_tags pt ON t.id = pt.tag_id 
                 WHERE pt.post_id = $1`
	rows, err := db.Query(tagQuery, id)
	if err != nil {
		return PostDetails{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return PostDetails{}, err
		}
		post.Tags = append(post.Tags, tag)
	}

	// Получение изображений
	imageQuery := `SELECT pi.image_url 
                     FROM post_images pi 
                     WHERE pi.post_id = $1`
	imageRows, err := db.Query(imageQuery, id)
	if err != nil {
		return PostDetails{}, err
	}
	defer imageRows.Close()

	for imageRows.Next() {
		var imageURL string
		if err := imageRows.Scan(&imageURL); err != nil {
			return PostDetails{}, err
		}
		post.ImageURLs = append(post.ImageURLs, imageURL)
	}

	return post, nil
}

func UpdatePost(db *sql.DB, post Post) error {
	query := `UPDATE posts SET description=$1, is_private=$2, is_nsfw=$3, reg_date=$4 WHERE id=$5`

	_, err := db.Exec(query, post.Description, post.IsPrivate, post.IsNSFW, post.RegDate, post.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeletePostAdmin(db *sql.DB, postID int) error {
	// Проверка принадлежности поста пользователю
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_posts WHERE post_id = $1)`, postID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user does not own this post")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Удаление связанных данных
	queries := []string{
		`DELETE FROM post_likes WHERE post_id = $1`,
		`DELETE FROM post_comments WHERE post_id = $1`,
		`DELETE FROM post_images WHERE post_id = $1`,
		`DELETE FROM post_tags WHERE post_id = $1`,
		`DELETE FROM post_stats WHERE post_id = $1`,
		`DELETE FROM user_posts WHERE post_id = $1`,
		`DELETE FROM ratios WHERE post_id = $1`,
		`DELETE FROM plan_posts WHERE post_id = $1`,
	}

	for _, query := range queries {
		_, err = tx.Exec(query, postID)
		if err != nil {
			return err
		}
	}

	// Удаление самого поста
	_, err = tx.Exec(`DELETE FROM posts WHERE id = $1`, postID)
	if err != nil {
		return err
	}

	return nil
}

func DeletePost(db *sql.DB, postID, uid int) error {
	// Проверка принадлежности поста пользователю
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_posts WHERE post_id = $1 AND user_id = $2)`, postID, uid).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user does not own this post")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Удаление связанных данных
	queries := []string{
		`DELETE FROM post_likes WHERE post_id = $1`,
		`DELETE FROM post_comments WHERE post_id = $1`,
		`DELETE FROM post_images WHERE post_id = $1`,
		`DELETE FROM post_tags WHERE post_id = $1`,
		`DELETE FROM post_stats WHERE post_id = $1`,
		`DELETE FROM user_posts WHERE post_id = $1`,
		`DELETE FROM ratios WHERE post_id = $1`,
		`DELETE FROM plan_posts WHERE post_id = $1`,
	}

	for _, query := range queries {
		_, err = tx.Exec(query, postID)
		if err != nil {
			return err
		}
	}

	// Удаление самого поста
	_, err = tx.Exec(`DELETE FROM posts WHERE id = $1`, postID)
	if err != nil {
		return err
	}

	return nil
}

// Теги как и просил с инсёртом
func TagIDByName(db *sql.DB, tagNames []string) (map[string]int, error) {
	tagIDs := make(map[string]int)

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, tagName := range tagNames {
		var tagID int
		query := `SELECT id FROM tags WHERE name = $1`
		err := tx.QueryRow(query, tagName).Scan(&tagID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Если тег не найден, вставляем новый тег
				insertQuery := `INSERT INTO tags (name) VALUES ($1) RETURNING id`
				err = tx.QueryRow(insertQuery, tagName).Scan(&tagID)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		tagIDs[tagName] = tagID
	}

	return tagIDs, nil
}

// часть с лайками
func IsLiked(db *sql.DB, postID, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`
	var exists bool
	err := db.QueryRow(query, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func LikePost(db *sql.DB, postID, userID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, postID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	updateQuery := `UPDATE post_stats SET likes = likes + 1 WHERE post_id = $1`
	_, err = tx.Exec(updateQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	AddUpdateRatio(db, postID)

	return nil
}

func UnlikePost(db *sql.DB, postID, userID int) error {
	// Start a transaction
	tx, err := db.Begin()
	isLiked, err := IsLiked(db, postID, userID)
	if err != nil {
		return err
	}

	if isLiked {
		// Delete the like from the post_likes table
		query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`
		_, err = tx.Exec(query, postID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Update the post_stats table to decrement the likes count
		updateQuery := `UPDATE post_stats SET likes = likes - 1 WHERE post_id = $1`
		_, err = tx.Exec(updateQuery, postID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			return err
		}

		AddUpdateRatioDown(db, postID)
	}

	return nil
}

// AddComment добавляет новый комментарий к посту
func AddComment(db *sql.DB, postID, userID int, content string) error {
	query := `INSERT INTO comments (content, reg_date) VALUES ($1, $2) RETURNING id`
	var commentID int
	err := db.QueryRow(query, content, time.Now()).Scan(&commentID)
	if err != nil {
		return err
	}

	query = `INSERT INTO post_comments (post_id, comment_id) VALUES ($1, $2)`
	_, err = db.Exec(query, postID, commentID)
	if err != nil {
		return err
	}

	query = `INSERT INTO user_comments (user_id, comment_id) VALUES ($1, $2)`
	_, err = db.Exec(query, userID, commentID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComment(db *sql.DB, commentID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Удаление связанных данных
	queries := []string{
		`DELETE FROM post_comments WHERE comment_id = $1`,
		`DELETE FROM user_comments WHERE comment_id = $1`,
	}

	for _, query := range queries {
		_, err = tx.Exec(query, commentID)
		if err != nil {
			return err
		}
	}

	// Удаление самого комментария
	_, err = tx.Exec(`DELETE FROM comments WHERE id = $1`, commentID)
	if err != nil {
		return err
	}

	return nil
}

func AddViewToPost(db *sql.DB, postID int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Update the post_stats table to increment the views count
	updateQuery := `UPDATE post_stats SET views = views + 1 WHERE post_id = $1`
	_, err = tx.Exec(updateQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	AddViewToPostAndUpdateRatio(db, postID)
	return nil

}

func AddViewToPostAndUpdateRatio(db *sql.DB, postID int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Get the current ratio from the ratios table
	var currentRatio int
	getRatioQuery := `SELECT ratio FROM ratios WHERE post_id = $1`
	err = tx.QueryRow(getRatioQuery, postID).Scan(&currentRatio)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Calculate the new ratio (current ratio - 7)
	newRatio := currentRatio - 7

	// Update the ratios table with the new ratio
	updateRatioQuery := `UPDATE ratios SET ratio = $1 WHERE post_id = $2`
	_, err = tx.Exec(updateRatioQuery, newRatio, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func AddUpdateRatio(db *sql.DB, postID int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Get the current ratio from the ratios table
	var currentRatio int
	getRatioQuery := `SELECT ratio FROM ratios WHERE post_id = $1`
	err = tx.QueryRow(getRatioQuery, postID).Scan(&currentRatio)
	if err != nil {
		tx.Rollback()
		return err
	}

	newRatio := currentRatio + 400

	// Update the ratios table with the new ratio
	updateRatioQuery := `UPDATE ratios SET ratio = $1 WHERE post_id = $2`
	_, err = tx.Exec(updateRatioQuery, newRatio, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func AddUpdateRatioDown(db *sql.DB, postID int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Get the current ratio from the ratios table
	var currentRatio int
	getRatioQuery := `SELECT ratio FROM ratios WHERE post_id = $1`
	err = tx.QueryRow(getRatioQuery, postID).Scan(&currentRatio)
	if err != nil {
		tx.Rollback()
		return err
	}

	newRatio := currentRatio - 400

	// Update the ratios table with the new ratio
	updateRatioQuery := `UPDATE ratios SET ratio = $1 WHERE post_id = $2`
	_, err = tx.Exec(updateRatioQuery, newRatio, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
