package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/minhvhd/go-monitor/internal/server"
)

func main() {
	mux := http.NewServeMux()
	mutex := &sync.Mutex{}

	srv := server.NewServer(
		server.WithMessageBuffer(10),
		server.WithMux(mux),
		server.WithMutex(mutex),
	)

	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/ws", srv.SubscribeHandler)

	go func(s *server.Server) {
		for {
			s.Broadcast([]byte("Hello"))
			time.Sleep(5 * time.Second)
		}
	}(srv)

	fmt.Println("Listening on port 8082")
	err := http.ListenAndServe(":8082", srv.GetMux())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
