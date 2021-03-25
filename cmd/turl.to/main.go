package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vimoppa/turl.to/internal/config"
	"github.com/vimoppa/turl.to/internal/router"
	"github.com/vimoppa/turl.to/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "[turl.to] ", 0)

	cfg, err := config.SetupConfigurationDefaults()
	if err != nil {
		log.Fatal(err)
	}

	s, err := storage.New(&cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		MaxHeaderBytes: 10, // 10 MB
		Addr:           ":" + cfg.Server.Port,
		WriteTimeout:   time.Second * time.Duration(cfg.Server.Timeout),
		ReadTimeout:    time.Second * time.Duration(cfg.Server.Timeout),
		IdleTimeout:    time.Second * time.Duration(cfg.Server.Timeout),
		Handler:        router.New(s),
	}

	logger.Printf("listening on %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
