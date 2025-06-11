package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"drive-mvp/internal/models"

	"github.com/google/uuid"
)

type FileService struct {
	filesDir     string
	chunkService *ChunkService
}

func NewFileService(filesDir string, chunkService *ChunkService) *FileService {
	return &FileService{
		filesDir:     filesDir,
		chunkService: chunkService,
	}
}

func (fs *FileService) UploadFile(name, contentType string, reader io.Reader) (*models.FileMetadata, error) {
	fileID := uuid.New().String()

	chunks, err := fs.chunkService.CreateChunks(reader)
	if err != nil {
		return nil, fmt.Errorf("error creating chunks: %w", err)
	}

	var totalSize int64
	for _, chunk := range chunks {
		totalSize += chunk.Size
	}

	fileMetadata := &models.FileMetadata{
		ID:          fileID,
		Name:        name,
		Size:        totalSize,
		ContentType: contentType,
		Chunks:      chunks,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     1,
	}

	if err := fs.saveFileMetadata(fileMetadata); err != nil {
		return nil, fmt.Errorf("error saving file metadata: %w", err)
	}

	return fileMetadata, nil
}

func (fs *FileService) UpdateFile(fileID string, reader io.ReadSeeker) (*models.FileMetadata, error) {
	existingFile, err := fs.GetFileMetadata(fileID)
	if err != nil {
		return nil, fmt.Errorf("error loading existing file: %w", err)
	}

	newChunks, err := fs.chunkService.CreateChunks(reader)
	if err != nil {
		return nil, fmt.Errorf("error creating new chunks: %w", err)
	}

	changedChunks, err := fs.chunkService.CompareChunks(existingFile.Chunks, newChunks)
	if err != nil {
		return nil, fmt.Errorf("error comparing chunks: %w", err)
	}

	if len(changedChunks) > 0 {
		if err := fs.chunkService.UpdateChunks(changedChunks, reader); err != nil {
			return nil, fmt.Errorf("error updating chunks: %w", err)
		}
	}

	var totalSize int64
	for _, chunk := range newChunks {
		totalSize += chunk.Size
	}

	existingFile.Chunks = newChunks
	existingFile.Size = totalSize
	existingFile.UpdatedAt = time.Now()
	existingFile.Version++

	if err := fs.saveFileMetadata(existingFile); err != nil {
		return nil, fmt.Errorf("error saving updated metadata: %w", err)
	}

	return existingFile, nil
}

func (fs *FileService) DownloadFile(fileID string, writer io.Writer) (*models.FileMetadata, error) {
	fileMetadata, err := fs.GetFileMetadata(fileID)
	if err != nil {
		return nil, fmt.Errorf("error loading file metadata: %w", err)
	}

	if err := fs.chunkService.ReconstructFile(fileMetadata.Chunks, writer); err != nil {
		return nil, fmt.Errorf("error reconstructing file: %w", err)
	}

	return fileMetadata, nil
}

func (fs *FileService) GetFileMetadata(fileID string) (*models.FileMetadata, error) {
	metadataPath := filepath.Join(fs.filesDir, fileID+".json")

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata file: %w", err)
	}

	var metadata models.FileMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("error unmarshaling metadata: %w", err)
	}

	return &metadata, nil
}

func (fs *FileService) ListFiles() ([]models.FileMetadata, error) {
	files, err := os.ReadDir(fs.filesDir)
	if err != nil {
		return nil, fmt.Errorf("error reading files directory: %w", err)
	}

	var fileList []models.FileMetadata
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fileID := file.Name()[:len(file.Name())-5]
			metadata, err := fs.GetFileMetadata(fileID)
			if err != nil {
				continue
			}
			fileList = append(fileList, *metadata)
		}
	}

	return fileList, nil
}

func (fs *FileService) DeleteFile(fileID string) error {
	metadata, err := fs.GetFileMetadata(fileID)
	if err != nil {
		return fmt.Errorf("error loading file metadata: %w", err)
	}

	metadataPath := filepath.Join(fs.filesDir, fileID+".json")
	if err := os.Remove(metadataPath); err != nil {
		return fmt.Errorf("error deleting metadata file: %w", err)
	}

	_ = metadata

	return nil
}

func (fs *FileService) saveFileMetadata(metadata *models.FileMetadata) error {
	metadataPath := filepath.Join(fs.filesDir, metadata.ID+".json")

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("error writing metadata file: %w", err)
	}

	return nil
}
