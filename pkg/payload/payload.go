package payload

import (
	"encoding/json"
	"fmt"
)

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

func CreateImageRequestPayload(model, instruction string, imageUrl string) (string, error) {
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{
				"content": instruction,
				"role":    "system",
			},
			{
				"content": imageUrl,
				"role":    "user",
				"type":    "image",
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(data), nil
}
