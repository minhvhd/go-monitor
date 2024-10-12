package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/minhvhd/go-monitor/internal/hardware"
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
			data := make(map[string]interface{})

			systemInfo, err := hardware.GetSystem()
			if err == nil {
				data["system"] = systemInfo
			} else {
				fmt.Printf("Error getting system info: %v\n", err)
			}

			cpuInfo, err := hardware.GetCPU()
			if err == nil {
				data["cpu"] = cpuInfo
			} else {
				fmt.Printf("Error getting CPU info: %v\n", err)
			}

			diskInfo, err := hardware.GetDisk()
			if err == nil {
				data["disk"] = diskInfo
			} else {
				fmt.Printf("Error getting disk info: %v\n", err)
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("Error marshaling data: %v\n", err)
				continue
			}

			s.Broadcast(jsonData)

			time.Sleep(1 * time.Second)
		}
	}(srv)

	fmt.Println("Listening on port 8082")
	err := http.ListenAndServe(":8082", srv.GetMux())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
