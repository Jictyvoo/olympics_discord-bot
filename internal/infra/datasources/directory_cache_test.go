package datasources

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewDirectoryCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_cache")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	cacheDuration := 10 * time.Minute
	cache, err := NewDirectoryCache(tempDir, cacheDuration)
	if err != nil {
		t.Fatalf("NewDirectoryCache() error = %v", err)
	}

	if cache.rootPath != tempDir {
		t.Errorf("Expected rootPath to be %s, got %s", tempDir, cache.rootPath)
	}
	if cache.cacheDuration != cacheDuration {
		t.Errorf("Expected cacheDuration to be %v, got %v", cacheDuration, cache.cacheDuration)
	}
	if cache.loadedCache == nil {
		t.Error("Expected loadedCache to be initialized, got nil")
	}
}

func TestDirectoryCache_ReadWrite(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_cache")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	cache, err := NewDirectoryCache(tempDir, 10*time.Minute)
	if err != nil {
		t.Fatalf("NewDirectoryCache() error = %v", err)
	}

	testData := []byte("test data")
	testKey := "testfile.txt"

	err = cache.Write(testKey, testData)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	readData, err := cache.Read(testKey)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if string(readData) != string(testData) {
		t.Errorf("Expected read data to be %s, got %s", string(testData), string(readData))
	}
}

func TestDirectoryCache_SubFolderName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_cache")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	cacheDuration := 10 * time.Minute
	cache, err := NewDirectoryCache(tempDir, cacheDuration)
	if err != nil {
		t.Fatalf("NewDirectoryCache() error = %v", err)
	}

	subFolderName := cache.subFolderName()
	expectedSubFolderNamePrefix := time.Now().Format("20060102")

	if !strings.HasPrefix(subFolderName, expectedSubFolderNamePrefix) {
		t.Errorf(
			"Expected subFolderName to start with %s, got %s",
			expectedSubFolderNamePrefix,
			subFolderName,
		)
	}
}

func TestLoadExistentFolderCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_cache")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	// Create a test file in the temp directory
	testFileName := "testfile.txt"
	testFilePath := filepath.Join(tempDir, testFileName)
	testData := []byte("test data")

	err = os.WriteFile(testFilePath, testData, 0600)
	if err != nil {
		t.Fatal(err)
	}

	folder, err := os.Open(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func(folder *os.File) {
		if err := folder.Close(); err != nil {
			t.Fatal(err)
		}
	}(folder)

	cache, err := loadExistentFolderCache(folder)
	if err != nil {
		t.Fatalf("loadExistentFolderCache() error = %v", err)
	}

	if len(cache) != 1 {
		t.Errorf("Expected cache to contain 1 entry, got %d", len(cache))
	}
	if string(cache[testFileName]) != string(testData) {
		t.Errorf(
			"Expected cache data to be %s, got %s",
			string(testData),
			string(cache[testFileName]),
		)
	}
}
