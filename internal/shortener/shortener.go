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
	Save(shortPath, srcURL string)
	GetShortPath(srcURL string) string
	GetSourceURL(shortPath string) string
	DeleteSourceURL(srcURL string)
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
	shortPath := s.storage.GetShortPath(srcURL)
	if shortPath != "" {
		return shortPath, nil
	}

	// If not exists in DB, then create new
	length := 6
	shortPath = s.generateRandomPath(length)

	s.storage.Save(shortPath, srcURL)

	return shortPath, nil
}

func (s ShortenerService) GetSourceURL(shortPath string) (string, error) {
	// Check if shortPath is correct

	// Get from DB if exists. If not, then return empty string
	srcURL := s.storage.GetSourceURL(shortPath)
	return srcURL, nil
}

func (s ShortenerService) DeleteSourceURL(srcURL string) error {
	s.storage.DeleteSourceURL(srcURL)
	return nil
}

//	func (s ShortenerService) ValidateURL(str string) (bool, error) {
//		u, err := url.Parse(str)
//		if err != nil {
//			return false, err
//		}
//
//		if u.Scheme == "" || !isGoodHost(u.Host) {
//			return false, nil
//		}
//		return true, nil
//	}
//
//	func isGoodHost(host string) bool {
//		if host == "" {
//			return false
//		}
//		if !strings.Contains(host, ".") {
//			return false
//		}
//		return true
//	}
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s ShortenerService) generateRandomPath(n int) string {
	var sb strings.Builder
	sb.Grow(n)

	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}

	return sb.String()
}
