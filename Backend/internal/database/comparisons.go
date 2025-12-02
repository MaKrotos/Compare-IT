package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/models"
)

// CreateComparisonCollection создает новую коллекцию сравнений
func CreateComparisonCollection(collection *models.ComparisonCollection) error {
	query := `
		INSERT INTO comparison_collections (user_id, name, items, public_link, price_rating_weight, pros_cons_rating_weight)
		VALUES ($1, $2, $3, NULL, $4, $5)
		RETURNING id, created_at, updated_at
	`

	itemsJSON, err := json.Marshal(collection.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	var id int
	var createdAt, updatedAt time.Time

	err = GetDB().QueryRow(query,
		collection.UserID,
		collection.Name,
		itemsJSON,
		collection.PriceRatingWeight,
		collection.ProsConsRatingWeight,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return fmt.Errorf("failed to create comparison collection: %w", err)
	}

	collection.ID = id
	collection.CreatedAt = createdAt
	collection.UpdatedAt = updatedAt

	return nil
}

// GetComparisonCollectionByID получает коллекцию сравнений по ID
func GetComparisonCollectionByID(collectionID, userID int) (*models.ComparisonCollection, error) {
	query := `
		SELECT id, user_id, name, items, public_link, price_rating_weight, pros_cons_rating_weight, created_at, updated_at
		FROM comparison_collections
		WHERE id = $1 AND user_id = $2
	`

	var collection models.ComparisonCollection
	var itemsJSON []byte
	var publicLink sql.NullString
	var createdAt, updatedAt time.Time

	err := GetDB().QueryRow(query, collectionID, userID).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.Name,
		&itemsJSON,
		&publicLink,
		&collection.PriceRatingWeight,
		&collection.ProsConsRatingWeight,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get comparison collection: %w", err)
	}

	// Парсим JSON с элементами
	err = json.Unmarshal(itemsJSON, &collection.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	// Устанавливаем публичную ссылку, если она есть
	if publicLink.Valid {
		collection.PublicLink = &publicLink.String
	}

	collection.CreatedAt = createdAt
	collection.UpdatedAt = updatedAt

	return &collection, nil
}

// GetComparisonCollectionsByUserID получает все коллекции сравнений пользователя
func GetComparisonCollectionsByUserID(userID int) ([]models.ComparisonCollection, error) {
	query := `
		SELECT id, user_id, name, items, public_link, price_rating_weight, pros_cons_rating_weight, created_at, updated_at
		FROM comparison_collections
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := GetDB().Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comparison collections: %w", err)
	}
	defer rows.Close()

	var collections []models.ComparisonCollection

	for rows.Next() {
		var collection models.ComparisonCollection
		var itemsJSON []byte
		var publicLink sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&collection.ID,
			&collection.UserID,
			&collection.Name,
			&itemsJSON,
			&publicLink,
			&collection.PriceRatingWeight,
			&collection.ProsConsRatingWeight,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comparison collection: %w", err)
		}

		// Парсим JSON с элементами
		err = json.Unmarshal(itemsJSON, &collection.Items)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
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

// UpdateComparisonCollection обновляет коллекцию сравнений
func UpdateComparisonCollection(collection *models.ComparisonCollection) error {
	// Сначала получаем текущую коллекцию из базы данных, чтобы сохранить публичную ссылку
	currentCollection, err := GetComparisonCollectionByID(collection.ID, collection.UserID)
	if err != nil {
		return fmt.Errorf("failed to get current comparison collection: %w", err)
	}
	
	// Если у текущей коллекции есть публичная ссылка, сохраняем её
	var publicLink *string
	if currentCollection != nil && currentCollection.PublicLink != nil {
		publicLink = currentCollection.PublicLink
	}
	
	query := `
		UPDATE comparison_collections
		SET name = $1, items = $2, public_link = $3, price_rating_weight = $4, pros_cons_rating_weight = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND user_id = $7
	`

	itemsJSON, err := json.Marshal(collection.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	result, err := GetDB().Exec(query,
		collection.Name,
		itemsJSON,
		publicLink, // Используем сохраненную публичную ссылку
		collection.PriceRatingWeight,
		collection.ProsConsRatingWeight,
		collection.ID,
		collection.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update comparison collection: %w", err)
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

// DeleteComparisonCollection удаляет коллекцию сравнений
func DeleteComparisonCollection(collectionID, userID int) error {
	query := `
		DELETE FROM comparison_collections
		WHERE id = $1 AND user_id = $2
	`
	
	result, err := GetDB().Exec(query, collectionID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete comparison collection: %w", err)
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

// GeneratePublicLink создает публичную ссылку для коллекции
func GeneratePublicLink(collectionID int) (string, error) {
	// Попытаемся сгенерировать уникальную ссылку, проверяя на дубликаты
	var link string
	var maxAttempts = 5
	var attempt = 0
	
	for attempt < maxAttempts {
		// Генерируем уникальную ссылку (можно использовать UUID или хеш)
		link = fmt.Sprintf("public_%d_%d_%d", collectionID, time.Now().Unix(), attempt)
		
		// Проверяем, существует ли уже такая ссылка
		var existingID int
		checkQuery := `SELECT id FROM comparison_collections WHERE public_link = $1`
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
		UPDATE comparison_collections
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

// generateCollectionHash генерирует хэш для коллекции на основе её содержимого
func generateCollectionHash(collection *models.ComparisonCollection) string {
	// Создаем строку, представляющую содержимое коллекции
	content := fmt.Sprintf("%s|%v", collection.Name, collection.Items)
	
	// Генерируем SHA256 хэш
	hash := sha256.Sum256([]byte(content))
	
	// Преобразуем хэш в строку hex
	return hex.EncodeToString(hash[:])
}

// RemovePublicLink удаляет публичную ссылку для коллекции
func RemovePublicLink(collectionID int) error {
	query := `
		UPDATE comparison_collections
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

// GetComparisonCollectionByPublicLink получает коллекцию по публичной ссылке
func GetComparisonCollectionByPublicLink(publicLink string) (*models.ComparisonCollection, error) {
	// Сначала пытаемся найти коллекцию по точному совпадению ссылки
	query := `
		SELECT id, user_id, name, items, public_link, price_rating_weight, pros_cons_rating_weight, created_at, updated_at
		FROM comparison_collections
		WHERE public_link = $1
	`
	
	var collection models.ComparisonCollection
	var itemsJSON []byte
	var publicLinkValue sql.NullString
	var createdAt, updatedAt time.Time
	
	err := GetDB().QueryRow(query, publicLink).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.Name,
		&itemsJSON,
		&publicLinkValue,
		&collection.PriceRatingWeight,
		&collection.ProsConsRatingWeight,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		// Если не нашли по точному совпадению, проверяем, может быть ссылка передана как часть URL
		if err == sql.ErrNoRows {
			// Пытаемся найти коллекцию, у которой ссылка совпадает с концом переданной строки
			query = `
				SELECT id, user_id, name, items, public_link, price_rating_weight, pros_cons_rating_weight, created_at, updated_at
				FROM comparison_collections
				WHERE public_link = $1 OR $1 LIKE '%' || public_link
			`
			
			err = GetDB().QueryRow(query, publicLink).Scan(
				&collection.ID,
				&collection.UserID,
				&collection.Name,
				&itemsJSON,
				&publicLinkValue,
				&collection.PriceRatingWeight,
				&collection.ProsConsRatingWeight,
				&createdAt,
				&updatedAt,
			)
			
			// Если все равно не нашли, возвращаем nil
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, nil
				}
				return nil, fmt.Errorf("failed to get comparison collection by public link: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get comparison collection by public link: %w", err)
		}
	}
	
	// Парсим JSON с элементами
	err = json.Unmarshal(itemsJSON, &collection.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}
	
	// Устанавливаем публичную ссылку, если она есть
	if publicLinkValue.Valid {
		collection.PublicLink = &publicLinkValue.String
	}
	
	// Генерируем хэш для коллекции
	collection.Hash = generateCollectionHash(&collection)
	
	collection.CreatedAt = createdAt
	collection.UpdatedAt = updatedAt
	
	return &collection, nil
}