package models

import (
	"time"
)

// ComparisonItem представляет элемент сравнения
type ComparisonItem struct {
	ID             string   `json:"id"`
	URL            string   `json:"url"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Images         []string `json:"images"`
	Price          float64  `json:"price"`
	Currency       string   `json:"currency"`
	Pros           []ProCon `json:"pros"`
	Cons           []ProCon `json:"cons"`
	Rating         float64  `json:"rating"`
	PriceRating    float64  `json:"priceRating"`
	ProsConsRating float64  `json:"prosConsRating"`
	CreatedDate    string   `json:"created_date"`
}

// ProCon представляет плюс или минус
type ProCon struct {
	Text   string `json:"text"`
	Impact int    `json:"impact"`
}

// ComparisonCollection представляет коллекцию сравнений пользователя
type ComparisonCollection struct {
	ID                int               `json:"id"`
	UserID            int               `json:"user_id"`
	Name              string            `json:"name"`
	Items             []ComparisonItem  `json:"items"`
	PublicLink        *string           `json:"public_link,omitempty"`
	PriceRatingWeight int               `json:"price_rating_weight"`
	ProsConsRatingWeight int             `json:"pros_cons_rating_weight"`
	Hash              string            `json:"hash,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}