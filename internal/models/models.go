package models

// InputData представляет входные данные с URL для сокращения.
type InputData struct {
	URL string `json:"url"`
}

// ResultJSON содержит результат операции по сокращению URL.
type ResultJSON struct {
	Result string `json:"result"`
}

// BatchItem представляет элемент запроса пакетной обработки URL.
type BatchItem struct {
	CorrelationID string `json:"correlation_id"` // Уникальный идентификатор корреляции
	OriginalURL   string `json:"original_url"`   // Исходный URL
}

// BatchResultItem содержит результат пакетной обработки URL.
type BatchResultItem struct {
	CorrelationID string `json:"correlation_id"` // Уникальный идентификатор корреляции
	ShortURL      string `json:"short_url"`      // Сокращенный URL
}

// ShortURLItem представляет связь между сокращенным и исходным URL.
type ShortURLItem struct {
	ShortURL    string `json:"short_url"`    // Сокращенный URL
	OriginalURL string `json:"original_url"` // Исходный URL
}

// StatsResponse представляет JSON-ответ для эндпоинта /api/internal/stats.
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
