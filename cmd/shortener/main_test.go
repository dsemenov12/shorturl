package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"io"

	"github.com/dsemenov12/shorturl/internal/structs/storage"
	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/dsemenov12/shorturl/internal/config"
)

func TestShortenPost(t *testing.T) {
	type want struct {
        code        int
        contentType string
    }
	tests := []struct {
		name string
		body string
		want want
	}{
		{
            name: "positive test #1",
			body: `{"url": "https://practicum.yandex.ru/"}`,
            want: want{
                code: http.StatusCreated,
        		contentType: "application/json",
            },
        },
		{
            name: "positive test #2",
			body: `{"url": "https://practicum.yandex.ru/2323"}`,
            want: want{
                code: http.StatusCreated,
        		contentType: "application/json",
            },
        },
		{
            name: "negative test empty body",
			body: ``,
            want: want{
                code: http.StatusBadRequest,
        		contentType: "application/json",
            },
        },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			handlers.ShortenPost(response, request)

			res := response.Result()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
            
            assert.Equal(t, test.want.code, res.StatusCode)
			if test.body != `` {
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
				assert.Contains(t, string(body), config.FlagBaseAddr)
			}
		})
	}
}

func TestPostURL(t *testing.T) {
	type want struct {
        code        int
        contentType string
    }
	tests := []struct {
		name string
		body string
		want want
	}{
		{
            name: "positive test #1",
			body: `https://practicum.yandex.ru/`,
            want: want{
                code: http.StatusCreated,
        		contentType: "text/plain",
            },
        },
		{
            name: "positive test #2",
			body: `https://practicum.yandex.ru/2323`,
            want: want{
                code: http.StatusCreated,
        		contentType: "text/plain",
            },
        },
		{
            name: "negative test empty body",
			body: ``,
            want: want{
                code: http.StatusBadRequest,
        		contentType: "text/plain",
            },
        },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			handlers.PostURL(response, request)

			res := response.Result()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
            
            assert.Equal(t, test.want.code, res.StatusCode)
			if test.body != `` {
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
				assert.Contains(t, string(body), config.FlagBaseAddr)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	var storageObj = storage.Storage{Data: make(map[string]string)}

	storageObj.Set("bmXrsnZk", "https://practicum.yandex.ru/profile/go-advanced/")
	storageObj.Set("NVbvbWXj", "https://practicum.yandex.ru/")
	storageObj.Set("CztkzbdO", "https://practicum.yandex.ru/profile/")

	type want struct {
        code int
		redirectURL string
    }
	tests := []struct {
		name string
		code string
		want want
	}{
		{
            name: "positive test #1",
			code: "bmXrsnZk",
            want: want{
                code: http.StatusTemporaryRedirect,
				redirectURL: "https://practicum.yandex.ru/profile/go-advanced/",
            },
        },
		{
            name: "positive test #2",
			code: "NVbvbWXj",
            want: want{
                code: http.StatusTemporaryRedirect,
				redirectURL: "https://practicum.yandex.ru/",
            },
        },
		{
            name: "positive test #3",
			code: "CztkzbdO",
            want: want{
                code: http.StatusTemporaryRedirect,
				redirectURL: "https://practicum.yandex.ru/profile/",
            },
        },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestURL := config.FlagBaseAddr + "/" + test.code
			request := httptest.NewRequest(http.MethodGet, requestURL, nil)
			response := httptest.NewRecorder()

			handlers.Redirect(response, request)

			res := response.Result()
			defer res.Body.Close()
            
            assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}