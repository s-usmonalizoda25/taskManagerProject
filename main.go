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
	"github.com/s-usmonalizoda25/taskManagerProject/router"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load("config/config.env"); err != nil {
		log.Fatalf("Error loading config.env file: %v", err)
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
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully!")

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	appRouter := router.NewRouter(taskHandler)

	log.Println("Server is running on port :8080...")
	if err := http.ListenAndServe(":8080", appRouter); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
