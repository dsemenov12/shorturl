package grpcserver_test

import (
	"context"
	"testing"
	"time"

	"github.com/dsemenov12/shorturl/internal/grpcserver"
	"github.com/dsemenov12/shorturl/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRunGRPCServer_StartStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	// Используем ":0" чтобы ОС сама выбрала свободный порт
	grpcAddr := "localhost:0"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем сервер в горутине, чтобы он не блокировал тест
	go func() {
		err := grpcserver.RunGRPCServer(ctx, mockStorage, grpcAddr)
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)
}

func TestRunGateway_StartStop(t *testing.T) {
	grpcAddr := "localhost:0"
	httpAddr := "localhost:0"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		err := grpcserver.RunGateway(ctx, grpcAddr, httpAddr)
		if err != nil && err.Error() != "http: Server closed" {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	err := <-errCh
	assert.NoError(t, err)
}
