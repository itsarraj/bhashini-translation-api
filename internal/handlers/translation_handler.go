package handlers

import (
	"database/sql"
	"fmt"
	"user-service/internal/constants"
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

		// Validate language codes
		if !constants.IsValidLanguage(req.SourceLang) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "source_lang '" + req.SourceLang + "' is not supported",
			})
		}
		if !constants.IsValidLanguage(req.TargetLang) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "target_lang '" + req.TargetLang + "' is not supported",
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

// TranslateBatchItem represents a single translation item in batch request
type TranslateBatchItem struct {
	SourceText string `json:"source_text" validate:"required"`
	SourceLang string `json:"source_lang" validate:"required"`
	TargetLang string `json:"target_lang" validate:"required"`
}

// TranslateBatchRequest represents the batch translation request
type TranslateBatchRequest struct {
	Items []TranslateBatchItem `json:"items" validate:"required"`
}

// TranslateBatchResponse represents the batch translation response
type TranslateBatchResponse struct {
	SourceTexts     []string `json:"source_texts"`
	SourceLangs     []string `json:"source_langs"`
	TargetLangs     []string `json:"target_langs"`
	TranslatedTexts []string `json:"translated_texts"`
}

// TranslateBatch handles batch translation requests
func TranslateBatch(db *sql.DB) fiber.Handler {
	// translate multiple texts at once with individual language pairs
	// input:
	// {
	// 	"items": [
	// 		{"source_text": "Hello", "source_lang": "en", "target_lang": "hi"},
	// 		{"source_text": "World", "source_lang": "en", "target_lang": "hi"}
	// 	]
	// }
	// output:
	// {
	// 	"source_texts": ["Hello", "World"],
	// 	"source_langs": ["en", "en"],
	// 	"target_langs": ["hi", "hi"],
	// 	"translated_texts": ["नमस्ते", "दुनिया"]
	// }

	return func(c *fiber.Ctx) error {
		var req TranslateBatchRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "Invalid request body: " + err.Error(),
			})
		}

		// Validate required fields
		if len(req.Items) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "items array is required and cannot be empty",
			})
		}

		// Validate each item
		for i, item := range req.Items {
			if item.SourceText == "" || item.SourceLang == "" || item.TargetLang == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"error":  fmt.Sprintf("item[%d]: source_text, source_lang, and target_lang are required", i),
				})
			}

			// Validate language codes
			if !constants.IsValidLanguage(item.SourceLang) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"error":  fmt.Sprintf("item[%d]: source_lang '%s' is not supported", i, item.SourceLang),
				})
			}
			if !constants.IsValidLanguage(item.TargetLang) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"error":  fmt.Sprintf("item[%d]: target_lang '%s' is not supported", i, item.TargetLang),
				})
			}
		}

		// Initialize services
		bhashiniClient := services.NewBhashiniClient()
		cacheRepo := repository.NewTranslationRepository(db)
		translationService := services.NewTranslationService(bhashiniClient, cacheRepo)

		// Prepare response arrays
		sourceTexts := make([]string, len(req.Items))
		sourceLangs := make([]string, len(req.Items))
		targetLangs := make([]string, len(req.Items))
		translatedTexts := make([]string, len(req.Items))

		// Perform translation for each item
		for i, item := range req.Items {
			translatedText, err := translationService.Translate(item.SourceText, item.SourceLang, item.TargetLang)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error",
					"error":  fmt.Sprintf("item[%d]: %s", i, err.Error()),
				})
			}

			sourceTexts[i] = item.SourceText
			sourceLangs[i] = item.SourceLang
			targetLangs[i] = item.TargetLang
			translatedTexts[i] = translatedText
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": TranslateBatchResponse{
				SourceTexts:     sourceTexts,
				SourceLangs:     sourceLangs,
				TargetLangs:     targetLangs,
				TranslatedTexts: translatedTexts,
			},
		})
	}
}

// Languages returns the list of supported languages (ISO-639 codes)
// Used by frontend for language dropdown and i18n localization
func Languages(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   constants.SupportedLanguages,
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
