package main

import (
	"testing"
)

func TestDB(t *testing.T) {
	const (
		shortPath = "/shorted"
		srcURL    = "https://github.com/nekidb"
	)

	storage := createAndFillBoltStorage(t, "test.db")

	err := storage.Put(shortPath, srcURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get source URL by short path", func(t *testing.T) {
		got, err := storage.GetSrcURL(shortPath)
		if err != nil {
			t.Fatal(err)
		}
		assertString(t, got, srcURL)
	})

	t.Run("get source URL by short path", func(t *testing.T) {
		got, err := storage.GetShortPath(srcURL)
		if err != nil {
			t.Fatal(err)
		}
		assertString(t, got, shortPath)
	})
}

func assertString(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func createAndFillBoltStorage(t *testing.T, dbName string) *BoltStorage {
	storage, err := NewBoltStorage("test.db")
	if err != nil {
		t.Fatal(err)
	}

	return storage
}
