package main

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/KAnggara75/IDXStock/internal/database"
	"github.com/KAnggara75/IDXStock/internal/repositories"
	"github.com/KAnggara75/IDXStock/internal/routes"
	"github.com/KAnggara75/IDXStock/internal/services"
	"github.com/KAnggara75/scc2go"
)

func init() {
	scc2go.GetEnv(os.Getenv("SCC_URL"), os.Getenv("SCC_AUTH"))
}

func main() {
	// 1. Initialize Logger
	config.InitLogger()

	// 2. Initialize Database
	// For local testing, you might need to export DATABASE_URL
	// e.g., export DATABASE_URL=postgres://user:pass@localhost:5432/dbname
	database.Connect()

	// 2.1 Run Seeders
	brokerRepo := repositories.NewBrokerRepository(database.Pool)
	seederService := services.NewSeederService(brokerRepo, nil)
	if err := seederService.SeedBrokersData(context.Background()); err != nil {
		logrus.Errorf("Failed to run brokers seeder: %v", err)
	}

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "IDXStock API",
	})

	// 4. Middlewares
	app.Use(recover.New())

	// Custom Logger Middleware using logrus
	app.Use(func(c fiber.Ctx) error {
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
