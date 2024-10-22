package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/stretchr/testify/assert"
)

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
                code: 201,
        		contentType: "text/plain",
            },
        },
		{
            name: "positive test #2",
			body: `https://practicum.yandex.ru/2323`,
            want: want{
                code: 201,
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
            
            assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestRedirect(t *testing.T) {
	type want struct {
        code int
    }
	tests := []struct {
		name string
		want want
	}{
		{
            name: "positive test #1",
            want: want{
                code: 307,
            },
        },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			handlers.Redirect(response, request)

			res := response.Result()
            
            assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
