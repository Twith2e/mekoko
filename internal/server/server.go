package server

import (
	"mekoko/internal/config"
	"net/http"
)

func NewServer(cfg config.Config) (*http.Server, error) {
	router, err := NewRouter(cfg)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	return s, nil
}
