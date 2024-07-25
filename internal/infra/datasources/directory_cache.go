package datasources

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jictyvoo/olympics_data_fetcher/internal/utils"
)

func loadExistentFolderCache(folder *os.File) (map[string][]byte, error) {
	result := make(map[string][]byte)

	// Read the directory contents
	files, err := os.ReadDir(folder.Name())
	if err != nil {
		return nil, err
	}

	// Iterate over the directory contents
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Read the file content
		filePath := filepath.Join(folder.Name(), file.Name())
		content, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return nil, readErr
		}

		// Add the file name and content to the map
		result[file.Name()] = content
	}

	return result, nil
}

type DirectoryCache struct {
	rootPath    string
	folderRef   *os.File
	loadedCache map[string][]byte
}

func NewDirectoryCache(rootPath string) (*DirectoryCache, error) {
	if err := utils.CreateDirIfNotExist(rootPath); err != nil {
		return nil, err
	}
	folderRef, err := os.Open(rootPath)
	instanceCache := DirectoryCache{rootPath: rootPath, folderRef: folderRef}
	if err != nil {
		return nil, err
	}

	instanceCache.loadedCache, err = loadExistentFolderCache(folderRef)
	return &instanceCache, err
}

func (d *DirectoryCache) Read(key string) ([]byte, error) {
	if content, ok := d.loadedCache[key]; ok {
		return content, nil
	}

	// Read from file
	filename := filepath.Join(d.rootPath, key)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Store the content in the cache
	d.loadedCache[key] = content

	return content, nil
}

func (d *DirectoryCache) Write(key string, data []byte) error {
	d.loadedCache[key] = data
	// Save on a file
	filename := filepath.Join(d.rootPath, key)
	// Ensure that the folder exists
	if err := utils.CreateDirIfNotExist(filepath.Dir(filename)); err != nil {
		return err
	}

	// Create the file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(data)
	return err
}
