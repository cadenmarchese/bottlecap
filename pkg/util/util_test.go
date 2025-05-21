package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadImageAsBytes(t *testing.T) {
	// Create a mock image server
	mockImage := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG header bytes
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockImage)
	}))
	defer server.Close()

	// Call the function with the mock server URL
	data, err := DownloadAndConvertImageToBytes(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(data) != len(mockImage) {
		t.Errorf("Expected %d bytes, got %d", len(mockImage), len(data))
	}

	if string(data) != string(mockImage) {
		t.Errorf("Expected data %v, got %v", mockImage, data)
	}
}
