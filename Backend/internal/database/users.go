package database

import (
	"database/sql"
	"fmt"

	"backend/internal/models"
)

// CreateUser создает нового пользователя в базе данных
func CreateUser(user *models.TelegramUser) error {
	query := `
		INSERT INTO telegram_users (telegram_id, first_name, last_name, username, generated_name, is_admin, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt string
	
	err := GetDB().QueryRow(query, 
		user.TelegramID, 
		user.FirstName, 
		user.LastName, 
		user.Username, 
		user.GeneratedName, 
		user.IsAdmin, 
		user.Role.Int(),
	).Scan(&id, &createdAt, &updatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	
	return nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func GetUserByTelegramID(telegramID int64) (*models.TelegramUser, error) {
	query := `
		SELECT id, telegram_id, first_name, last_name, username, generated_name, is_admin, role, created_at, updated_at, last_active
		FROM telegram_users
		WHERE telegram_id = $1
	`

	user := &models.TelegramUser{}
	
	err := GetDB().QueryRow(query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.GeneratedName,
		&user.IsAdmin,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastActive,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user by telegram_id: %w", err)
	}
	
	return user, nil
}

// UpdateUserLastActive обновляет время последней активности пользователя
func UpdateUserLastActive(userID int, lastActive string) error {
	query := `
		UPDATE telegram_users
		SET last_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	
	_, err := GetDB().Exec(query, lastActive, userID)
	if err != nil {
		return fmt.Errorf("failed to update user last_active: %w", err)
	}
	
	return nil
}