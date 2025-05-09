package handlers

import (
	"context"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/rand"
	"github.com/dsemenov12/shorturl/internal/storage"
	"github.com/dsemenov12/shorturl/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer реализует методы для работы с сокращением URL через gRPC.
// Он использует хранилище для сохранения и извлечения данных о сокращенных URL.
type GRPCServer struct {
	proto.UnimplementedShortenerServiceServer                 // Встраиваемая структура для реализации интерфейса сервиса
	Storage                                   storage.Storage // Хранилище для хранения сокращенных URL
}

// Shorten сокращает один URL и возвращает сокращенную ссылку.
// Получает оригинальный URL, генерирует уникальный короткий ключ и сохраняет его в хранилище.
// Возвращает сокращенный URL или ошибку, если произошла ошибка при сохранении.
func (s *GRPCServer) Shorten(ctx context.Context, req *proto.ShortenRequest) (*proto.ShortenResponse, error) {
	shortKey := rand.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	_, err := s.Storage.Set(ctx, shortKey, req.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store the URL: %v", err)
	}

	return &proto.ShortenResponse{
		Result: shortURL,
	}, nil
}

// ShortenBatchPost обрабатывает пакет запросов на сокращение URL.
// Для каждого элемента в пакете генерируется уникальный короткий ключ и сохраняется в хранилище.
// Возвращает список сокращенных URL с привязанными корреляционными идентификаторами.
func (s *GRPCServer) ShortenBatchPost(ctx context.Context, req *proto.ShortenBatchRequest) (*proto.ShortenBatchResponse, error) {
	var responseItems []*proto.ShortenBatchResponseItem

	for _, item := range req.Items {
		shortKey := rand.RandStringBytes(8)
		shortURL := config.FlagBaseAddr + "/" + shortKey

		_, err := s.Storage.Set(ctx, shortKey, item.OriginalUrl)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to store the URL: %v", err)
		}

		responseItems = append(responseItems, &proto.ShortenBatchResponseItem{
			CorrelationId: item.CorrelationId,
			ShortUrl:      shortURL,
		})
	}

	return &proto.ShortenBatchResponse{
		Items: responseItems,
	}, nil
}

// Redirect обрабатывает запрос на перенаправление по сокращенному URL.
// Извлекает оригинальный URL из хранилища по идентификатору и возвращает его для перенаправления.
func (s *GRPCServer) Redirect(ctx context.Context, req *proto.RedirectRequest) (*proto.RedirectResponse, error) {
	redirectLink, _, _, err := s.Storage.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve the URL: %v", err)
	}

	return &proto.RedirectResponse{
		Url: redirectLink,
	}, nil
}

// UserUrls возвращает список URL, сохранённых пользователем.
func (s *GRPCServer) UserUrls(ctx context.Context, _ *proto.Empty) (*proto.UserUrlsResponse, error) {
	urls, err := s.Storage.GetUserURL(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve user URLs: %v", err)
	}

	if len(urls) == 0 {
		return &proto.UserUrlsResponse{Urls: []*proto.URL{}}, nil
	}

	var response []*proto.URL
	for _, u := range urls {
		response = append(response, &proto.URL{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}

	return &proto.UserUrlsResponse{
		Urls: response,
	}, nil
}

// DeleteUserUrls удаляет список сокращенных URL, предоставленных пользователем.
func (s *GRPCServer) DeleteUserUrls(ctx context.Context, req *proto.DeleteUserUrlsRequest) (*proto.Empty, error) {
	if len(req.ShortUrls) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no URLs provided to delete")
	}
	for _, shortURL := range req.ShortUrls {
		err := s.Storage.Delete(ctx, shortURL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to delete URL %s: %v", shortURL, err)
		}
	}

	return &proto.Empty{}, nil
}

// InternalStats возвращает статистику по количеству пользователей и URL.
func (s *GRPCServer) InternalStats(ctx context.Context, req *proto.Empty) (*proto.StatsResponse, error) {
	usersCount, err := s.Storage.CountUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users count: %v", err)
	}

	urlsCount, err := s.Storage.CountURLs(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get URLs count: %v", err)
	}

	return &proto.StatsResponse{
		Urls:  int64(urlsCount),
		Users: int64(usersCount),
	}, nil
}
