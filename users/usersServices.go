package users

import (
	"context"
	"errors"
	"gop/db"
)

// Used for initial authentication
func GetDBUserIDForProvider(provider string, providerUserID string) (int, error) {
	var userId int

	query := `
        SELECT user_id
        FROM user_auth_providers
        WHERE provider_title = $1 AND provider_user_id = $2
        LIMIT 1;
    `

	row := db.DBPool.QueryRow(context.Background(), query, provider, providerUserID)
	err := row.Scan(&userId)
	if err != nil {
		return 0, errors.New("user not found for provider")
	}

	return userId, nil
}

func CreateDBUser(provider string, userIDFromProvider string, emailFromProvider string) (int, error) {
	var userId int

	// Try to find an existing user first
	querySelect := `
        SELECT id
        FROM users
        WHERE auth_provider_title = $1 AND auth_provider_user_id = $2 AND auth_provider_email = $3
        LIMIT 1;
    `
	err := db.DBPool.QueryRow(context.Background(), querySelect, provider, userIDFromProvider, emailFromProvider).Scan(&userId)
	if err == nil {
		// User already exists
		return userId, nil
	}

	// If not found, insert a new user
	queryInsert := `
        INSERT INTO users (auth_provider_title, auth_provider_user_id, auth_provider_email)
        VALUES ($1, $2, $3)
        RETURNING id;
    `
	err = db.DBPool.QueryRow(context.Background(), queryInsert, provider, userIDFromProvider, emailFromProvider).Scan(&userId)
	if err != nil {
		return 0, errors.New("failed to create user")
	}

	return userId, nil
}
