package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cadenmarchese/bottlecap/pkg/payload"
	"github.com/cadenmarchese/bottlecap/pkg/types"
)

// Client parses the user input provided, constructs a payload, sends, and then parses the response from LLM.
func Client(subcommand, argument string) (string, error) {
	config, err := loadConfig("config.json")
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	var requestPayload string
	var responseBody []byte

	failedToSendError := fmt.Errorf("failed to send request: %w", err)
	chatCompletionsApi := "v1/chat/completions"
	imageGenerationsApi := "v1/images/generations"

	if subcommand == "ask" {
		requestPayload, err = payload.CreateChatRequestPayload(config.Model, config.ChatInstructions, argument)
		if err != nil {
			return "", err
		}

		responseBody, err = sendRequest("POST", fmt.Sprintf("%s/%s", config.URL, chatCompletionsApi), requestPayload)
		if err != nil {
			return "", failedToSendError
		}

	} else if subcommand == "image" {
		requestPayload, err = payload.CreateImageRequestPayload(config.Model, config.ImageInstructions, argument)
		if err != nil {
			return "", err
		}

		responseBody, err = sendRequest("POST", fmt.Sprintf("%s/%s", config.URL, chatCompletionsApi), requestPayload)
		if err != nil {
			return "", failedToSendError
		}

	} else if subcommand == "generate" {
		requestPayload, err = payload.CreateImageGenerationPayload(config.Model, argument)
		if err != nil {
			return "", err
		}

		responseBody, err = sendRequest("POST", fmt.Sprintf("%s/%s", config.URL, imageGenerationsApi), requestPayload)
		if err != nil {
			return "", failedToSendError
		}

	}

	return parseResponse(responseBody)
}

// loadConfig reads and parses the config.json file
func loadConfig(filename string) (*types.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config types.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// sendRequest sends an HTTP request to the LLM server
func sendRequest(method, url, payload string) ([]byte, error) {
	config, err := loadConfig("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if config.BearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.BearerToken))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

// parseResponse extracts meaningful content from the response JSON
func parseResponse(body []byte) (string, error) {
	var chat types.ChatCompletion
	if err := json.Unmarshal(body, &chat); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	if len(chat.Choices) == 0 {
		return "", fmt.Errorf("no content found in response: %v", chat)
	}

	choice := chat.Choices[0]
	return fmt.Sprintf(
		"%s\n\nCompletion Tokens: %d\nPrompt Tokens: %d\nTotal Tokens: %d",
		choice.Message.Content, chat.Usage.CompletionTokens, chat.Usage.PromptTokens, chat.Usage.TotalTokens,
	), nil
}
