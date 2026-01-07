package models

// Bhashini API Request/Response Models

// PipelineConfigRequest represents the request for Pipeline Config API
type PipelineConfigRequest struct {
	PipelineTasks         []PipelineTask        `json:"pipelineTasks"`
	PipelineRequestConfig PipelineRequestConfig `json:"pipelineRequestConfig"`
}

// PipelineTask represents a task in the pipeline
type PipelineTask struct {
	TaskType string `json:"taskType"` // e.g., "translation"
}

// PipelineRequestConfig contains pipeline configuration
type PipelineRequestConfig struct {
	PipelineID string `json:"pipelineId"`
}

// PipelineConfigResponse represents the response from Pipeline Config API
type PipelineConfigResponse struct {
	PipelineInferenceAPIEndPoint PipelineInferenceAPIEndPoint `json:"pipelineInferenceAPIEndPoint"`
	PipelineResponseConfig       []PipelineResponseConfigItem `json:"pipelineResponseConfig"`
}

// PipelineInferenceAPIEndPoint contains the endpoint and auth details
type PipelineInferenceAPIEndPoint struct {
	CallbackURL     string          `json:"callbackUrl"`
	InferenceAPIKey InferenceAPIKey `json:"inferenceApiKey"`
}

// InferenceAPIKey contains the auth key name and value
type InferenceAPIKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PipelineResponseConfigItem represents a single pipeline response config item
type PipelineResponseConfigItem struct {
	TaskType string       `json:"taskType"`
	Config   []ConfigItem `json:"config"`
}

// ConfigItem represents a configuration item
type ConfigItem struct {
	ServiceID string       `json:"serviceId"`
	ModelID   string       `json:"modelId"`
	Language  LanguagePair `json:"language"`
}

// LanguagePair represents source and target languages (API returns flat structure)
type LanguagePair struct {
	SourceLanguage   string `json:"sourceLanguage"`
	SourceScriptCode string `json:"sourceScriptCode,omitempty"`
	TargetLanguage   string `json:"targetLanguage"`
	TargetScriptCode string `json:"targetScriptCode,omitempty"`
}

// PipelineComputeRequest represents the request for Pipeline Compute API
type PipelineComputeRequest struct {
	PipelineTasks []PipelineComputeTask `json:"pipelineTasks"`
	InputData     InputData             `json:"inputData"`
}

// PipelineComputeTask represents a task with config for compute API
type PipelineComputeTask struct {
	TaskType string                    `json:"taskType"`
	Config   PipelineComputeTaskConfig `json:"config"`
}

// PipelineComputeTaskConfig contains task configuration
type PipelineComputeTaskConfig struct {
	Language  TaskLanguage `json:"language"`
	ServiceID string       `json:"serviceId"`
}

// TaskLanguage represents language configuration for a task
type TaskLanguage struct {
	SourceLanguage string `json:"sourceLanguage,omitempty"`
	TargetLanguage string `json:"targetLanguage,omitempty"`
}

// InputData represents input data for the pipeline
type InputData struct {
	Input []InputItem `json:"input"`
}

// InputItem represents input data for translation
type InputItem struct {
	Source string `json:"source"`
}

// PipelineComputeResponse represents the response from Pipeline Compute API
type PipelineComputeResponse struct {
	PipelineResponse []PipelineResponseItem `json:"pipelineResponse"`
}

// PipelineResponseItem represents a single pipeline response item
type PipelineResponseItem struct {
	TaskType string       `json:"taskType"`
	Config   interface{}  `json:"config"`
	Output   []OutputItem `json:"output"`
	Audio    interface{}  `json:"audio"`
}

// OutputItem represents output data from translation
type OutputItem struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// TranslationCache represents a cached translation entry
type TranslationCache struct {
	ID             string `json:"id"`
	SourceText     string `json:"source_text"`
	SourceLang     string `json:"source_lang"`
	TargetLang     string `json:"target_lang"`
	TranslatedText string `json:"translated_text"`
	CreatedAt      string `json:"created_at"`
	ExpiresAt      string `json:"expires_at"`
}

// PipelineSearchRequest represents the request for Pipeline Search API
type PipelineSearchRequest struct {
	TaskType []string `json:"taskType,omitempty"`
}

// PipelineSearchResponse represents the response from Pipeline Search API
type PipelineSearchResponse struct {
	Pipelines []PipelineInfo `json:"pipelines"`
}

// PipelineInfo represents pipeline information
type PipelineInfo struct {
	PipelineID string   `json:"pipelineId"`
	TaskType   []string `json:"taskType"`
}
