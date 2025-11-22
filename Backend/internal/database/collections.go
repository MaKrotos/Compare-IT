package database

import (
	"database/sql"
	"fmt"
	"time"

	"backend/internal/models"
)

// CreateCollection создает новую подборку
func CreateCollection(collection *models.Collection) error {
	query := `
		INSERT INTO collections (user_id, name, is_pinned)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt time.Time

	err := GetDB().QueryRow(query,
		collection.UserID,
		collection.Name,
		collection.IsPinned,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	collection.ID = id
	collection.CreatedAt = createdAt
	collection.UpdatedAt = updatedAt

	return nil
}

// GetCollectionsByUserID получает все подборки пользователя
func GetCollectionsByUserID(userID int) ([]models.Collection, error) {
	query := `
		SELECT id, user_id, name, is_pinned, created_at, updated_at
		FROM collections
		WHERE user_id = $1
		ORDER BY is_pinned DESC, created_at DESC
	`

	rows, err := GetDB().Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	defer rows.Close()

	var collections []models.Collection

	for rows.Next() {
		var collection models.Collection
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&collection.ID,
			&collection.UserID,
			&collection.Name,
			&collection.IsPinned,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}

		collection.CreatedAt = createdAt
		collection.UpdatedAt = updatedAt

		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return collections, nil
}

// UpdateCollection обновляет подборку
func UpdateCollection(collection *models.Collection) error {
	query := `
		UPDATE collections
		SET name = $1, is_pinned = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
	`

	result, err := GetDB().Exec(query,
		collection.Name,
		collection.IsPinned,
		collection.ID,
		collection.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Обновляем время обновления
	collection.UpdatedAt = time.Now()

	return nil
}

// DeleteCollection удаляет подборку
func DeleteCollection(collectionID, userID int) error {
	query := `
		DELETE FROM collections
		WHERE id = $1 AND user_id = $2
	`

	result, err := GetDB().Exec(query, collectionID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetCollectionByID получает подборку по ID
func GetCollectionByID(collectionID, userID int) (*models.Collection, error) {
	query := `
		SELECT id, user_id, name, is_pinned, created_at, updated_at
		FROM collections
		WHERE id = $1 AND user_id = $2
	`

	var collection models.Collection
	var createdAt, updatedAt time.Time

	err := GetDB().QueryRow(query, collectionID, userID).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.Name,
		&collection.IsPinned,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	collection.CreatedAt = createdAt
	collection.UpdatedAt = updatedAt

	return &collection, nil
}