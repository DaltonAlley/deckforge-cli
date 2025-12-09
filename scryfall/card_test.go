package scryfall

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindCardByID(t *testing.T) {
	testID := "a65e485b-03a2-4634-9218-f5bb7c104d41"
	foundCard, err := FindCardByID(testID)
	require.NoError(t, err)
	fmt.Println(foundCard.Name)
}

func TestDownloadCardImage(t *testing.T) {
	t.Run("download normal quality image", func(t *testing.T) {
		// Use a known card ID
		testID := "a65e485b-03a2-4634-9218-f5bb7c104d41"

		// Create temp directory for test
		tempDir, err := os.MkdirTemp("", "card_images_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		imagePath, err := DownloadCardImage(testID, tempDir, "normal")
		require.NoError(t, err)
		require.NotEmpty(t, imagePath)

		// Verify file exists and has content
		info, err := os.Stat(imagePath)
		require.NoError(t, err)
		require.Greater(t, info.Size(), int64(0))
	})

	t.Run("download with invalid card ID", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "card_images_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		_, err = DownloadCardImage("invalid-id", tempDir, "normal")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to find card")
	})

	t.Run("download with caching", func(t *testing.T) {
		testID := "a65e485b-03a2-4634-9218-f5bb7c104d41"

		tempDir, err := os.MkdirTemp("", "card_images_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// First download
		imagePath1, err := DownloadCardImage(testID, tempDir, "normal")
		require.NoError(t, err)

		// Get file info for first download
		info1, err := os.Stat(imagePath1)
		require.NoError(t, err)

		// Second download (should use cache)
		imagePath2, err := DownloadCardImage(testID, tempDir, "normal")
		require.NoError(t, err)

		// Should be the same file
		require.Equal(t, imagePath1, imagePath2)

		// File should not have been re-downloaded (same size, same mod time)
		info2, err := os.Stat(imagePath2)
		require.NoError(t, err)
		require.Equal(t, info1.Size(), info2.Size())
		require.Equal(t, info1.ModTime(), info2.ModTime())
	})
}
