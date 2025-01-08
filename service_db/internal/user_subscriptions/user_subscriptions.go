package usersubscriptions

import "database/sql"

type UserSubscription struct {
	UserID             int `json:"user_id"`
	SubscriptionPlanID int `json:"subscription_plan_id"`
}

func CreateUserSubscription(db *sql.DB, subscription UserSubscription) error {
	query := `
	INSERT INTO user_subscriptions (user_id, subscription_plan_id)
	VALUES ($1, $2)
	`
	_, err := db.Exec(query, subscription.UserID, subscription.SubscriptionPlanID)
	if err != nil {
		return err
	}
	return nil
}

func GetUserSubscription(db *sql.DB, userID int) (*UserSubscription, error) {
	query := `
	SELECT user_id, subscription_plan_id
	FROM user_subscriptions
	WHERE user_id = $1
	`
	var subscription UserSubscription
	err := db.QueryRow(query, userID).Scan(&subscription.UserID, &subscription.SubscriptionPlanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}

func UpdateUserSubscription(db *sql.DB, subscription UserSubscription) error {
	query := `
	UPDATE user_subscriptions
	SET subscription_plan_id = $1
	WHERE user_id = $2
	`
	_, err := db.Exec(query, subscription.SubscriptionPlanID, subscription.UserID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserSubscription(db *sql.DB, userID int) error {
	query := `
	DELETE FROM user_subscriptions
	WHERE user_id = $1
	`
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
