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
