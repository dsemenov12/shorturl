package storage

import (
	"sync"
)

type Storage struct {
	mx sync.Mutex
    Data map[string]string
}

func (s *Storage) Get(key string) (string, error) {
	s.mx.Lock()
    defer s.mx.Unlock()
    return s.Data[key], nil
}

func (s *Storage) Set(key string, value string) {
	s.mx.Lock()
    defer s.mx.Unlock()
	s.Data[key] = value
}