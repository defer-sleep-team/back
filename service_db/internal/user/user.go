package user

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"strconv"
	"strings"
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
	SubscribersCount int    `json:"subscribers_count"`
	Background       string `json:"background"`
	IsSubscribed     bool   `json:"is_subscribed"`
}

var ErrUserBlocked = errors.New("user is blocked")

func CreateUser(db *sql.DB, user *User) error {
	var err error
	var query string
	user.Password = HashPassword(user.Password)

	log.Print(user.Email + " " + user.Username)
	if user.Avatar != "" {
		if user.Background != "" {
			query = `INSERT INTO users (yid, username, email, password, avatar, bio, privilege_level, payments, background)
					RETURNING id;`
			err = db.QueryRow(query, user.Yid, user.Username, user.Email, user.Password, user.Avatar, user.Bio, user.PrivilegeLevel, user.Payments, user.Background).Scan(&user.ID)
		} else {
			user.Background = "2.jpg"
			query = `INSERT INTO users (yid, username, email, password, avatar, bio, privilege_level, payments, background)
					RETURNING id;`
			err = db.QueryRow(query, user.Yid, user.Username, user.Email, user.Password, user.Avatar, user.Bio, user.PrivilegeLevel, user.Payments, user.Background).Scan(&user.ID)
		}

	} else {
		if user.Background != "" {
			user.Avatar = "1.jpg"
			query = `INSERT INTO users (yid, username, email, password, bio, privilege_level, payments, background)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id;`
			err = db.QueryRow(query, user.Yid, user.Username, user.Email, user.Password, user.Bio, user.PrivilegeLevel, user.Payments, user.Background).Scan(&user.ID)
		} else {
			user.Avatar = "1.jpg"
			user.Background = "2.jpg"
			query = `INSERT INTO users (yid, username, email, password, bio, privilege_level, payments, background)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id;`
			err = db.QueryRow(query, user.Yid, user.Username, user.Email, user.Password, user.Bio, user.PrivilegeLevel, user.Payments, user.Background).Scan(&user.ID)
		}

	}
	return err
}

func GetUser(db *sql.DB, id int, requestingUserID int) (*User, error) {
	var user User
	query := `
		SELECT u.id, u.yid, u.username, u.email, u.avatar, u.bio, u.privilege_level, u.payments, u.block, u.background,
		       (SELECT COUNT(*) FROM user_followers uf WHERE uf.followee_id = u.id) AS subscribers_count,
		       EXISTS (SELECT 1 FROM user_followers uf WHERE uf.follower_id = $2 AND uf.followee_id = u.id) AS is_subscribed
		FROM users u
		WHERE u.id = $1`

	err := db.QueryRow(query, id, requestingUserID).Scan(&user.ID, &user.Yid, &user.Username, &user.Email, &user.Avatar, &user.Bio, &user.PrivilegeLevel, &user.Payments, &user.IsBlock, &user.Background, &user.SubscribersCount, &user.IsSubscribed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if user.IsBlock {
		return nil, ErrUserBlocked
	}
	return &user, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	email = strings.ReplaceAll(email, "%40", "@")
	var user User
	query := `
		SELECT u.id, u.yid, u.username, u.email, u.avatar, u.bio, u.privilege_level, u.payments, u.block, u.background,
		       (SELECT COUNT(*) FROM user_followers uf WHERE uf.followee_id = u.id) AS subscribers_count
		FROM users u
		WHERE u.email = $1`

	err := db.QueryRow(query, email).Scan(&user.ID, &user.Yid, &user.Username, &user.Email, &user.Avatar, &user.Bio, &user.PrivilegeLevel, &user.Payments, &user.IsBlock, &user.Background, &user.SubscribersCount)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAuthUser(db *sql.DB, us User) (*User, error) {
	var err error

	log.Printf("%v", us)
	query := `SELECT id, yid, privilege_level, block, background
              FROM users WHERE email = $1 AND password = $2`
	log.Print(us.Email, "`"+us.Password+"`")

	err = db.QueryRow(query, us.Email, us.Password).Scan(&us.ID, &us.Yid, &us.PrivilegeLevel, &us.IsBlock, &us.Background)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if us.IsBlock {
		return nil, ErrUserBlocked
	}
	log.Print(us)
	return &us, nil
}

func GetEmailUser(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT email
              FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsernameUser(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT username
              FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAvatarAndBackgroundUser(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT avatar, background
              FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.Avatar, &user.Background)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetPrivilegeLevelUser(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT privilege_level
              FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.PrivilegeLevel)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetBlockStatus(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT block
              FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.IsBlock)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(db *sql.DB, id int, updatedUser *User) error {
	// Начинаем с базового запроса
	query := `UPDATE users SET `
	var args []interface{}
	var setClauses []string

	// Проверяем каждое поле и добавляем его в запрос, если оно не пустое
	if updatedUser.Yid != "" {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Yid)
	}
	if updatedUser.Username != "" {
		setClauses = append(setClauses, "username = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Username)
	}
	if updatedUser.Email != "" {
		setClauses = append(setClauses, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Email)
	}
	if updatedUser.Password != "" {
		setClauses = append(setClauses, "password = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Password)
	}
	if updatedUser.Avatar != "" {
		setClauses = append(setClauses, "avatar = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Avatar)
	}
	if updatedUser.Bio != "" {
		setClauses = append(setClauses, "bio = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Bio)
	}
	if updatedUser.PrivilegeLevel != 0 { // Предполагаем, что 0 - это значение по умолчанию
		setClauses = append(setClauses, "privilege_level = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.PrivilegeLevel)
	}
	if updatedUser.Background != "" {
		setClauses = append(setClauses, "background = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Background)
	}

	// Если нет полей для обновления, возвращаем nil
	if len(setClauses) == 0 {
		return nil
	}

	// Объединяем все части запроса
	query += strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)

	// Выполняем запрос
	_, err := db.Exec(query, args...)
	return err
}

func DeleteUser(db *sql.DB, id int) error {
	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Удаление всех постов пользователя
	_, err = tx.Exec(`DELETE FROM posts WHERE id IN (SELECT post_id FROM user_posts WHERE user_id = $1)`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление всех связей пользователя с постами
	_, err = tx.Exec(`DELETE FROM user_posts WHERE user_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление всех связей пользователя с тегами
	_, err = tx.Exec(`DELETE FROM user_tags WHERE user_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление всех связей пользователя с подписками
	_, err = tx.Exec(`DELETE FROM user_subscriptions WHERE user_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление всех связей пользователя с комментариями
	_, err = tx.Exec(`DELETE FROM user_comments WHERE user_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление всех связей пользователя с подписчиками
	_, err = tx.Exec(`DELETE FROM user_followers WHERE follower_id = $1 OR followee_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление всех связей пользователя с IP-адресами
	_, err = tx.Exec(`DELETE FROM user_ips WHERE user_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление пользователя
	_, err = tx.Exec(`DELETE FROM users WHERE id = $1`, id)
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

func BlockUser(db *sql.DB, id int) error {
	query := `UPDATE users SET block = true WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func HashPassword(password string) string {
	// Создаем новый хеш
	hash := sha256.New()
	// Записываем пароль в хеш
	hash.Write([]byte(password))
	// Получаем хеш в виде байтового массива
	hashedBytes := hash.Sum(nil)
	// Преобразуем байты в строку в шестнадцатеричном формате
	return hex.EncodeToString(hashedBytes)
}

//func CheckPasswordHash(password, hash string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	return err == nil
//}
