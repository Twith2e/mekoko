package server

import (
	"database/sql"
	"fmt"
	"log"
	"mekoko/internal/config"
	"mekoko/internal/database"
	"mekoko/internal/middleware"
	"mekoko/internal/modules/auth"
	"mekoko/internal/modules/cart"
	"mekoko/internal/modules/order"
	"mekoko/internal/modules/product"
	"mekoko/internal/modules/waitlist"
	"mekoko/internal/providers/email"
	tokenGenerator "mekoko/internal/providers/tokens"
	"mekoko/internal/providers/upload"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) (*gin.Engine, error) {
	db, err := connectPostgresWithRetries(cfg.DBUrl, 5, time.Second)
	if err != nil {
		return nil, err
	}
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		AllowCredentials: true,
	}))
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

	isProd := cfg.IsProd

	// resend := email.NewResend(cfg.ResendApiKey)
	brevo := email.NewBrevo(cfg.BrevoApiKey)
	cloudinary, err := upload.NewCloudinary(cfg.CloudinaryCloudName, cfg.CloudinaryApiKey, cfg.CloudinaryApiSecret)
	if err != nil {
		log.Printf("failed to create cloudinary: %v", err)
		return nil, fmt.Errorf("failed to create cloudinary: %w", err)
	}

	authRepository := auth.NewRepository(db)
	authGuard := middleware.AuthGuard(generator, authRepository)
	adminGuard := middleware.AdminGuard()
	authService := auth.NewService(authRepository, db, generator, brevo, cfg.MekokoClientBaseURL, cfg.AppName)
	authHandler := auth.NewHandler(authService, stringToBool(isProd))
	adminHandler := auth.NewHandler(authService, stringToBool(isProd))
	auth.RegisterRoutes(apiV1, authGuard, adminGuard, authHandler, adminHandler)

	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository, db)
	productHandler := product.NewHandler(productService, cloudinary)
	adminProductHandler := product.NewAdminHandler(productService)
	product.RegisterRoutes(apiV1, authGuard, adminGuard, productHandler, adminProductHandler)

	cartRepository := cart.NewRepository(db)
	cartService := cart.NewService(cartRepository, db)
	cartHandler := cart.NewHandler(cartService)
	cart.RegisterRoutes(apiV1, authGuard, cartHandler)

	orderRepository := order.NewRepository(db)
	orderService := order.NewService(orderRepository, db)
	orderHandler := order.NewHandler(orderService)
	order.RegisterRoutes(apiV1, authGuard, orderHandler)

	waitlistRepository := waitlist.NewRepository(db)
	waitlistService := waitlist.NewService(waitlistRepository, brevo, cfg.AppName, cfg.EmailSender)
	waitlistHandler := waitlist.NewHandler(waitlistService)
	waitlist.AddRoute(apiV1, waitlistHandler)

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

func stringToBool(isProd string) bool {
	switch strings.ToLower(isProd) {
	case "false":
		return false
	case "true":
		return true
	default:
		return false
	}
}
