package models

// PreviewData представляет данные предварительного просмотра URL
type PreviewData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
}