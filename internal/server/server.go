package server

import (
	"net/http"
	"sync"
)

type subscriber struct {
	messages chan []byte
}

type server struct {
	messageBuffer int
	mux           *http.ServeMux
	mutex         *sync.Mutex
	subscribers   map[*subscriber]struct{}
}

type Option func(*server)

func WithMessageBuffer(buffer int) Option {
	return func(s *server) {
		s.messageBuffer = buffer
	}
}

func WithSubscribers(subscribers map[*subscriber]struct{}) Option {
	return func(s *server) {
		s.subscribers = subscribers
	}
}

func WithMux(mux *http.ServeMux) Option {
	return func(s *server) {
		s.mux = mux
	}
}

func WithMutex(mutex *sync.Mutex) Option {
	return func(s *server) {
		s.mutex = mutex
	}
}

func NewServer(opts ...Option) *server {
	s := &server{
		subscribers: make(map[*subscriber]struct{}),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *server) GetMux() *http.ServeMux {
	return s.mux
}
