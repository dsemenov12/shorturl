package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/models"
	mock_storage "github.com/dsemenov12/shorturl/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestShortenBatchPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mock_storage.NewMockStorage(ctrl)

	m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	app := NewApp(m)

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
			body: `[{"correlation_id": "JJUQVrJ12","original_url": "https://practicum.yandex.ru/"},{"correlation_id": "JJUQVrJ22","original_url": "https://mail.ru/"}]`,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name: "positive test #2",
			body: `[{"correlation_id": "JJUQVrJ12","original_url": "https://practicum.yandex.ru/123"},{"correlation_id": "JJUQVrJ22","original_url": "https://mail.ru/1234"}]`,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name: "test empty body",
			body: ``,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			app.ShortenBatchPost(response, request)

			res := response.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestShortenPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mock_storage.NewMockStorage(ctrl)

	m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	app := NewApp(m)

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
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name: "positive test #2",
			body: `{"url": "https://practicum.yandex.ru/2323"}`,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name: "test empty body",
			body: ``,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			app.ShortenPost(response, request)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mock_storage.NewMockStorage(ctrl)

	m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	app := NewApp(m)

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
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #2",
			body: `https://practicum.yandex.ru/2323`,
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "test empty body",
			body: ``,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			app.PostURL(response, request)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mock_storage.NewMockStorage(ctrl)

	m.EXPECT().Get(gomock.Any(), gomock.Any()).Return("bmXrsnZk", "https://practicum.yandex.ru/profile/go-advanced/", false, nil).AnyTimes()

	// создадим экземпляр приложения и передадим ему «хранилище»
	app := NewApp(m)

	type want struct {
		code        int
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
				code:        http.StatusTemporaryRedirect,
				redirectURL: "https://practicum.yandex.ru/profile/go-advanced/",
			},
		},
		{
			name: "positive test #2",
			code: "NVbvbWXj",
			want: want{
				code:        http.StatusTemporaryRedirect,
				redirectURL: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "positive test #3",
			code: "CztkzbdO",
			want: want{
				code:        http.StatusTemporaryRedirect,
				redirectURL: "https://practicum.yandex.ru/profile/",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestURL := config.FlagBaseAddr + "/" + test.code
			request := httptest.NewRequest(http.MethodGet, requestURL, nil)
			response := httptest.NewRecorder()

			app.Redirect(response, request)

			res := response.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_storage.NewMockStorage(ctrl)

	result := []models.ShortURLItem{}
	result = append(result, models.ShortURLItem{
		OriginalURL: "https://practicum.yandex.ru/",
		ShortURL:    "http://127.0.0.1:8080/qsd54gFg/gh5dEm34",
	})

	app := NewApp(m)

	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name: "test not found",
			want: want{
				code:        http.StatusNoContent,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.want.code == http.StatusNoContent {
				m.EXPECT().GetUserURL(gomock.Any()).Return(nil, errors.New("Not found"))
			} else {
				m.EXPECT().GetUserURL(gomock.Any()).Return(result, nil)
			}

			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			response := httptest.NewRecorder()

			app.UserUrls(response, request)

			res := response.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestDeleteUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_storage.NewMockStorage(ctrl)

	m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	app := NewApp(m)

	type want struct {
		code int
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive test #1",
			body: `["JJUQVrJ12", "JJUQVrJ22", "Jlfd67ds", "cdpFuzqh"]`,
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "test empty body",
			body: ``,
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			app.DeleteUserUrls(response, request)

			res := response.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
