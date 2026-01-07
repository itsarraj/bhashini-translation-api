package router

import (
	"database/sql"
	"user-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api/v1/translation")

	// Translation routes
	api.Post("/translate", handlers.Translate(db))
	api.Post("/cache/clean", handlers.CleanCache(db))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
			"service": "translation-service",
		})
	})
}

