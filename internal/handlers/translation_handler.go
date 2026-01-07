package handlers

import (
	"database/sql"
	"user-service/internal/repository"
	"user-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

// TranslateRequest represents the translation request
type TranslateRequest struct {
	SourceText string `json:"source_text" validate:"required"`
	SourceLang string `json:"source_lang" validate:"required"`
	TargetLang string `json:"target_lang" validate:"required"`
}

// TranslateResponse represents the translation response
type TranslateResponse struct {
	SourceText     string `json:"source_text"`
	SourceLang     string `json:"source_lang"`
	TargetLang     string `json:"target_lang"`
	TranslatedText string `json:"translated_text"`
}

// Translate handles translation requests
func Translate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req TranslateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "Invalid request body: " + err.Error(),
			})
		}

		// Validate required fields
		if req.SourceText == "" || req.SourceLang == "" || req.TargetLang == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "source_text, source_lang, and target_lang are required",
			})
		}

		// Initialize services
		bhashiniClient := services.NewBhashiniClient()
		cacheRepo := repository.NewTranslationRepository(db)
		translationService := services.NewTranslationService(bhashiniClient, cacheRepo)

		// Perform translation
		translatedText, err := translationService.Translate(req.SourceText, req.SourceLang, req.TargetLang)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": TranslateResponse{
				SourceText:     req.SourceText,
				SourceLang:     req.SourceLang,
				TargetLang:     req.TargetLang,
				TranslatedText: translatedText,
			},
		})
	}
}

// CleanCache handles cache cleanup requests
func CleanCache(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cacheRepo := repository.NewTranslationRepository(db)

		if err := cacheRepo.CleanExpiredTranslations(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Expired cache entries cleaned",
		})
	}
}
