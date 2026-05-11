package server

import (
	"database/sql"
	"fmt"
	"log"
	"mekoko/internal/config"
	"mekoko/internal/database"
	"mekoko/internal/middleware"
	"mekoko/internal/modules/auth"
	tokenGenerator "mekoko/internal/providers/tokens"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) (*gin.Engine, error) {
	db, err := connectPostgresWithRetries(cfg.DBUrl, 5, time.Second)
	if err != nil {
		return nil, err
	}
	r := gin.New()
	r.Use(middleware.RequestContextLogger())
	r.Use(gin.Recovery())

	r.HEAD("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	apiV1 := api.Group("/v1")

	generator := tokenGenerator.NewJWT(cfg.AccessSecret, cfg.RefreshSecret)

	authRepository := auth.NewRepository(db)
	authService := auth.NewService(authRepository, db, generator)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(apiV1, authHandler)

	return r, nil
}

func connectPostgresWithRetries(dsn string, attempts int, baseDelay time.Duration) (*sql.DB, error) {
	if attempts <= 0 {
		attempts = 1
	}

	if baseDelay <= 0 {
		baseDelay = time.Second
	}

	var lastErr error

	for attempt := 1; attempt < attempts; attempt++ {
		db, err := database.NewPostgresDB(dsn)
		if err == nil {
			return db, nil
		}

		lastErr = err
		if attempt == attempts {
			break
		}

		delay := baseDelay * time.Duration(attempt)
		log.Printf("database connection attempt %d/%d failed: %v (retrying in %s)", attempt, attempts, err, delay)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("database connection failed after %d attempts: %w", attempts, lastErr)

}
