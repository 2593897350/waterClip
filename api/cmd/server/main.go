package main

import (
	"log"
	"net/http"

	"waterclip/api/internal/config"
	httpx "waterclip/api/internal/http"
)

func main() {
	cfg := config.Default()
	if err := http.ListenAndServe(cfg.Address, httpx.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
