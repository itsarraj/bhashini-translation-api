package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"user-service/internal/models"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BhashiniClient handles communication with Bhashini API
type BhashiniClient struct {
	BaseURL    string
	UserID     string
	APIKey     string
	HTTPClient *http.Client
}

// NewBhashiniClient creates a new Bhashini client
func NewBhashiniClient() *BhashiniClient {
	baseURL := os.Getenv("BHASHINI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://meity-auth.ulcacontrib.org"
	}

	userID := os.Getenv("BHASHINI_USER_ID")
	apiKey := os.Getenv("BHASHINI_API_KEY")

	// Trim whitespace in case there are any spaces
	if userID != "" {
		userID = strings.TrimSpace(userID)
	}
	if apiKey != "" {
		apiKey = strings.TrimSpace(apiKey)
	}

	return &BhashiniClient{
		BaseURL:    baseURL,
		UserID:     userID,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SearchPipelines searches for available pipelines that support translation
func (c *BhashiniClient) SearchPipelines() (*models.PipelineSearchResponse, error) {
	if c.UserID == "" || c.APIKey == "" {
		return nil, errors.New("BHASHINI_USER_ID and BHASHINI_API_KEY must be set")
	}

	req := models.PipelineSearchRequest{
		TaskType: []string{"translation"},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/ulca/apis/v0/model/getModelsPipeline", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("userID", c.UserID)
	httpReq.Header.Set("ulcaApiKey", c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp models.PipelineSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		// If unmarshal fails, try to find pipeline ID in the response
		// Some APIs return pipeline info directly
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &searchResp, nil
}

// FindTranslationPipeline finds a valid pipeline ID for translation
// Returns the first known valid pipeline ID for translation
func (c *BhashiniClient) FindTranslationPipeline() (string, error) {
	// Known valid pipeline IDs for translation (from Bhashini documentation)
	// These are pre-validated pipeline IDs that support translation
	validPipelines := []string{
		"64392f96daac500b55c543cd", // Initial Pipeline (supports Translation, ASR, Transliteration, TTS)
		"660f813c0413087224435d2c", // IIT Bombay (Translation)
		"660f866443e53d4133f65317", // IIIT Hyderabad (Translation)
	}

	// Return the first one (Initial Pipeline is most comprehensive)
	if len(validPipelines) > 0 {
		return validPipelines[0], nil
	}

	return "", errors.New("no valid translation pipeline found")
}

// GetPipelineConfig retrieves pipeline configuration for translation
func (c *BhashiniClient) GetPipelineConfig(pipelineID, sourceLang, targetLang string) (*models.PipelineConfigResponse, error) {
	if c.UserID == "" {
		return nil, errors.New("BHASHINI_USER_ID is not set or empty")
	}
	if c.APIKey == "" {
		return nil, errors.New("BHASHINI_API_KEY is not set or empty. Please get your API key from https://bhashini.gov.in/ulca/dashboard")
	}
	if len(c.APIKey) < 20 {
		return nil, fmt.Errorf("API key seems too short (%d chars) - verify you're using the correct ulcaApiKey from dashboard", len(c.APIKey))
	}

	req := models.PipelineConfigRequest{
		PipelineTasks: []models.PipelineTask{
			{TaskType: "translation"},
		},
		PipelineRequestConfig: models.PipelineRequestConfig{
			PipelineID: pipelineID,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/ulca/apis/v0/model/getModelsPipeline", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("userID", c.UserID)
	httpReq.Header.Set("ulcaApiKey", c.APIKey)

	// Debug logging (only in development - can be removed or made conditional)
	// fmt.Printf("[DEBUG] Making request with UserID: %s, APIKey length: %d\n", c.UserID, len(c.APIKey))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Provide more helpful error message for 400 errors
		if resp.StatusCode == 400 {
			return nil, fmt.Errorf("API returned status %d: %s. Verify: 1) API key is the 'ulcaApiKey' from 'My Profile' section (not other keys), 2) Key is active/enabled in dashboard, 3) No extra spaces or quotes in .env file", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var configResp models.PipelineConfigResponse
	if err := json.Unmarshal(body, &configResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &configResp, nil
}

// Translate performs translation using Bhashini API
func (c *BhashiniClient) Translate(config *models.PipelineConfigResponse, sourceText, sourceLang, targetLang string) (*models.PipelineComputeResponse, error) {
	if config.PipelineInferenceAPIEndPoint.CallbackURL == "" {
		return nil, errors.New("callback URL not found in pipeline config")
	}

	// Extract service ID from config
	// Find the translation task config
	var serviceID string
	for _, taskConfig := range config.PipelineResponseConfig {
		if taskConfig.TaskType == "translation" && len(taskConfig.Config) > 0 {
			// Find config matching source and target language
			for _, cfg := range taskConfig.Config {
				if cfg.Language.SourceLanguage == sourceLang && cfg.Language.TargetLanguage == targetLang {
					serviceID = cfg.ServiceID
					break
				}
			}
			if serviceID == "" && len(taskConfig.Config) > 0 {
				// Fallback to first config item if exact match not found
				serviceID = taskConfig.Config[0].ServiceID
			}
			break
		}
	}

	if serviceID == "" {
		return nil, errors.New("could not find service ID for translation task")
	}

	req := models.PipelineComputeRequest{
		PipelineTasks: []models.PipelineComputeTask{
			{
				TaskType: "translation",
				Config: models.PipelineComputeTaskConfig{
					Language: models.TaskLanguage{
						SourceLanguage: sourceLang,
						TargetLanguage: targetLang,
					},
					ServiceID: serviceID,
				},
			},
		},
		InputData: models.InputData{
			Input: []models.InputItem{
				{Source: sourceText},
			},
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", config.PipelineInferenceAPIEndPoint.CallbackURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	// Set auth headers from pipeline config
	authKeyName := config.PipelineInferenceAPIEndPoint.InferenceAPIKey.Name
	authKeyValue := config.PipelineInferenceAPIEndPoint.InferenceAPIKey.Value
	if authKeyName != "" && authKeyValue != "" {
		httpReq.Header.Set(authKeyName, authKeyValue)
	}

	// Debug logging (only in development - can be removed or made conditional)
	// fmt.Printf("[DEBUG] Translation request URL: %s\n", config.PipelineInferenceAPIEndPoint.CallbackURL)
	// fmt.Printf("[DEBUG] Translation request payload: %s\n", string(jsonData))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s. Request payload was: %s", resp.StatusCode, string(body), string(jsonData))
	}

	var computeResp models.PipelineComputeResponse
	if err := json.Unmarshal(body, &computeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &computeResp, nil
}
