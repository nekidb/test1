package shortener

import "testing"

func TestShortener(t *testing.T) {
	const (
		shortPath = "shorted"
		srcURL    = "htts://github.com"
	)

	t.Run("writes to DB if source URL not exists", func(t *testing.T) {
		storage := NewMockStorage()
		shortener := createShortener(t, storage)

		if _, err := shortener.GetShortPath(srcURL); err != nil {
			t.Fatal(err)
		}

		got := len(storage.data)
		want := 1
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("gets shorted path from DB", func(t *testing.T) {
		storage := NewMockStorage()
		storage.Save(shortPath, srcURL)
		shortener := createShortener(t, storage)

		got, err := shortener.GetShortPath(srcURL)
		if err != nil {
			t.Fatal(err)
		}
		assertString(t, got, shortPath)
	})
	t.Run("gets source URL from DB", func(t *testing.T) {
		storage := NewMockStorage()
		storage.Save(shortPath, srcURL)
		shortener := createShortener(t, storage)

		got, err := shortener.GetSourceURL(shortPath)
		if err != nil {
			t.Fatal(err)
		}
		assertString(t, got, srcURL)
	})
	t.Run("deletes source URL", func(t *testing.T) {
		storage := NewMockStorage()
		storage.Save(shortPath, srcURL)
		shortener := createShortener(t, storage)

		shortener.DeleteSourceURL(srcURL)

		got := len(storage.data)
		want := 0
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

func createShortener(t *testing.T, storage URLStorage) *ShortenerService {
	t.Helper()
	shortener, err := NewShortenerService(storage)
	if err != nil {
		t.Fatal(err)
	}

	return shortener
}

func assertString(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

type MockStorage struct {
	data map[string]string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]string),
	}
}

func (s *MockStorage) Save(str1, str2 string) {
	s.data[str1] = str2
}

func (s *MockStorage) GetShortPath(str string) string {
	for k, v := range s.data {
		if v == str {
			return k
		}
	}

	return ""
}

func (s *MockStorage) GetSourceURL(str string) string {
	return s.data[str]
}

func (s *MockStorage) DeleteSourceURL(str string) {
	for k, v := range s.data {
		if v == str {
			delete(s.data, k)
		}
	}
}
