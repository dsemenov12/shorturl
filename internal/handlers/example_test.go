package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/storage/memory"
)

func ExampleApp_ShortenPost() {
	store := memory.NewStorage()
	app := handlers.NewApp(store)
	reqBody := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	app.ShortenPost(res, req)

	respBody := res.Body.String()
	var result models.ResultJSON
	json.Unmarshal([]byte(respBody), &result)

	fmt.Println(res.Code)

	// Output:
	// 201
}

func ExampleApp_ShortenBatchPost() {
	store := memory.NewStorage()
	app := handlers.NewApp(store)

	reqBody := `[
		{"correlation_id": "1", "original_url": "https://example.com"},
		{"correlation_id": "2", "original_url": "https://golang.org"}
	]`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	app.ShortenBatchPost(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Body.String())

	// Output:
	// 201
	// [
	//     {
	//         "correlation_id": "1",
	//         "short_url": "/1"
	//     },
	//     {
	//         "correlation_id": "2",
	//         "short_url": "/2"
	//     }
	// ]
}

func ExampleApp_PostURL() {
	store := memory.NewStorage()
	app := handlers.NewApp(store)

	reqBody := "https://example.com"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	res := httptest.NewRecorder()

	app.PostURL(res, req)

	fmt.Println(res.Code)

	// Output:
	// 201
}

func ExampleApp_Redirect() {
	store := memory.NewStorage()
	store.Set(nil, "abc123", "https://example.com")
	app := handlers.NewApp(store)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	res := httptest.NewRecorder()

	app.Redirect(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Header().Get("Location"))

	// Output:
	// 307
	// /
}

func ExampleApp_UserUrls() {
	store := memory.NewStorage()
	app := handlers.NewApp(store)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	res := httptest.NewRecorder()

	app.UserUrls(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Body.String())

	// Output:
	// 200
	// null
}

func ExampleApp_DeleteUserUrls() {
	store := memory.NewStorage()
	app := handlers.NewApp(store)

	reqBody := `[
		"http://localhost:8080/1",
		"http://localhost:8080/2"
	]`
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	app.DeleteUserUrls(res, req)

	fmt.Println(res.Code)

	// Output:
	// ok
	// ok
	// 202
}
