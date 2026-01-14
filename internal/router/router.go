package router

import (
	"database/sql"
	"user-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "ok",
			"service": "translation-service",
		})
	})

	api := app.Group("/v1")

	// Translation routes
	api.Post("/translate", handlers.Translate(db))
	api.Post("/translate/batch", handlers.TranslateBatch(db)) // translate multiple texts at once
	api.Get("/languages", handlers.Languages(db))             // return the list of languages, like en,hi, all iso-639 codes from readme file

	// manage routes
	manage := api.Group("/manage")
	manage.Post("/cache/clean", handlers.CleanCache(db))

}

// clean cache every 8 hours using redis FLUSHALL
