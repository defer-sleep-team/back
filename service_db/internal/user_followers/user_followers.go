package userfollowers

import "database/sql"

type UserFollower struct {
	FollowerID int `json:"follower_id"`
	FolloweeID int `json:"followee_id"`
}

func InsertUserFollower(db *sql.DB, follower UserFollower) error {
	query := `
	INSERT INTO user_followers (follower_id, followee_id)
	VALUES ($1, $2)
	`
	_, err := db.Exec(query, follower.FollowerID, follower.FolloweeID)
	if err != nil {
		return err
	}
	return nil
}

func GetUserFollowers(db *sql.DB, followerID int) ([]UserFollower, error) {
	query := `
	SELECT follower_id, followee_id FROM user_followers WHERE follower_id = $1
	`
	rows, err := db.Query(query, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []UserFollower
	for rows.Next() {
		var follower UserFollower
		err := rows.Scan(&follower.FollowerID, &follower.FolloweeID)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return followers, nil
}

func GetUserFollowees(db *sql.DB, followeeID int) ([]UserFollower, error) {
	query := `
	SELECT follower_id, followee_id FROM user_followers WHERE followee_id = $1
	`
	rows, err := db.Query(query, followeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followees []UserFollower
	for rows.Next() {
		var followee UserFollower
		err := rows.Scan(&followee.FollowerID, &followee.FolloweeID)
		if err != nil {
			return nil, err
		}
		followees = append(followees, followee)
	}
	return followees, nil
}

func UpdateUserFollower(db *sql.DB, currentFollowerID, currentFolloweeID, newFolloweeID int) error {
	query := `
	UPDATE user_followers SET followee_id = $1 WHERE follower_id = $2 AND followee_id = $3
	`
	_, err := db.Exec(query, newFolloweeID, currentFollowerID, currentFolloweeID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserFollower(db *sql.DB, followerID, followeeID int) error {
	query := `
	DELETE FROM user_followers WHERE follower_id = $1 AND followee_id = $2
	`
	_, err := db.Exec(query, followerID, followeeID)
	if err != nil {
		return err
	}
	return nil
}
