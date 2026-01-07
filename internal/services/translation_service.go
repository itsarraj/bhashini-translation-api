package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"user-service/internal/repository"
)

// TranslationService handles translation business logic with caching
type TranslationService struct {
	bhashiniClient *BhashiniClient
	cacheRepo      *repository.TranslationRepository
	defaultPipelineID string
	cacheTTL       time.Duration
}

// NewTranslationService creates a new translation service
func NewTranslationService(bhashiniClient *BhashiniClient, cacheRepo *repository.TranslationRepository) *TranslationService {
	// Default pipeline ID for translation (can be overridden via env)
	pipelineID := os.Getenv("BHASHINI_PIPELINE_ID")
	if pipelineID == "" {
		// Use a known valid pipeline ID for translation
		// Valid IDs: 64392f96daac500b55c543cd (Initial), 660f813c0413087224435d2c (IIT Bombay), 660f866443e53d4133f65317 (IIIT Hyderabad)
		pipelineID = "64392f96daac500b55c543cd" // Initial Pipeline - supports Translation
	}

	// Parse cache TTL from env (default 24 hours)
	cacheTTL := 24 * time.Hour
	if ttlStr := os.Getenv("TRANSLATION_CACHE_TTL"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			cacheTTL = parsed
		}
	}

	return &TranslationService{
		bhashiniClient:    bhashiniClient,
		cacheRepo:        cacheRepo,
		defaultPipelineID: pipelineID,
		cacheTTL:         cacheTTL,
	}
}

// Translate translates text from source language to target language with caching
func (s *TranslationService) Translate(sourceText, sourceLang, targetLang string) (string, error) {
	// Normalize input
	sourceText = strings.TrimSpace(sourceText)
	if sourceText == "" {
		return "", fmt.Errorf("source text cannot be empty")
	}

	// Check cache first
	if cached, found, err := s.cacheRepo.GetCachedTranslation(sourceText, sourceLang, targetLang); err == nil && found {
		return cached, nil
	} else if err != nil {
		// Log error but continue with API call
		fmt.Printf("Cache lookup error: %v\n", err)
	}

	// Get pipeline config
	config, err := s.bhashiniClient.GetPipelineConfig(s.defaultPipelineID, sourceLang, targetLang)
	if err != nil {
		// If pipeline ID fails, try to find a valid one
		if pipelineID, findErr := s.bhashiniClient.FindTranslationPipeline(); findErr == nil {
			s.defaultPipelineID = pipelineID
			config, err = s.bhashiniClient.GetPipelineConfig(s.defaultPipelineID, sourceLang, targetLang)
		}
		if err != nil {
			return "", fmt.Errorf("failed to get pipeline config: %w", err)
		}
	}

	// Perform translation
	response, err := s.bhashiniClient.Translate(config, sourceText, sourceLang, targetLang)
	if err != nil {
		return "", fmt.Errorf("failed to translate: %w", err)
	}

	// Extract translation from pipeline response
	if len(response.PipelineResponse) == 0 {
		return "", fmt.Errorf("no pipeline response received")
	}
	
	// Find translation task output
	var translatedText string
	for _, pipelineItem := range response.PipelineResponse {
		if pipelineItem.TaskType == "translation" && len(pipelineItem.Output) > 0 {
			translatedText = pipelineItem.Output[0].Target
			break
		}
	}
	
	if translatedText == "" {
		return "", fmt.Errorf("no translation output received")
	}

	// Cache the translation
	if err := s.cacheRepo.CacheTranslation(sourceText, sourceLang, targetLang, translatedText, s.cacheTTL); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Cache storage error: %v\n", err)
	}

	return translatedText, nil
}

// GenerateCacheKey generates a unique cache key for translation
func (s *TranslationService) GenerateCacheKey(sourceText, sourceLang, targetLang string) string {
	data := fmt.Sprintf("%s:%s:%s", sourceText, sourceLang, targetLang)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CleanExpiredCache removes expired cache entries
func (s *TranslationService) CleanExpiredCache() error {
	return s.cacheRepo.CleanExpiredTranslations()
}
