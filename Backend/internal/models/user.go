package models

import (
	"database/sql"
	"fmt"
)

// UserRole представляет роль пользователя
type UserRole int

const (
	RoleUser UserRole = iota
	RoleArchitect
	RoleAdmin
)

// String возвращает строковое представление роли
func (r UserRole) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleArchitect:
		return "architect"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// Int возвращает целочисленное представление роли
func (r UserRole) Int() int {
	return int(r)
}

// ParseRoleFromInt преобразует целое число в UserRole
func ParseRoleFromInt(roleInt int) UserRole {
	switch roleInt {
	case 0:
		return RoleUser
	case 1:
		return RoleArchitect
	case 2:
		return RoleAdmin
	default:
		return RoleUser
	}
}

// TelegramUser представляет пользователя Telegram
type TelegramUser struct {
	ID            int            `json:"id"`
	TelegramID    int64          `json:"telegram_id"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	Username      string         `json:"username"`
	GeneratedName string         `json:"generated_name"`
	IsAdmin       bool           `json:"is_admin"`
	Role          UserRole       `json:"role"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	LastActive    sql.NullString `json:"last_active"`
}

// String возвращает строковое представление пользователя
func (u *TelegramUser) String() string {
	return fmt.Sprintf("TelegramUser{ID: %d, TelegramID: %d, Username: %s, Role: %s}", 
		u.ID, u.TelegramID, u.Username, u.Role.String())
}