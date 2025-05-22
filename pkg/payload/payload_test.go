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

func TestCreateImageGenerationPayload(t *testing.T) {
	tests := []struct {
		name       string
		model      string
		userInput  string
		wantErr    bool
		wantKeyLen int
	}{
		{
			name:       "success",
			model:      "model",
			userInput:  "generate an image of a cat",
			wantErr:    false,
			wantKeyLen: 4,
		},
		{
			name:      "empty model",
			model:     "",
			userInput: "generate an image of a cat",
			wantErr:   false,
			wantKeyLen: 4,
		},
		{
			name:       "empty user input",
			model:      "model",
			userInput:  "",
			wantErr:    false,
			wantKeyLen: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := CreateImageGenerationPayload(tt.model, tt.userInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateImageGenerationPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(payload) == 0 {
				t.Errorf("expected a payload but got an empty string")
			}

			var payloadMap map[string]interface{}
			err = json.Unmarshal([]byte(payload), &payloadMap)
			if err != nil {
				t.Errorf("failed to unmarshal payload: %v", err)
			}

			if len(payloadMap) != tt.wantKeyLen {
				t.Errorf("expected payload map to have %d keys but got %d", tt.wantKeyLen, len(payloadMap))
			}
		})
	}
}