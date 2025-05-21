package client

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cadenmarchese/bottlecap/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		wantErr  bool
		expected string
	}{
		{
			name:     "valid response",
			body:     []byte(`{"choices":[{"message":{"content":"hello"}}],"usage":{"completion_tokens":1,"prompt_tokens":1,"total_tokens":2}}`),
			wantErr:  false,
			expected: "hello",
		},
		{
			name:    "invalid JSON",
			body:    []byte(`invalid`),
			wantErr: true,
		},
		{
			name:    "empty choices",
			body:    []byte(`{"choices":[]}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := parseResponse(tt.body)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, resp, tt.expected)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	config := types.Config{
		Model:             "test-model",
		ChatInstructions:  "chat",
		ImageInstructions: "image",
		URL:               "http://mock.url",
		BearerToken:       "token",
	}
	enc := json.NewEncoder(tmpFile)
	err = enc.Encode(config)
	assert.NoError(t, err)
	tmpFile.Close()

	t.Run("valid config", func(t *testing.T) {
		cfg, err := loadConfig(tmpFile.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test-model", cfg.Model)
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := loadConfig("nonexistent.json")
		assert.Error(t, err)
	})
}
