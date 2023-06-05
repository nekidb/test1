package shortener

import (
	"math/rand"
	"strings"
)

type Shortener interface {
	GetShortPath(srcURL string) (string, error)
	GetSourceURL(shortPath string) (string, error)
	DeleteSourceURL(srcURL string) error
	// GetCustomShortURL(source, custom string) error
}

type URLStorage interface {
	Save(shortPath, srcURL string) error
	GetShortPath(srcURL string) (string, error)
	GetSourceURL(shortPath string) (string, error)
	DeleteSourceURL(srcURL string) error
}

type ShortenerService struct {
	storage URLStorage
}

func NewShortenerService(storage URLStorage) (*ShortenerService, error) {
	return &ShortenerService{
		storage: storage,
	}, nil
}

func (s ShortenerService) GetShortPath(srcURL string) (string, error) {
	// Check if exists in DB
	shortPath, err := s.storage.GetShortPath(srcURL)
	if err != nil {
		return "", err
	}
	if shortPath != "" {
		return shortPath, nil
	}

	// If not exists in DB, then create new
	length := 6
	shortPath = s.generateRandomPath(length)

	if err := s.storage.Save(shortPath, srcURL); err != nil {
		return "", err
	}

	return shortPath, nil
}

func (s ShortenerService) GetSourceURL(shortPath string) (string, error) {
	// Check if shortPath is correct

	// Get from DB if exists. If not, then return empty string
	return s.storage.GetSourceURL(shortPath)
}

func (s ShortenerService) DeleteSourceURL(srcURL string) error {
	return s.storage.DeleteSourceURL(srcURL)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s ShortenerService) generateRandomPath(n int) string {
	var sb strings.Builder
	sb.Grow(n)

	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}

	return sb.String()
}
