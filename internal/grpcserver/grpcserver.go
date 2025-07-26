package grpcserver

import (
	"context"
	"net"
	"net/http"

	"github.com/dsemenov12/shorturl/internal/grpchandlers"
	"github.com/dsemenov12/shorturl/internal/middlewares/authinterceptor"
	"github.com/dsemenov12/shorturl/internal/middlewares/logger"
	"github.com/dsemenov12/shorturl/internal/storage"
	pb "github.com/dsemenov12/shorturl/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RunGRPCServer запускает gRPC сервер с указанным адресом и хранилищем.
// Внутри сервера используется Unary Interceptor для аутентификации пользователей с помощью JWT-токенов.
// После успешной аутентификации создаётся gRPC сервер и регистрируется обработчик сервиса ShortenerService.
//
// ctx: Контекст для управления жизненным циклом сервера (например, отмена через сигнал).
// storage: Реализация интерфейса Storage для работы с данными.
// grpcAddr: Адрес (host:port), на котором запускается gRPC сервер.
//
// Возвращаемое значение: ошибка запуска сервера (если есть).
func RunGRPCServer(ctx context.Context, storage storage.Storage, grpcAddr string) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(authinterceptor.AuthUnaryInterceptor()),
	)
	pb.RegisterShortenerServiceServer(grpcSrv, grpchandlers.NewGRPCServer(storage))

	go func() {
		<-ctx.Done()
		grpcSrv.GracefulStop()
	}()

	logger.Log.Info("Starting gRPC server", zap.String("address", grpcAddr))
	return grpcSrv.Serve(lis)
}

// RunGateway запускает HTTP сервер grpc-gateway, который проксирует REST-запросы в gRPC сервер.
//
// ctx: Контекст для управления жизненным циклом сервера.
// grpcAddr: Адрес gRPC сервера, к которому grpc-gateway будет подключаться.
// httpAddr: Адрес, на котором запускается HTTP сервер grpc-gateway.
//
// Возвращаемое значение: ошибка запуска HTTP сервера (если есть).
func RunGateway(ctx context.Context, grpcAddr, httpAddr string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterShortenerServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	logger.Log.Info("Starting grpc-gateway server", zap.String("address", httpAddr))
	return srv.ListenAndServe()
}
