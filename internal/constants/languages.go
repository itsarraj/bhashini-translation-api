package constants

// SupportedLanguages returns the list of supported language codes (ISO-639)
// These are used for frontend dropdowns and i18n localization
var SupportedLanguages = []string{"en", "hi", "mr", "ta", "te", "gu", "pa", "or", "ml"}

// LanguageNames maps language codes to their display names
var LanguageNames = map[string]string{
	"en": "English",
	"hi": "Hindi",
	"mr": "Marathi",
	"ta": "Tamil",
	"te": "Telugu",
	"gu": "Gujarati",
	"pa": "Punjabi",
	"or": "Odia",
	"ml": "Malayalam",
}

// IsValidLanguage checks if a language code is supported
func IsValidLanguage(langCode string) bool {
	for _, lang := range SupportedLanguages {
		if lang == langCode {
			return true
		}
	}
	return false
}
