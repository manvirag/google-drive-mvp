package models

import (
	"time"
)

type FileMetadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Chunks      []Chunk   `json:"chunks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"`
}

type Chunk struct {
	ID     string `json:"id"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
	Index  int    `json:"index"`
	Offset int64  `json:"offset"`
}

type UploadRequest struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type UploadResponse struct {
	FileID  string `json:"file_id"`
	Message string `json:"message"`
}

type UpdateRequest struct {
	FileID string `json:"file_id"`
	Name   string `json:"name,omitempty"`
}

type FileListResponse struct {
	Files []FileMetadata `json:"files"`
	Total int            `json:"total"`
}
