package userips

import (
	"database/sql"
	"time"
)

type UserIP struct {
	UserID    int       `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	LastLogin time.Time `json:"last_login"`
}

func СreateUserIpsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS user_ips (
		user_id INT NOT NULL,
		ip_address VARCHAR(45) NOT NULL,
		last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertUserIP(db *sql.DB, userIP UserIP) error {
	query := `
	INSERT INTO user_ips (user_id, ip_address, last_login)
	VALUES ($1, $2, $3)
	`
	_, err := db.Exec(query, userIP.UserID, userIP.IPAddress, userIP.LastLogin)
	if err != nil {
		return err
	}
	return nil
}

func GetUserIPs(db *sql.DB, userID int) ([]UserIP, error) {
	query := `
	SELECT user_id, ip_address, last_login FROM user_ips WHERE user_id = $1
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIPs []UserIP
	for rows.Next() {
		var userIP UserIP
		err := rows.Scan(&userIP.UserID, &userIP.IPAddress, &userIP.LastLogin)
		if err != nil {
			return nil, err
		}
		userIPs = append(userIPs, userIP)
	}
	return userIPs, nil
}

func UpdateUserIP(db *sql.DB, userIP UserIP) error {
	query := `
	UPDATE user_ips SET ip_address = $1, last_login = $2 WHERE user_id = $3
	`
	_, err := db.Exec(query, userIP.IPAddress, userIP.LastLogin, userIP.UserID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserIP(db *sql.DB, userID int) error {
	query := `
	DELETE FROM user_ips WHERE user_id = $1
	`
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
