package types

type Config struct {
	URL               string `json:"url"`
	BearerToken       string `json:"bearerToken"`
	Model             string `json:"model"`
	ChatInstructions  string `json:"chatInstructions"`
	ImageInstructions string `json:"imageInstructions"`
}

// ChatCompletion represents the overall response structure
type ChatCompletion struct {
	ID      string     `json:"id"`
	Object  string     `json:"object"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Choices []Choice   `json:"choices"`
	Usage   TokenUsage `json:"usage"`
}

// Choice represents an individual choice in the response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Message represents a message from the assistant
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// TokenUsage represents token usage details
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}