package main

import (
	"log"
	"net/http"

	"waterclip/api/internal/config"
	httpx "waterclip/api/internal/http"
)

func main() {
	cfg := config.Load()
	if err := http.ListenAndServe(cfg.Address, httpx.NewRouter(cfg.ProcessorBaseURL)); err != nil {
		log.Fatal(err)
	}
}
