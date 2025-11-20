package preview

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"backend/internal/models"
)

// extractTextContent извлекает текстовое содержимое из HTML узла
func extractTextContent(n *html.Node) string {
	var sb strings.Builder

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.TrimSpace(sb.String())
}

// getFirstParagraph находит первый значимый параграф текста
func getFirstParagraph(doc *goquery.Document) string {
	// Ищем первый параграф с содержимым
	paragraph := ""
	doc.Find("p").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := strings.TrimSpace(s.Text())
		// Проверяем, что параграф содержит достаточно текста
		if len(text) > 50 {
			paragraph = text
			return false // останавливаем поиск
		}
		return true // продолжаем поиск
	})

	// Если не нашли параграф, ищем в div или span
	if paragraph == "" {
		doc.Find("div, span, h1, h2, h3").EachWithBreak(func(i int, s *goquery.Selection) bool {
			text := strings.TrimSpace(s.Text())
			// Проверяем, что элемент содержит достаточно текста
			if len(text) > 50 {
				// Берем только первые 200 символов
				if len(text) > 200 {
					text = text[:200] + "..."
				}
				paragraph = text
				return false // останавливаем поиск
			}
			return true // продолжаем поиск
		})
	}

	return paragraph
}

// tryDifferentUserAgents пытается получить данные с разными User-Agent заголовками
func tryDifferentUserAgents(url string) (*models.PreviewData, error) {
	// Список различных User-Agent заголовков
	userAgents := []string{
		"", // Пустой User-Agent
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"TelegramBot (like TwitterBot)",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}

	for _, userAgent := range userAgents {
		preview, err := fetchWithUserAgent(url, userAgent)
		if err == nil && preview.Title != "Доступ ограничен" && preview.Title != "Ошибка загрузки" {
			return preview, nil
		}
		// Если получили валидные данные, возвращаем их
		if err == nil {
			return preview, nil
		}
	}

	// Если ни один User-Agent не сработал, возвращаем ошибку
	return nil, fmt.Errorf("не удалось получить данные с любого User-Agent")
}

// fetchWithUserAgent получает данные с определенным User-Agent
func fetchWithUserAgent(url string, userAgent string) (*models.PreviewData, error) {
	// Создаем HTTP клиент с пользовательским агентом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Создаем запрос с пользовательским агентом
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки для запроса HTML только если User-Agent не пустой
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Cache-Control", "max-age=0")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		// Если получили 401 или 403, пробуем получить хотя бы мета-данные
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			// Возвращаем минимальные данные вместо ошибки
			return &models.PreviewData{
				Title:       "Доступ ограничен",
				Description: "Сайт не предоставляет публичный доступ к содержимому",
				Image:       "", // Пустое изображение
				URL:         url,
			}, nil
		}
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	preview := &models.PreviewData{
		URL: url,
	}

	// Получаем заголовок из тега title
	if title := doc.Find("title").Text(); title != "" {
		preview.Title = strings.TrimSpace(title)
	} else if ogTitle, exists := doc.Find("meta[property='og:title']").Attr("content"); exists {
		preview.Title = strings.TrimSpace(ogTitle)
	} else {
		preview.Title = "Без названия"
	}

	// Получаем описание из мета-тега description
	if desc, exists := doc.Find("meta[name='description']").Attr("content"); exists {
		preview.Description = strings.TrimSpace(desc)
	} else if ogDesc, exists := doc.Find("meta[property='og:description']").Attr("content"); exists {
		preview.Description = strings.TrimSpace(ogDesc)
	} else {
		// Если нет мета-описания, извлекаем первый значимый параграф
		preview.Description = getFirstParagraph(doc)
		if preview.Description == "" {
			preview.Description = "Описание отсутствует"
		}
	}

	// Получаем изображение из мета-тегов
	if img, exists := doc.Find("meta[property='og:image']").Attr("content"); exists {
		preview.Image = img
	} else if img, exists := doc.Find("meta[name='twitter:image']").Attr("content"); exists {
		preview.Image = img
	} else if img, exists := doc.Find("meta[property='og:image:url']").Attr("content"); exists {
		preview.Image = img
	} else if img, exists := doc.Find("link[rel='image_src']").Attr("href"); exists {
		preview.Image = img
	}

	return preview, nil
}

// GetPreviewData получает данные предварительного просмотра для URL
func GetPreviewData(url string) (*models.PreviewData, error) {
	// Пытаемся получить данные с разными User-Agent заголовками
	preview, err := tryDifferentUserAgents(url)
	if err != nil {
		// Если не удалось получить данные, возвращаем ошибку
		return nil, err
	}

	return preview, nil
}