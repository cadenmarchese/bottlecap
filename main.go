package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// ChatCompletion represents the overall response structure
type ChatCompletion struct {
	ID      string     `json:"id"`
	Object  string     `json:"object"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Choices []Choice   `json:"choices"`
	Usage   TokenUsage `json:"usage"`
}

// Choice represents a choice in the response
type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	Logprobs     *Logprobs `json:"logprobs"`
	FinishReason string    `json:"finish_reason"`
}

// Message represents a message from the assistant
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// Logprobs represents log probability information (can be nil)
type Logprobs struct {
	// Define fields if needed, assuming null means no data provided
}

// TokenUsage represents token usage details
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func main() {
	cmd := &cli.Command{
		Name:  "ask",
		Usage: "ask the local LLM a question",
		Action: func(context.Context, *cli.Command) error {
			args := os.Args[1:]

			if len(args) > 1 {
				return fmt.Errorf("Error: multiple arguments not yet supported")
			}
			arg := args[0]

			out, err := client(arg)
			if err != nil {
				return err
			}

			fmt.Printf(out)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func client(input string) (output string, err error) {
	request := fmt.Sprintf(`{
		"messages": [
		  {
			"content": "You are a helpful assistant.",
			"role": "system"
		  },
		  {
			"content": "%s",
			"role": "user"
		  }
		]
	  }`, input)

	url := "http://localhost:58930/v1/chat/completions"
	method := "POST"
	payload := strings.NewReader(request)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var chat ChatCompletion
	err = json.Unmarshal([]byte(body), &chat)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	// Extract and print the content only
	if len(chat.Choices) > 0 {
		return chat.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("No content found.")
}
