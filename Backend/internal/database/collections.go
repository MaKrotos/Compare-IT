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
		INSERT INTO collections (user_id, name, is_pinned, public_link)
		VALUES ($1, $2, $3, NULL)
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
		SELECT id, user_id, name, is_pinned, public_link, created_at, updated_at
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
		var publicLink sql.NullString

		err := rows.Scan(
			&collection.ID,
			&collection.UserID,
			&collection.Name,
			&collection.IsPinned,
			&publicLink,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}

		// Устанавливаем публичную ссылку, если она есть
		if publicLink.Valid {
			collection.PublicLink = &publicLink.String
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
	// Сначала получаем текущую коллекцию из базы данных, чтобы сохранить публичную ссылку
	currentCollection, err := GetCollectionByID(collection.ID, collection.UserID)
	if err != nil {
		return fmt.Errorf("failed to get current collection: %w", err)
	}
	
	// Если у текущей коллекции есть публичная ссылка, сохраняем её
	var publicLink *string
	if currentCollection != nil && currentCollection.PublicLink != nil {
		publicLink = currentCollection.PublicLink
	}
	
	query := `
		UPDATE collections
		SET name = $1, is_pinned = $2, public_link = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
	`

	result, err := GetDB().Exec(query,
		collection.Name,
		collection.IsPinned,
		publicLink, // Используем сохраненную публичную ссылку
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

// GetCollectionByPublicLink получает подборку по публичной ссылке
func GetCollectionByPublicLink(publicLink string) (*models.Collection, error) {
	// Сначала пытаемся найти коллекцию по точному совпадению ссылки
	query := `
		SELECT id, user_id, name, is_pinned, public_link, created_at, updated_at
		FROM collections
		WHERE public_link = $1
	`

	var collection models.Collection
	var createdAt, updatedAt time.Time
	var publicLinkValue sql.NullString

	err := GetDB().QueryRow(query, publicLink).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.Name,
		&collection.IsPinned,
		&publicLinkValue,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		// Если не нашли по точному совпадению, проверяем, может быть ссылка передана как часть URL
		if err == sql.ErrNoRows {
			// Пытаемся найти коллекцию, у которой ссылка совпадает с концом переданной строки
			query = `
				SELECT id, user_id, name, is_pinned, public_link, created_at, updated_at
				FROM collections
				WHERE public_link = $1 OR $1 LIKE '%' || public_link
			`
			
			err = GetDB().QueryRow(query, publicLink).Scan(
				&collection.ID,
				&collection.UserID,
				&collection.Name,
				&collection.IsPinned,
				&publicLinkValue,
				&createdAt,
				&updatedAt,
			)
			
			// Если все равно не нашли, возвращаем nil
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, nil
				}
				return nil, fmt.Errorf("failed to get collection by public link: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get collection by public link: %w", err)
		}
	}

	// Устанавливаем публичную ссылку, если она есть
	if publicLinkValue.Valid {
		collection.PublicLink = &publicLinkValue.String
	}

	collection.CreatedAt = createdAt
	collection.UpdatedAt = updatedAt

	return &collection, nil
}

// GeneratePublicLinkForCollection создает публичную ссылку для коллекции
func GeneratePublicLinkForCollection(collectionID int) (string, error) {
	// Попытаемся сгенерировать уникальную ссылку, проверяя на дубликаты
	var link string
	var maxAttempts = 5
	var attempt = 0
	
	for attempt < maxAttempts {
		// Генерируем уникальную ссылку (можно использовать UUID или хеш)
		link = fmt.Sprintf("public_%d_%d_%d", collectionID, time.Now().Unix(), attempt)
		
		// Проверяем, существует ли уже такая ссылка
		var existingID int
		checkQuery := `SELECT id FROM collections WHERE public_link = $1`
		err := GetDB().QueryRow(checkQuery, link).Scan(&existingID)
		
		if err != nil {
			if err == sql.ErrNoRows {
				// Ссылка уникальна, можно использовать
				break
			}
			// Другая ошибка базы данных
			return "", fmt.Errorf("failed to check for duplicate public link: %w", err)
		}
		
		// Ссылка уже существует, попробуем другую
		attempt++
	}
	
	// Если не удалось сгенерировать уникальную ссылку за maxAttempts попыток
	if attempt >= maxAttempts {
		return "", fmt.Errorf("failed to generate unique public link after %d attempts", maxAttempts)
	}
	
	query := `
		UPDATE collections
		SET public_link = $1
		WHERE id = $2
	`
	
	result, err := GetDB().Exec(query, link, collectionID)
	if err != nil {
		return "", fmt.Errorf("failed to generate public link: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}
	
	return link, nil
}

// RemovePublicLinkForCollection удаляет публичную ссылку для коллекции
func RemovePublicLinkForCollection(collectionID int) error {
	query := `
		UPDATE collections
		SET public_link = NULL
		WHERE id = $1
	`
	
	result, err := GetDB().Exec(query, collectionID)
	if err != nil {
		return fmt.Errorf("failed to remove public link: %w", err)
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