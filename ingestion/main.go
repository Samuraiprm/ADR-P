package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adr-p/ingestion/handlers"
	"github.com/adr-p/ingestion/middleware"
	"github.com/adr-p/ingestion/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize Redis
	redisClient := redis.NewClient(
		redis.WithAddr(os.Getenv("REDIS_ADDR")),
		redis.WithPassword(os.Getenv("REDIS_PASSWORD")),
		redis.WithDB(0),
	)

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.Metrics())

	// Health check
	router.GET("/healthz", handlers.HealthCheck(redisClient))

	// Metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API endpoints
	v1 := router.Group("/api/v1")
	{
		v1.POST("/events", handlers.HandleEvent(redisClient))
	}

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	redisClient.Close()
	log.Println("Server exiting")
}
