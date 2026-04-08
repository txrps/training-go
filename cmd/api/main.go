package main

import (
	"log"
	"net/http"
	"time"
	"training-go/internal/config"
	"training-go/internal/database"
	"training-go/internal/handlers"

	_ "training-go/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// migrate create -ext sql -dir migrations -seq create_testing_api_table

// @title           Go Training API
// @version         1.0
// @description     This is a training API using Gin + GORM + PostgreSQL
// @host            localhost:3000
// @BasePath        /

// @contact.name    API Support
// @contact.email   support@example.com

// @license.name    MIT
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("❌ Failed to load configuration:", err)
	}

	dbCfg := database.Config{
		DatabaseURL:    cfg.DatabaseURL,
		MaxConns:       25,
		MinConns:       5,
		RetryAttempts:  3,
		ConnectTimeout: 5 * time.Second,
		PingTimeout:    3 * time.Second,
	}

	pool, err := database.Connect(dbCfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)

	gormDB, err := database.NewGormFromPool(sqlDB)

	_ = gormDB

	if err != nil {
		log.Fatal("Failed to init GORM:", err)
	}

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/todos", handlers.CreateTodoHandlerGorm(gormDB))
	router.GET("/todos", handlers.GetAllTodosHandlerGorm(gormDB))
	router.GET("/todos/:id", handlers.GetTodoByIDHandlerGorm(gormDB))
	router.PUT("/todos/:id", handlers.UpdateTodoHandlerGorm(gormDB))
	router.DELETE("/todos/:id", handlers.DeleteTodoHandlerGorm(gormDB))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("🚀 Server running on port %s", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("❌ Server failed:", err)
	}
}
