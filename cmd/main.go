package main

import (
	"fmt"
	"keerja-backend/internal/config"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	config.InitLogger()
	appLogger := config.GetLogger()
	appLogger.Info("Starting Keerja Backend API...")

	if err := config.InitDatabase(); err != nil {
		appLogger.WithError(err).Fatal("Failed to initialize database")
	}
	defer config.CloseDatabase()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":", port)
	appLogger.Info("Server running on " + addr)

	if err := app.Listen(addr); err != nil {
		appLogger.WithError(err).Fatal("Failed to start")
	}
}
