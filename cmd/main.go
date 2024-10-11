package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/minhvhd/go-monitor/internal/hardware"
	"github.com/minhvhd/go-monitor/internal/server"
)

func main() {

	go func() {
		for {
			system, err := hardware.GetSystem()
			if err != nil {
				fmt.Println(err)
				continue
			}

			cpu, err := hardware.GetCPU()
			if err != nil {
				fmt.Println(err)
				continue
			}

			disk, err := hardware.GetDisk()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(system)
			fmt.Println(cpu)
			fmt.Println(disk)

			time.Sleep(3 * time.Second)
		}
	}()

	mux := http.NewServeMux()
	mutex := &sync.Mutex{}

	mux.Handle("/", http.FileServer(http.Dir("../web")))

	srv := server.NewServer(
		server.WithMessageBuffer(10),
		server.WithMux(mux),
		server.WithMutex(mutex),
	)

	fmt.Println("Listening on port 8081")
	err := http.ListenAndServe(":8081", srv.GetMux())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
