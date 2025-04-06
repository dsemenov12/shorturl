package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputData_MarshalJSON(t *testing.T) {
	input := InputData{
		URL: "https://example.com",
	}

	// Тестируем маршалинг в JSON
	result, err := json.Marshal(input)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"url":"https://example.com"}`, string(result))
}

func TestResultJSON_MarshalJSON(t *testing.T) {
	result := ResultJSON{
		Result: "https://short.url",
	}

	// Тестируем маршалинг в JSON
	resultJSON, err := json.Marshal(result)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"result":"https://short.url"}`, string(resultJSON))
}

func TestBatchItem_MarshalJSON(t *testing.T) {
	batchItem := BatchItem{
		CorrelationID: "12345",
		OriginalURL:   "https://example.com",
	}

	// Тестируем маршалинг в JSON
	result, err := json.Marshal(batchItem)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"correlation_id":"12345","original_url":"https://example.com"}`, string(result))
}

func TestBatchResultItem_MarshalJSON(t *testing.T) {
	batchResult := BatchResultItem{
		CorrelationID: "12345",
		ShortURL:      "https://short.url/12345",
	}

	// Тестируем маршалинг в JSON
	result, err := json.Marshal(batchResult)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"correlation_id":"12345","short_url":"https://short.url/12345"}`, string(result))
}

func TestShortURLItem_MarshalJSON(t *testing.T) {
	shortURLItem := ShortURLItem{
		ShortURL:    "https://short.url",
		OriginalURL: "https://example.com",
	}

	// Тестируем маршалинг в JSON
	result, err := json.Marshal(shortURLItem)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"short_url":"https://short.url","original_url":"https://example.com"}`, string(result))
}

func TestInputData_UnmarshalJSON(t *testing.T) {
	inputJSON := `{"url":"https://example.com"}`

	var input InputData
	err := json.Unmarshal([]byte(inputJSON), &input)

	assert.Nil(t, err)
	assert.Equal(t, "https://example.com", input.URL)
}

func TestResultJSON_UnmarshalJSON(t *testing.T) {
	resultJSON := `{"result":"https://short.url"}`

	var result ResultJSON
	err := json.Unmarshal([]byte(resultJSON), &result)

	assert.Nil(t, err)
	assert.Equal(t, "https://short.url", result.Result)
}

func TestBatchItem_UnmarshalJSON(t *testing.T) {
	batchItemJSON := `{"correlation_id":"12345","original_url":"https://example.com"}`

	var batchItem BatchItem
	err := json.Unmarshal([]byte(batchItemJSON), &batchItem)

	assert.Nil(t, err)
	assert.Equal(t, "12345", batchItem.CorrelationID)
	assert.Equal(t, "https://example.com", batchItem.OriginalURL)
}

func TestBatchResultItem_UnmarshalJSON(t *testing.T) {
	batchResultJSON := `{"correlation_id":"12345","short_url":"https://short.url/12345"}`

	var batchResult BatchResultItem
	err := json.Unmarshal([]byte(batchResultJSON), &batchResult)

	assert.Nil(t, err)
	assert.Equal(t, "12345", batchResult.CorrelationID)
	assert.Equal(t, "https://short.url/12345", batchResult.ShortURL)
}

func TestShortURLItem_UnmarshalJSON(t *testing.T) {
	shortURLItemJSON := `{"short_url":"https://short.url","original_url":"https://example.com"}`

	var shortURLItem ShortURLItem
	err := json.Unmarshal([]byte(shortURLItemJSON), &shortURLItem)

	assert.Nil(t, err)
	assert.Equal(t, "https://short.url", shortURLItem.ShortURL)
	assert.Equal(t, "https://example.com", shortURLItem.OriginalURL)
}
