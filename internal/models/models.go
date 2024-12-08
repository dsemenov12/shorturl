package models

type InputData struct {
	URL string `json:"url"`
}

type ResultJSON struct {
	Result string `json:"result"`
}

type BatchItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchResultItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type ShortUrlItem struct {
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}