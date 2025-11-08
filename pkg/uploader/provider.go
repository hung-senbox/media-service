package uploader

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type UploadMode int

const (
	Unknown UploadMode = iota
	UploadPrivate
	UploadPublic
)

// String returns the string representation of the enum
func (s UploadMode) String() string {
	return [...]string{"UploadPrivate", "UploadPublic", "Unknown"}[s]
}

// UploadModeFromString converts a string to UploadMode
func UploadModeFromString(s string) (UploadMode, error) {
	switch strings.ToLower(s) {
	case "private":
		return UploadPrivate, nil
	case "public":
		return UploadPublic, nil
	default:
		return Unknown, errors.New("invalid upload mode string")
	}
}

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, dest string, mode UploadMode) (*string, error)
	SaveFileUploadedReader(ctx context.Context, r io.Reader, dest string, contentType string, mode UploadMode) (*string, error)
	GetFileUploaded(ctx context.Context, key string, duration *time.Duration) (*string, error)
	DeleteFileUploaded(ctx context.Context, key string) error
}
