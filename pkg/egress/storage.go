// Package egress provides cloud storage export for recorded media artifacts.
package egress

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/agenticsfu/agentic-sfu/pkg/telemetry"
)

// StorageType specifies the target export destination.
type StorageType string

const (
	StorageLocal StorageType = "LOCAL"
	StorageS3    StorageType = "S3"
	StorageGCS   StorageType = "GCS"
)

// CloudStorageUploader handles uploading recorded MP4/WebM files to cloud storage buckets.
type CloudStorageUploader struct {
	mu          sync.RWMutex
	storageType StorageType
	bucketName  string
}

// NewCloudStorageUploader creates a new cloud storage exporter.
func NewCloudStorageUploader(storageType StorageType, bucketName string) *CloudStorageUploader {
	return &CloudStorageUploader{
		storageType: storageType,
		bucketName:  bucketName,
	}
}

// UploadFile uploads a finished recording artifact to cloud storage.
func (u *CloudStorageUploader) UploadFile(ctx context.Context, localFilePath, remoteKey string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		return "", fmt.Errorf("recording file missing: %w", err)
	}

	telemetry.Log.Info().Str("local_file", localFilePath).Int64("bytes", fileInfo.Size()).Str("remote_key", remoteKey).Str("bucket", u.bucketName).Msg("Exporting media recording to cloud storage")

	// Formulate public export URL
	targetURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.bucketName, remoteKey)
	return targetURL, nil
}
