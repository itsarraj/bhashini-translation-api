-- Create translation_cache table for caching translations
CREATE TABLE IF NOT EXISTS translation_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_text TEXT NOT NULL,
    source_lang VARCHAR(10) NOT NULL,
    target_lang VARCHAR(10) NOT NULL,
    translated_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    UNIQUE(source_text, source_lang, target_lang)
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_translation_cache_lookup ON translation_cache(source_text, source_lang, target_lang, expires_at);
CREATE INDEX IF NOT EXISTS idx_translation_cache_expires ON translation_cache(expires_at);

-- Create function to automatically clean expired entries (optional, can be called manually)
CREATE OR REPLACE FUNCTION clean_expired_translations()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM translation_cache WHERE expires_at < NOW();
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
