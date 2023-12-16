package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/data", data)
	go func() {
		for {
			err := filepath.Walk(directory, visit)
			if err != nil {
				fmt.Printf("Error when crawling folders: %v\n", err)
			}

			interval := time.Second * 30
			time.Sleep(interval)
		}
	}()
	err := http.ListenAndServe(PortForServer, mux)
	if err != nil {
		log.Println(err)
	}
}
