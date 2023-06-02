package config

import (
	"encoding/json"
	"testing"
	"testing/fstest"
)

func TestConfig(t *testing.T) {
	filename := "somefile"

	config := Config{
		Host: "somehost",
		Port: "someport",
		DB:   "somedb",
	}

	fs := createFsWithConfig(t, filename, config)

	got, _ := Get(fs, filename)
	want := config
	if *got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func createFsWithConfig(t *testing.T, filename string, config Config) fstest.MapFS {
	t.Helper()

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatal("marshal test config: ", err)
	}

	fs := fstest.MapFS{
		filename: &fstest.MapFile{
			Data: data,
		},
	}

	return fs
}
