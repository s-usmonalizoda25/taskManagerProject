package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/s-usmonalizoda25/taskManagerProject/handlers"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/repository"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/service"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/logger"
	"github.com/s-usmonalizoda25/taskManagerProject/router"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	mainLog, err := logger.New(true)
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	defer mainLog.Sync()

	if err := godotenv.Load("config/config.env"); err != nil {
		mainLog.Fatal("Error loading config.env file", zap.Error(err))
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		mainLog.Fatal("failed to connect to database", zap.Error(err))
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		mainLog.Fatal("failed to migrate database", zap.Error(err))
	}
	mainLog.Info("Database migration completed successfully!")

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, mainLog)
	taskHandler := handlers.NewTaskHandler(taskService, mainLog)

	appRouter := router.NewRouter(taskHandler)

	mainLog.Info("Server is running on port :8080...")
	if err := http.ListenAndServe(":8080", appRouter); err != nil {
		mainLog.Fatal("failed to start server", zap.Error(err))
	}
}
