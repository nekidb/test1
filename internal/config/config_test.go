package config

import (
	"testing"
	"testing/fstest"
)

func TestConfig(t *testing.T) {
	filename := "somefile"

	data := []byte(`{"host":"somehost","port":"someport","db":"somedb"}`)
	want := Config{
		Host: "somehost",
		Port: "someport",
		DB:   "somedb",
	}

	fs := createFsWithData(t, filename, data)

	got, err := Get(fs, filename)
	if err != nil {
		t.Fatal(err)
	}

	if *got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func createFsWithData(t *testing.T, filename string, data []byte) fstest.MapFS {
	t.Helper()

	fs := fstest.MapFS{
		filename: &fstest.MapFile{
			Data: data,
		},
	}

	return fs
}
