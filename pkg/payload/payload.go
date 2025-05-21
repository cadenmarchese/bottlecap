package payload

import (
	"encoding/base64"
	"encoding/json"

	"github.com/cadenmarchese/bottlecap/pkg/util"
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

// experimental - will download an image from URL and provide it to the LLM
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
				"role":    "system",
				"content": instruction,
			},
			{
				"role":    "user",
				"content": instruction,
			},
		},
		"input": map[string]interface{}{
			"text": instruction,
		},
		"image_data": base64.StdEncoding.EncodeToString(image),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
