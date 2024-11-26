package models

type InputData struct {
	URL string `json:"url"`
}

type ResultJSON struct {
	Result string `json:"result"`
}

type BatchItem struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchResultItem struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type BatchResult struct {
	BatchItem []BatchResultItem
}