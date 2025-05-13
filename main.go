package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// Config represents the structure of config.json
type Config struct {
	URL          string `json:"url"`
	BearerToken  string `json:"bearerToken"`
	Instructions string `json:"instructions"`
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

const subcommandAsk = "ask"

func main() {
	app := &cli.Command{
		Name:  subcommandAsk,
		Usage: `Ask the local LLM a question by providing it as the first argument to the "ask" subcommand.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := os.Args[1:]

			if len(args) < 2 || args[0] != subcommandAsk {
				return errors.New(`Usage: ask "Your question in quotes"`)
			}

			question := args[1]
			response, err := client(question)
			if err != nil {
				return err
			}

			fmt.Println(response)
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// client interacts with the local LLM based on the config file
func client(input string) (string, error) {
	config, err := loadConfig("config.json")
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	requestPayload := createRequestPayload(config.Instructions, input)
	responseBody, err := sendRequest("POST", config.URL, requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	return parseResponse(responseBody)
}

// loadConfig reads and parses the config.json file
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// createRequestPayload constructs the JSON request body
func createRequestPayload(instruction, userInput string) string {
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"content": instruction, "role": "system"},
			{"content": userInput, "role": "user"},
		},
	}
	data, _ := json.Marshal(payload)
	return string(data)
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
	var chat ChatCompletion
	if err := json.Unmarshal(body, &chat); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	if len(chat.Choices) == 0 {
		return "", errors.New("no content found in response")
	}

	choice := chat.Choices[0]
	return fmt.Sprintf(
		"%s\n\nCompletion Tokens: %d\nPrompt Tokens: %d\nTotal Tokens: %d",
		choice.Message.Content, chat.Usage.CompletionTokens, chat.Usage.PromptTokens, chat.Usage.TotalTokens,
	), nil
}
