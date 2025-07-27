package models_test

import (
	"testing"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	// Arrange
	name := "test-file.pdf"
	originalName := "original-test-file.pdf"
	description := "Test file description"
	fileType := "pdf"
	s3Key := "files/1/test-file.pdf"
	s3URL := "https://bucket.s3.amazonaws.com/files/1/test-file.pdf"
	fileSize := int64(1024 * 1024) // 1MB
	creatorID := uint(1)

	// Act
	file := models.NewFile(name, originalName, description, fileType, s3Key, s3URL, fileSize, creatorID)

	// Assert
	assert.NotNil(t, file)
	assert.Equal(t, name, file.Name)
	assert.Equal(t, originalName, file.OriginalName)
	assert.Equal(t, description, file.Description)
	assert.Equal(t, fileType, file.FileType)
	assert.Equal(t, s3Key, file.S3Key)
	assert.Equal(t, s3URL, file.S3URL)
	assert.Equal(t, fileSize, file.FileSize)
	assert.Equal(t, creatorID, file.CreatorID)
	assert.True(t, file.Status)
}

func TestGetFileSizeFormatted(t *testing.T) {
	tests := []struct {
		name     string
		fileSize int64
		expected string
	}{
		{
			name:     "Bytes",
			fileSize: 512,
			expected: "512 B",
		},
		{
			name:     "Kilobytes",
			fileSize: 1024 * 2,
			expected: "2.0 KB",
		},
		{
			name:     "Megabytes",
			fileSize: 1024 * 1024 * 3,
			expected: "3.0 MB",
		},
		{
			name:     "Gigabytes",
			fileSize: 1024 * 1024 * 1024 * 4,
			expected: "4.0 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &models.File{FileSize: tt.fileSize}
			result := file.GetFileSizeFormatted()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPDF(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		expected bool
	}{
		{
			name:     "PDF file",
			fileType: "pdf",
			expected: true,
		},
		{
			name:     "Document file",
			fileType: "document",
			expected: false,
		},
		{
			name:     "Image file",
			fileType: "image",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &models.File{FileType: tt.fileType}
			result := file.IsPDF()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		expected bool
	}{
		{
			name:     "Image file",
			fileType: "image",
			expected: true,
		},
		{
			name:     "PDF file",
			fileType: "pdf",
			expected: false,
		},
		{
			name:     "Document file",
			fileType: "document",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &models.File{FileType: tt.fileType}
			result := file.IsImage()
			assert.Equal(t, tt.expected, result)
		})
	}
}
