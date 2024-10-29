package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"io"

	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/dsemenov12/shorturl/internal/config"
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
		{
            name: "negative test empty body",
			body: ``,
            want: want{
                code: 400,
        		contentType: "text/plain",
            },
        },
	}

	config.ParseFlags()

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
			defer res.Body.Close()
            
            assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
