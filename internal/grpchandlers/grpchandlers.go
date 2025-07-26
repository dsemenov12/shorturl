package grpchandlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dsemenov12/shorturl/internal/auth"
	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/rand"
	"github.com/dsemenov12/shorturl/internal/storage"
	pb "github.com/dsemenov12/shorturl/proto"
)

// GRPCServer реализует gRPC сервер с методами для работы с сокращением URL.
type GRPCServer struct {
	pb.UnimplementedShortenerServiceServer
	storage storage.Storage
}

// NewGRPCServer создаёт новый экземпляр GRPCServer с указанным хранилищем.
func NewGRPCServer(storage storage.Storage) *GRPCServer {
	return &GRPCServer{storage: storage}
}

// PostURL генерирует короткий ключ для URL, сохраняет его в хранилище и возвращает сокращённый URL.
func (s *GRPCServer) PostURL(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	shortKey := rand.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	shortKeyResult, err := s.storage.Set(ctx, shortKey, req.Url)
	if err != nil {
		shortURL = config.FlagBaseAddr + "/" + shortKeyResult
	}

	return &pb.ShortenResponse{Result: shortURL}, nil
}

// ShortenPost дублирует логику PostURL, предоставляя альтернативный gRPC метод для сокращения URL.
func (s *GRPCServer) ShortenPost(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	return s.PostURL(ctx, req)
}

// ShortenBatchPost обрабатывает пакет запросов на сокращение URL.
// Для каждого элемента из входного списка создаёт короткий URL и возвращает список результатов.
func (s *GRPCServer) ShortenBatchPost(ctx context.Context, req *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	var items []*pb.ShortenBatchResponseItem
	for _, item := range req.Items {
		if item.CorrelationId == "" || item.OriginalUrl == "" {
			continue
		}

		shortURL := config.FlagBaseAddr + "/" + item.CorrelationId
		shortKeyResult, err := s.storage.Set(ctx, item.CorrelationId, item.OriginalUrl)
		if err != nil {
			shortURL = config.FlagBaseAddr + "/" + shortKeyResult
		}

		items = append(items, &pb.ShortenBatchResponseItem{
			CorrelationId: item.CorrelationId,
			ShortUrl:      shortURL,
		})
	}

	return &pb.ShortenBatchResponse{Items: items}, nil
}

// Redirect возвращает оригинальный URL по короткому ключу.
// Возвращает ошибку, если URL был удалён или отсутствует.
func (s *GRPCServer) Redirect(ctx context.Context, req *pb.RedirectRequest) (*pb.RedirectResponse, error) {
	url, _, isDeleted, err := s.storage.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, errors.New("url was deleted")
	}
	return &pb.RedirectResponse{Url: url}, nil
}

// UserUrls возвращает список всех URL, сохранённых пользователем.
// Пользователь определяется по userID, извлечённому из контекста.
func (s *GRPCServer) UserUrls(ctx context.Context, _ *pb.Empty) (*pb.UserUrlsResponse, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	urls, err := s.storage.GetUserURL(ctx)
	if err != nil {
		return nil, err
	}

	var pbUrls []*pb.URL
	for _, u := range urls {
		pbUrls = append(pbUrls, &pb.URL{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}

	return &pb.UserUrlsResponse{Urls: pbUrls}, nil
}

// DeleteUserUrls обрабатывает запрос на удаление списка коротких URL пользователя.
// Для каждого URL вызывает метод удаления из хранилища.
func (s *GRPCServer) DeleteUserUrls(ctx context.Context, req *pb.DeleteUserUrlsRequest) (*pb.Empty, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Можно запускать удаление в фоне или сразу делать синхронно
	for _, shortURL := range req.ShortUrls {
		s.storage.Delete(ctx, shortURL)
	}
	return &pb.Empty{}, nil
}

// InternalStats возвращает статистику по количеству сохранённых URL и пользователей.
func (s *GRPCServer) InternalStats(ctx context.Context, _ *pb.Empty) (*pb.StatsResponse, error) {
	countUrls, err := s.storage.CountURLs(ctx)
	if err != nil {
		return nil, err
	}
	countUsers, err := s.storage.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.StatsResponse{
		Urls:  int64(countUrls),
		Users: int64(countUsers),
	}, nil
}
