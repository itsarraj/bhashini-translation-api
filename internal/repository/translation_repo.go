package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TranslationRepository handles translation cache operations
type TranslationRepository struct {
	db *sql.DB
}

// NewTranslationRepository creates a new translation repository
func NewTranslationRepository(db *sql.DB) *TranslationRepository {
	return &TranslationRepository{db: db}
}

// GetCachedTranslation retrieves a cached translation if it exists and hasn't expired
func (r *TranslationRepository) GetCachedTranslation(sourceText, sourceLang, targetLang string) (string, bool, error) {
	var translatedText string
	var expiresAt time.Time

	query := `
		SELECT translated_text, expires_at 
		FROM translation_cache 
		WHERE source_text = $1 
		AND source_lang = $2 
		AND target_lang = $3 
		AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, sourceText, sourceLang, targetLang).Scan(&translatedText, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return translatedText, true, nil
}

// CacheTranslation stores a translation in the cache
func (r *TranslationRepository) CacheTranslation(sourceText, sourceLang, targetLang, translatedText string, ttl time.Duration) error {
	id := uuid.New().String()
	expiresAt := time.Now().Add(ttl)

	query := `
		INSERT INTO translation_cache (id, source_text, source_lang, target_lang, translated_text, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6)
		ON CONFLICT (source_text, source_lang, target_lang) 
		DO UPDATE SET 
			translated_text = EXCLUDED.translated_text,
			created_at = NOW(),
			expires_at = EXCLUDED.expires_at
	`

	_, err := r.db.Exec(query, id, sourceText, sourceLang, targetLang, translatedText, expiresAt)
	return err
}

// CleanExpiredTranslations removes expired cache entries
func (r *TranslationRepository) CleanExpiredTranslations() error {
	query := `DELETE FROM translation_cache WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}
