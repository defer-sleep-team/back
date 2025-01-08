package post

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Comments struct {
	ID       int       `json:"id"`
	Content  string    `json:"content"`
	RegDate  time.Time `json:"reg_date"`
	UserID   int       `json:"user_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	PostID   int       `json:"post_id"` // Добавлено поле PostID
}

type CommentRequest struct {
	Content string `json:"content"`
}

var ErrUser = errors.New("user have wrong id")

// CreateComment создает новый комментарий
func CreateComment(db *sql.DB, comment *Comments) error {
	query := `INSERT INTO comments (content, reg_date)
              VALUES ($1, $2)
              RETURNING id;`
	err := db.QueryRow(query, comment.Content, comment.RegDate).Scan(&comment.ID)
	if err != nil {
		return err
	}

	// Создание связи с пользователем
	_, err = db.Exec(`INSERT INTO user_comments (user_id, comment_id) VALUES ($1, $2)`, comment.UserID, comment.ID)
	if err != nil {
		return err
	}

	// Создание связи с постом
	_, err = db.Exec(`INSERT INTO post_comments (post_id, comment_id) VALUES ($1, $2)`, comment.PostID, comment.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetComment получает комментарий по его идентификатору
func GetComment(db *sql.DB, id int) (*Comments, error) {
	var comment Comments
	query := `SELECT c.id, c.content, c.reg_date, u.id AS user_id, u.username, u.avatar, pc.post_id
              FROM comments c
              JOIN user_comments uc ON c.id = uc.comment_id
              JOIN users u ON uc.user_id = u.id
              JOIN post_comments pc ON c.id = pc.comment_id
              WHERE c.id = $1`
	err := db.QueryRow(query, id).Scan(&comment.ID, &comment.Content, &comment.RegDate, &comment.UserID, &comment.Username, &comment.Avatar, &comment.PostID)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetAllComments получает все комментарии
func GetAllComments(db *sql.DB) ([]Comments, error) {
	var comments []Comments
	query := `SELECT c.id, c.content, c.reg_date, u.id AS user_id, u.username, u.avatar, pc.post_id
              FROM comments c
              JOIN user_comments uc ON c.id = uc.comment_id
              JOIN users u ON uc.user_id = u.id
              JOIN post_comments pc ON c.id = pc.comment_id`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.RegDate, &comment.UserID, &comment.Username, &comment.Avatar, &comment.PostID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// UpdateComment обновляет комментарий по его идентификатору
func UpdateComment(db *sql.DB, id int, updatedComment *Comments) error {
	query := `UPDATE comments SET content = $1, reg_date = $2 WHERE id = $3`
	_, err := db.Exec(query, updatedComment.Content, updatedComment.RegDate, id)
	return err
}

// DeleteComment удаляет комментарий по его идентификатору
func DeleteComment(db *sql.DB, id int, userID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Проверка владельца комментария в таблице user_comments
	var ownerID int
	err = tx.QueryRow("SELECT user_id FROM user_comments WHERE comment_id = $1", id).Scan(&ownerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Если владелец комментария не совпадает с текущим пользователем, возвращаем ошибку
	if ownerID != userID {
		tx.Rollback()
		return ErrUser
	}

	// Удаление записи из таблицы post_comments
	_, err = tx.Exec("DELETE FROM post_comments WHERE comment_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление записи из таблицы user_comments
	_, err = tx.Exec("DELETE FROM user_comments WHERE comment_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление комментария из таблицы comments
	_, err = tx.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment удаляет комментарий по его идентификатору
func DeleteCommentSudo(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Удаление записи из таблицы post_comments
	_, err = tx.Exec("DELETE FROM post_comments WHERE comment_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление записи из таблицы user_comments
	_, err = tx.Exec("DELETE FROM user_comments WHERE comment_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление комментария из таблицы comments
	_, err = tx.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetCommentsForPost получает n комментариев для поста с заданным идентификатором
func GetCommentsForPost(db *sql.DB, postID, n int) ([]Comments, error) {
	query := `
		SELECT c.id, c.content, c.reg_date, u.id AS user_id, u.username, u.avatar, pc.post_id
		FROM comments c
		JOIN user_comments uc ON c.id = uc.comment_id
		JOIN users u ON uc.user_id = u.id
		JOIN post_comments pc ON c.id = pc.comment_id
		WHERE pc.post_id = $1
		ORDER BY c.reg_date DESC
		LIMIT $2
	`
	rows, err := db.Query(query, postID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.RegDate, &comment.UserID, &comment.Username, &comment.Avatar, &comment.PostID); err != nil {
			return nil, err
		}
		log.Println("comment.ID:", comment.ID)
		log.Println("comment.Content:", comment.Content)
		log.Println("comment.RegDate:", comment.RegDate)
		log.Println("comment.UserID:", comment.UserID)
		log.Println("comment.Username:", comment.Username)

		comments = append(comments, comment)
	}

	return comments, nil
}
