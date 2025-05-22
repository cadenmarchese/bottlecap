package payload

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cadenmarchese/bottlecap/pkg/util"
)

// CreateChatRequestPayload constructs a basic chat request, using the config file and user's input
// https://platform.openai.com/docs/api-reference/chat/create
func CreateChatRequestPayload(model, instruction, userInput string) (string, error) {
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"content": instruction,
				"role":    "system",
			},
			{
				"content": userInput,
				"role":    "user",
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CreateImageRequestPayload downloads an image from URL, and provides it to the LLM via a base64 string
// it does this because some LLMs don't like downloading images from the internet, especially local LLMs
// https://platform.openai.com/docs/api-reference/chat/create
func CreateImageRequestPayload(model, instruction string, imageUrl string) (string, error) {
	var image []byte
	image, err := util.DownloadAndConvertImageToBytes(imageUrl)
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": instruction,
					},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url":    fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(image)),
							"detail": "low",
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CreateImageGenerationPayload asks the LLM to generate an image. If invoked, the url in config.json
// must point to the v1/images/generations API. For example:
// https://platform.openai.com/docs/api-reference/images
func CreateImageGenerationPayload(model, userInput string) (string, error) {
	payload := map[string]interface{}{
		"model":  model,
		"prompt": userInput,
		"n":      1,
		"size":   "1024x1024",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
