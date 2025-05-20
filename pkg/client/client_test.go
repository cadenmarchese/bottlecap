package client

import (
	"testing"
)

func TestClientLoadConfigFails(t *testing.T) {
	_, err := Client("subcommand", "argument")
	if err == nil {
		t.Errorf("expected an error when loading config fails")
	}
	if err.Error() != "failed to load config: open config.json: no such file or directory" {
		t.Errorf("expected error message 'failed to load config: <your specific error message>', got '%v'", err)
	}
}
