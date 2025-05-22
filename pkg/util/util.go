package util

import (
	"fmt"
	"io"
	"net/http"
)

// DownloadAndConvertImageToBytes pulls an image from URL and then converts it to []byte
// for use in the request payload.
func DownloadAndConvertImageToBytes(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageData, nil
}
