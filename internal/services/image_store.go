package services

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(taskID string, imageType string, imageData bytes.Buffer) (string, error)
}

type diskImageStore struct {
	mutex       sync.Mutex
	imageFolder string
	images      map[string]*ImageInfo
}

type ImageInfo struct {
	TaskID string
	Type   string
	Path   string
}

func NewDiskImageStore(imageFolder string) ImageStore {
	return &diskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

func (store *diskImageStore) Save(taskID string, imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot crate image file: %w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageID.String()] = &ImageInfo{
		TaskID: taskID,
		Type:   imageType,
		Path:   imagePath,
	}

	return imageID.String(), nil
}
