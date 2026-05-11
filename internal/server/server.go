package server

import (
	"mekoko/internal/config"
	"net/http"
	"time"
)

func NewServer(cfg config.Config) (*http.Server, error) {
	router, err := NewRouter(cfg)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	return s, nil
}
