package subscriptionplans

import "database/sql"

type SubscriptionPlan struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func InsertSubscriptionPlan(db *sql.DB, plan SubscriptionPlan) error {
	query := `
	INSERT INTO subscription_plans (user_id, name, price)
	VALUES ($1, $2, $3)
	`
	_, err := db.Exec(query, plan.UserID, plan.Name, plan.Price)
	if err != nil {
		return err
	}
	return nil
}

func GetSubscriptionPlans(db *sql.DB, userID int) ([]SubscriptionPlan, error) {
	query := `
	SELECT id, user_id, name, price FROM subscription_plans WHERE user_id = $1
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []SubscriptionPlan
	for rows.Next() {
		var plan SubscriptionPlan
		err := rows.Scan(&plan.ID, &plan.UserID, &plan.Name, &plan.Price)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func UpdateSubscriptionPlan(db *sql.DB, plan SubscriptionPlan) error {
	query := `
	UPDATE subscription_plans SET name = $1, price = $2 WHERE id = $3
	`
	_, err := db.Exec(query, plan.Name, plan.Price, plan.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSubscriptionPlan(db *sql.DB, id int) error {
	query := `
	DELETE FROM subscription_plans WHERE id = $1
	`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
