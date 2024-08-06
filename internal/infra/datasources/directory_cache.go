package datasources

import (
	"io"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/utils"
)

func loadExistentFolderCache(folder string) (map[string]cacheData, error) {
	result := make(map[string]cacheData)

	// Read the directory contents
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	// Iterate over the directory contents
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			subDirCache, subErr := loadExistentFolderCache(filepath.Join(folder, file.Name()))
			if subErr == nil {
				maps.Copy(result, subDirCache)
			}
			continue
		}

		// Read the file content
		filePath := filepath.Join(folder, file.Name())
		content, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return nil, readErr
		}

		// Add the file name and content to the map
		result[file.Name()] = cacheData{
			key:        file.Name(),
			identifier: filePath,
			content:    content,
		}
	}

	return result, nil
}

type (
	cacheData struct {
		key        string
		identifier string
		content    []byte
	}
	DirectoryCache struct {
		rootPath      string
		folderRef     *os.File
		loadedCache   map[string]cacheData
		cacheDuration time.Duration
	}
)

func NewDirectoryCache(rootPath string, cacheDuration time.Duration) (*DirectoryCache, error) {
	if err := utils.CreateDirIfNotExist(rootPath); err != nil {
		return nil, err
	}
	folderRef, err := os.Open(rootPath)
	instanceCache := DirectoryCache{
		rootPath:      rootPath,
		folderRef:     folderRef,
		cacheDuration: cacheDuration,
	}
	if err != nil {
		return nil, err
	}

	instanceCache.loadedCache, err = loadExistentFolderCache(folderRef.Name())
	if instanceCache.cacheDuration == 0 {
		instanceCache.cacheDuration = 4 * time.Minute
	}
	return &instanceCache, err
}

func (d *DirectoryCache) subFolderName() string {
	now := time.Now()
	subFolderName := now.Format("20060102")
	day := now.Day()
	hour := now.Hour()
	var divisionResult int
	switch {
	case d.cacheDuration < time.Hour:
		divisionResult = now.Minute() / int(d.cacheDuration.Minutes())
	case d.cacheDuration < 24*time.Hour:
		divisionResult = now.Hour() / int(d.cacheDuration.Hours())
		hour = 0
	}
	subFolderName += d.cacheDuration.String() +
		"#" + strconv.Itoa(day) +
		"@" + strconv.Itoa(hour) +
		"__" + strconv.Itoa(divisionResult)

	return subFolderName
}

func (d *DirectoryCache) Read(key string) ([]byte, error) {
	subFolderName := d.subFolderName()
	filename := filepath.Join(d.rootPath, subFolderName, key)
	if cData, ok := d.loadedCache[key]; ok {
		// Check if the identifier is the filename
		if cData.identifier == filename {
			return cData.content, nil
		}
	}

	// Read from file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Store the cData in the cache
	d.writeMemory(key, filename, content)

	return content, nil
}

func (d *DirectoryCache) writeMemory(key string, filename string, data []byte) {
	d.loadedCache[key] = cacheData{
		key:        key,
		identifier: filename,
		content:    data,
	}
}

func (d *DirectoryCache) Write(key string, data []byte) error {
	subFolderName := filepath.Join(d.rootPath, d.subFolderName())
	filename := filepath.Join(subFolderName, key)

	d.writeMemory(key, filename, data)
	// Ensure that the folder exists
	if err := utils.CreateDirIfNotExist(subFolderName); err != nil {
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
