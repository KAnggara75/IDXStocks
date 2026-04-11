package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStocks/internal/config"
	"github.com/KAnggara75/IDXStocks/internal/database"
	"github.com/KAnggara75/IDXStocks/internal/routes"
)

func main() {
	// 1. Initialize Logger
	config.InitLogger()

	// 2. Initialize Database
	// For local testing, you might need to export DATABASE_URL
	// e.g., export DATABASE_URL=postgres://user:pass@localhost:5432/dbname
	database.Connect()

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "IDXStocks API",
	})

	// 4. Middlewares
	app.Use(recover.New())

	// Custom Logger Middleware using logrus
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		logrus.WithFields(logrus.Fields{
			"method":  c.Method(),
			"path":    c.Path(),
			"status":  c.Response().StatusCode(),
			"latency": stop.Sub(start).String(),
			"ip":      c.IP(),
		}).Info("HTTP Request")

		return err
	})

	// 5. Setup Routes
	routes.Setup(app)

	// 6. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logrus.Infof("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
