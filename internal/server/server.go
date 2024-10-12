package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

type Subscriber struct {
	messages chan []byte
}

type Server struct {
	messageBuffer int
	mux           *http.ServeMux
	mutex         *sync.Mutex
	subscribers   map[*Subscriber]struct{}
}

type Option func(*Server)

func WithMessageBuffer(buffer int) Option {
	return func(s *Server) {
		s.messageBuffer = buffer
	}
}

func WithSubscribers(subscribers map[*Subscriber]struct{}) Option {
	return func(s *Server) {
		s.subscribers = subscribers
	}
}

func WithMux(mux *http.ServeMux) Option {
	return func(s *Server) {
		s.mux = mux
	}
}

func WithMutex(mutex *sync.Mutex) Option {
	return func(s *Server) {
		s.mutex = mutex
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		subscribers: make(map[*Subscriber]struct{}),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Getter methods
func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

func (s *Server) GetSubscribers() map[*Subscriber]struct{} {
	return s.subscribers
}

func (s *Server) GetMessageBuffer() int {
	return s.messageBuffer
}

func (s *Server) GetMutex() *sync.Mutex {
	return s.mutex
}

func (s *Server) Broadcast(message []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for sub := range s.subscribers {
		sub.messages <- message
	}
}

func (s *Server) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.Subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) addSubscriber(sub *Subscriber) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.subscribers[sub] = struct{}{}
	fmt.Printf("Added subscriber: %p\n", sub)
}

func (s *Server) Subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sub := &Subscriber{
		messages: make(chan []byte, s.messageBuffer),
	}
	s.addSubscriber(sub)

	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	}

	ws, err := websocket.Accept(w, r, options)
	if err != nil {
		if websocket.CloseStatus(err) == websocket.StatusPolicyViolation {
			http.Error(w, "WebSocket upgrade required", http.StatusUpgradeRequired)
		} else {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		}
		s.Unsubscribe(sub)
		return err
	}
	defer ws.Close(websocket.StatusNormalClosure, "")

	wsCtx := ws.CloseRead(ctx)

	for {
		select {
		case <-wsCtx.Done():
			fmt.Println("wsCtx.Done()")
			return wsCtx.Err()
		case msg := <-sub.messages:
			fmt.Println(string(msg))
			err := ws.Write(wsCtx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		}
	}
}

func (s *Server) Unsubscribe(sub *Subscriber) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.subscribers, sub)
}
