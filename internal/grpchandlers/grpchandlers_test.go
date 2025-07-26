package grpchandlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dsemenov12/shorturl/internal/auth"
	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/grpchandlers"
	"github.com/dsemenov12/shorturl/internal/models"
	mock_storage "github.com/dsemenov12/shorturl/internal/storage/mocks"
	pb "github.com/dsemenov12/shorturl/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCServer_PostURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	srv := grpchandlers.NewGRPCServer(mockStorage)

	mockStorage.EXPECT().
		Set(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("shortkey", nil).
		Times(1)

	req := &pb.ShortenRequest{Url: "https://example.com"}

	resp, err := srv.PostURL(context.Background(), req)

	assert.NoError(t, err)
	assert.Contains(t, resp.Result, config.FlagBaseAddr)
}

func TestGRPCServer_ShortenBatchPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	srv := grpchandlers.NewGRPCServer(mockStorage)

	req := &pb.ShortenBatchRequest{
		Items: []*pb.ShortenBatchItem{
			{CorrelationId: "id1", OriginalUrl: "https://a.com"},
			{CorrelationId: "id2", OriginalUrl: "https://b.com"},
			{CorrelationId: "", OriginalUrl: "https://skip.com"},
			{CorrelationId: "id4", OriginalUrl: ""},
		},
	}

	// для валидных элементов ожидаем вызовы Set
	mockStorage.EXPECT().Set(gomock.Any(), "id1", "https://a.com").Return("id1", nil)
	mockStorage.EXPECT().Set(gomock.Any(), "id2", "https://b.com").Return("id2", nil)

	resp, err := srv.ShortenBatchPost(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Items, 2)

	for _, item := range resp.Items {
		assert.Contains(t, item.ShortUrl, config.FlagBaseAddr)
		assert.NotEmpty(t, item.CorrelationId)
	}
}

func TestGRPCServer_Redirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	srv := grpchandlers.NewGRPCServer(mockStorage)

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().
			Get(gomock.Any(), "short").
			Return("https://example.com", "", false, nil)

		resp, err := srv.Redirect(context.Background(), &pb.RedirectRequest{Id: "short"})
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", resp.Url)
	})

	t.Run("deleted", func(t *testing.T) {
		mockStorage.EXPECT().
			Get(gomock.Any(), "deleted").
			Return("", "", true, nil)

		resp, err := srv.Redirect(context.Background(), &pb.RedirectRequest{Id: "deleted"})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "url was deleted")
	})

	t.Run("not found error", func(t *testing.T) {
		mockStorage.EXPECT().
			Get(gomock.Any(), "missing").
			Return("", "", false, errors.New("not found"))

		resp, err := srv.Redirect(context.Background(), &pb.RedirectRequest{Id: "missing"})
		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}

func TestGRPCServer_UserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	srv := grpchandlers.NewGRPCServer(mockStorage)

	userCtx := context.WithValue(context.Background(), auth.UserIDKey, "user1")

	t.Run("authenticated user with urls", func(t *testing.T) {
		mockStorage.EXPECT().
			GetUserURL(userCtx).
			Return([]models.ShortURLItem{
				{ShortURL: "short1", OriginalURL: "https://orig1"},
				{ShortURL: "short2", OriginalURL: "https://orig2"},
			}, nil)

		resp, err := srv.UserUrls(userCtx, &pb.Empty{})
		assert.NoError(t, err)
		assert.Len(t, resp.Urls, 2)
	})

	t.Run("unauthenticated user", func(t *testing.T) {
		resp, err := srv.UserUrls(context.Background(), &pb.Empty{})
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage.EXPECT().
			GetUserURL(userCtx).
			Return(nil, errors.New("db error"))

		resp, err := srv.UserUrls(userCtx, &pb.Empty{})
		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}

func TestGRPCServer_DeleteUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	srv := grpchandlers.NewGRPCServer(mockStorage)

	userCtx := context.WithValue(context.Background(), auth.UserIDKey, "user1")

	t.Run("authenticated user", func(t *testing.T) {
		mockStorage.EXPECT().Delete(userCtx, "short1").Return(nil)
		mockStorage.EXPECT().Delete(userCtx, "short2").Return(nil)

		req := &pb.DeleteUserUrlsRequest{
			ShortUrls: []string{"short1", "short2"},
		}

		resp, err := srv.DeleteUserUrls(userCtx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("unauthenticated user", func(t *testing.T) {
		req := &pb.DeleteUserUrlsRequest{
			ShortUrls: []string{"short1"},
		}

		resp, err := srv.DeleteUserUrls(context.Background(), req)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestGRPCServer_InternalStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	srv := grpchandlers.NewGRPCServer(mockStorage)

	mockStorage.EXPECT().CountURLs(gomock.Any()).Return(10, nil)
	mockStorage.EXPECT().CountUsers(gomock.Any()).Return(5, nil)

	resp, err := srv.InternalStats(context.Background(), &pb.Empty{})

	assert.NoError(t, err)
	assert.Equal(t, int64(10), resp.Urls)
	assert.Equal(t, int64(5), resp.Users)
}
