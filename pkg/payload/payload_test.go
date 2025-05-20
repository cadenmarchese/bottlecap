package payload

import (
	"encoding/json"
	"testing"
)

func TestCreateChatRequestPayload(t *testing.T) {
	model := "testModel"
	instruction := "This is an instruction."
	userInput := "This is a user input."

	payload, err := CreateChatRequestPayload(model, instruction, userInput)
	if err != nil {
		t.Errorf("CreateChatRequestPayload() error = %v", err)
	}

	var payloadMap map[string]interface{}
	err = json.Unmarshal([]byte(payload), &payloadMap)
	if err != nil {
		t.Errorf("Failed to unmarshal payload: %v", err)
	}

	// Basic validation
	if payloadMap["model"] != model {
		t.Errorf("Expected model name %s in payload, got %v", model, payloadMap["model"])
	}

	messagesList, ok := payloadMap["messages"].([]interface{})
	if !ok || len(messagesList) != 2 {
		t.Errorf("Expected messages list with 2 entries, got %v", messagesList)
	}

	for i, msg := range messagesList {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			t.Errorf("Message %d is not a map", i)
			continue
		}

		switch i {
		case 0:
			if msgMap["content"].(string) != instruction || msgMap["role"].(string) != "system" {
				t.Errorf("Instruction message is incorrect: %+v", msgMap)
			}
		case 1:
			if msgMap["content"].(string) != userInput || msgMap["role"].(string) != "user" {
				t.Errorf("User input message is incorrect: %+v", msgMap)
			}
		}
	}
}

func TestCreateImageRequestPayload(t *testing.T) {
	model := "someModel"
	instruction := "create an image of a cat"
	imageUrl := "https://example.com/cat.jpg"

	expectedPayload := map[string]interface{}{
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

	expectedJSON, err := json.Marshal(expectedPayload)
	if err != nil {
		t.Fatal(err)
	}

	payload, err := CreateImageRequestPayload(model, instruction, imageUrl)
	if err != nil {
		t.Errorf("CreateImageRequestPayload() error = %v", err)
	}

	if string(expectedJSON) != payload {
		t.Errorf("CreateImageRequestPayload() payload = %v, want %v", payload, string(expectedJSON))
	}
}
