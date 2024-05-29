package cache

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var cacheDir string
var dirPerm os.FileMode = 0755
var filePerm os.FileMode = 0644

func SetupCacheDir() error {
	cacheDir = os.Getenv("CACHE_DIR")

	if _, err := os.Stat(cacheDir); err != nil {
		err := os.Mkdir(cacheDir, dirPerm)
		if err != nil {
			log.Println("Failed to create cache dir:", cacheDir)
			return err
		}
		log.Println("Created cache dir:", cacheDir)
	}

	return nil
}

func ReadProject(projectId string) ([]string, error) {
	projectDir := filepath.Join(cacheDir, projectId)

	files, err := os.ReadDir(projectDir)
	if err != nil {
		return make([]string, 0), nil
	}

	fileNames := make([]string, 0)
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	return fileNames, nil
}

func ReadImage(projectId, fileName string) ([]byte, error) {
	filePath := filepath.Join(cacheDir, projectId, fileName)

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read image %s", filePath))
	}

	return file, nil
}

func CacheImage(img []byte, name, projectId string) error {
	projectDir := filepath.Join(cacheDir, projectId)
	if _, err := os.Stat(projectDir); err != nil {
		if err := os.Mkdir(projectDir, dirPerm); err != nil {
			return errors.New(fmt.Sprintf("Failed to create project directory %s", projectDir))
		}
	}

	filePath := filepath.Join(projectDir, name)
	if os.WriteFile(filePath, img, filePerm) != nil {
		return errors.New(fmt.Sprintf("Failed to write image to cache %s", filePath))
	}

	return nil
}

func DeleteImage(projectId, fileName string) error {
	filePath := filepath.Join(cacheDir, projectId, fileName)
	if os.Remove(filePath) != nil {
		return errors.New(fmt.Sprintf("Failed to delete image from cache %s", filePath))
	}

	return nil
}

func ClearCache(projectId string) error {
	projectDir := filepath.Join(cacheDir, projectId)
	if _, err := os.Stat(projectDir); err != nil {
		return errors.New(fmt.Sprintf("Project directory %s not found", projectDir))
	}

	if err := os.RemoveAll(projectDir); err != nil {
		return errors.New(fmt.Sprintf("Failed to delete project cache %s", projectDir))
	}

	return nil
}
