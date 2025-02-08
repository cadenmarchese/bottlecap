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

type Config struct {
	Url          string `json:"url"`
	Method       string `json:"method"`
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

// Choice represents a choice in the response
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
// TODO: implement token handling
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

var subcommandAsk = "ask"

func main() {
	cmd := &cli.Command{
		Name:  subcommandAsk,
		Usage: `ask the local LLM a question, in quotes, as the first argument to the "ask" subcommand`,
		Action: func(context.Context, *cli.Command) error {

			// collect the user input and pass it to client
			// TODO: migrate the user input handling to bufio
			args := os.Args[1:]
			if len(args) == 0 {
				return fmt.Errorf(`Error: no arguments provided`)
			}

			cmd := args[0]
			if cmd != subcommandAsk {
				return fmt.Errorf(`Error: subcommand '%s' not yet implemeted`, cmd)
			}

			// user should provide exactly 2 arguments - the subcommand, and the question in quotes
			if len(args) != 2 {
				return fmt.Errorf(`Please ask your question in quotes. For example, "Why is the sky blue?"`)
			}

			arg := args[1]
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

// client uses the config file to interact locally with the LLM
func client(input string) (output string, err error) {

	// Read from the config file
	file, err := os.Open("config.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)

	// construct the HTTP request from config file content
	request := fmt.Sprintf(`{
		"messages": [
		  {
			"content": "%s",
			"role": "system"
		  },
		  {
			"content": "%s",
			"role": "user"
		  }
		]
	  }`, config.Instructions, input)

	url := config.Url
	method := config.Method
	payload := strings.NewReader(request)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	// send request to the local LLM via provided URL
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

	// unmarshal response into ChatCompletion struct
	var chat ChatCompletion
	err = json.Unmarshal([]byte(body), &chat)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	// extract and print the content (answer) only
	if len(chat.Choices) > 0 {
		return fmt.Sprintf("%s\n\nCompletion Tokens: %v\nPrompt Tokens: %v\nTotal Tokens: %v",
			chat.Choices[0].Message.Content, chat.Usage.CompletionTokens,
			chat.Usage.PromptTokens, chat.Usage.TotalTokens), nil
	}

	return "", fmt.Errorf("No content found.")
}
