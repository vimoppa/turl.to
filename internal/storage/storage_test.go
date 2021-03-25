package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vimoppa/turl.to/internal/config"
)

func teardown(cfg *config.StorageConfiguration) {
	_, err := os.Stat(cfg.File)
	if err == nil {
		os.Remove(cfg.File)
	}
}

func TestStorage(t *testing.T) {
	file := filepath.Join(os.TempDir(), "test.txt")

	cfg := &config.StorageConfiguration{
		File: file,
	}
	defer teardown(cfg)

	store, err := New(cfg)
	if err != nil {
		t.Errorf("failed to create new store: %v", err)
	}

	longURL := "https://www.google.com"
	shortURL := "a3rnvsd"

	if err := store.WriteOnce(shortURL, longURL); err != nil {
		t.Errorf("failed to write to store: %v", err)
	}

	if exists := store.LookUp(longURL); !exists {
		t.Error("expected to find longURL in store")
	}

	result, err := store.ReadOne(shortURL)
	if err != nil {
		t.Error("expected to find longURL in store")
	}
	if result != longURL {
		t.Errorf("expected result (longURL) does not match")
	}

	_, err = store.ReadAll()
	if err != nil {
		t.Errorf("failed to read all records in store")
	}
}
