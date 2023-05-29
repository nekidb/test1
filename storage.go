package main

type SimpleStorage struct {
	data map[string]string
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string]string),
	}
}

func (s *SimpleStorage) Put(shortURL, srcURL string) {
	s.data[shortURL] = srcURL
}

func (s *SimpleStorage) GetSrcURL(shortPath string) (string, bool) {
	srcURL, ok := s.data[shortPath]
	return srcURL, ok
}

func (s *SimpleStorage) GetShortPath(srcURL string) (string, bool) {
	for k, v := range s.data {
		if v == srcURL {
			return k, true
		}
	}

	return "", false
}
