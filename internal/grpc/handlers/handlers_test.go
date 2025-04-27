package handlers

import (
	"context"
	"testing"

	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(ctx context.Context, key string, value string) (string, error) {
	args := m.Called(ctx, key, value)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Get(ctx context.Context, key string) (string, string, bool, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.String(1), args.Bool(2), args.Error(3)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorage) Bootstrap(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) CountURLs(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) CountUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return int(args.Get(0).(int64)), args.Error(1)
}

func (m *MockStorage) GetUserURL(ctx context.Context) ([]models.ShortURLItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.ShortURLItem), args.Error(1)
}

// TestShorten тестирует метод Shorten
func TestShorten(t *testing.T) {
	mockStorage := new(MockStorage)

	mockStorage.On("Set", mock.Anything, mock.Anything, mock.Anything).Return("short_url", nil)

	server := &GRPCServer{
		Storage: mockStorage,
	}

	req := &proto.ShortenRequest{
		Url: "http://example.com",
	}

	resp, err := server.Shorten(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockStorage.AssertExpectations(t)
}

// TestShorten_FailSet тестирует ошибку при вызове Set
func TestShorten_FailSet(t *testing.T) {
	mockStorage := new(MockStorage)

	mockStorage.On("Set", mock.Anything, mock.Anything, mock.Anything).Return("", assert.AnError)

	server := &GRPCServer{
		Storage: mockStorage,
	}

	req := &proto.ShortenRequest{
		Url: "http://example.com",
	}

	resp, err := server.Shorten(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to store the URL")
	mockStorage.AssertExpectations(t)
}

// TestShortenBatchPost тестирует метод ShortenBatchPost
func TestShortenBatchPost(t *testing.T) {
	mockStorage := new(MockStorage)

	mockStorage.On("Set", mock.Anything, mock.Anything, mock.Anything).Return("short_url", nil)

	server := &GRPCServer{
		Storage: mockStorage,
	}

	req := &proto.ShortenBatchRequest{
		Items: []*proto.ShortenBatchItem{
			{CorrelationId: "1", OriginalUrl: "http://example1.com"},
			{CorrelationId: "2", OriginalUrl: "http://example2.com"},
		},
	}

	resp, err := server.ShortenBatchPost(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 2)
	mockStorage.AssertExpectations(t)
}

// TestRedirect тестирует метод Redirect
func TestRedirect(t *testing.T) {
	mockStorage := new(MockStorage)
	server := &GRPCServer{Storage: mockStorage}

	mockStorage.On("Get", mock.Anything, mock.Anything).Return("http://example.com", "", true, nil)

	resp, err := server.Redirect(context.Background(), &proto.RedirectRequest{Id: "some_id"})
	assert.Nil(t, err)
	assert.Equal(t, "http://example.com", resp.Url)
	mockStorage.AssertExpectations(t)
}

// TestInternalStats проверяет корректность работы метода InternalStats
func TestInternalStats(t *testing.T) {
	mockStorage := new(MockStorage)

	mockStorage.On("CountUsers", mock.Anything).Return(int64(10), nil)
	mockStorage.On("CountURLs", mock.Anything).Return(int(100), nil)

	server := &GRPCServer{
		Storage: mockStorage,
	}

	response, err := server.InternalStats(context.Background(), &proto.Empty{})

	assert.NoError(t, err)
	assert.Equal(t, int64(10), response.Users)
	assert.Equal(t, int(100), int(response.Urls))
	mockStorage.AssertExpectations(t)
}
