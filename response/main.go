package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adr-p/response/db"
	"github.com/adr-p/response/handlers"
	"github.com/adr-p/response/middleware"
	"github.com/adr-p/response/redis"
	"github.com/adr-p/response/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://adr:adr@localhost:5432/adr?sslmode=disable"
	}

	psqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer psqlDB.Close()

	if err := psqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(
		redis.WithAddr(os.Getenv("REDIS_ADDR")),
	)

	// Initialize services
	verdictService := services.NewVerdictService(redisClient, psqlDB)
	responseService := services.NewResponseService(psqlDB)

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		if err := psqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Metrics
	router.GET("/metrics", gin.WrapH(prometheus.Handler()))

	// API endpoints
	v1 := router.Group("/api/v1")
	{
		v1.GET("/rules", handlers.GetRules(psqlDB))
		v1.POST("/rules", handlers.CreateRule(psqlDB))
		v1.GET("/stats", handlers.GetStats(psqlDB))
		v1.POST("/telegram/callback", handlers.TelegramCallback(psqlDB))
	}

	// Start verdict consumer
	go verdictService.StartConsuming(context.Background(), responseService)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	verdictService.Stop()
	log.Println("Server exiting")
}
