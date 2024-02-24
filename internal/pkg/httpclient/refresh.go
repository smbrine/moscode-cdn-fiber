package httpclient

import (
	"fmt"
	"log"
	"moscode-cdn-fiber/internal/pkg/iosync"
	"path/filepath"
)

func RefreshFile(url string, file string, tempDir string, staticDir string) error {

	files := []string{file, file + ".br", file + ".gz"}

	for _, file := range files {
		if err := iosync.DeleteFile(filepath.Join(tempDir, file)); err != nil {
			log.Println(fmt.Sprintf("Error deleting file from temp directory: %s", err))
			return nil // Since you asked to return nil
		}
	}

	if err := iosync.FetchAndSaveFile(url, file, tempDir); err != nil {
		log.Println(fmt.Sprintf("Error fetching and saving file: %s", err))
		return nil // Since you asked to return nil
	}

	for _, file := range files {
		if err := iosync.DeleteFile(filepath.Join(staticDir, file)); err != nil {
			log.Println(fmt.Sprintf("Error deleting file from static directory: %s", err))
			return nil // Since you asked to return nil
		}
	}

	for _, file := range files {
		if err := iosync.CopyFile(filepath.Join(tempDir, file), filepath.Join(staticDir, file)); err != nil {
			log.Println(fmt.Sprintf("Error copying file to static directory: %s", err))
			return nil // Since you asked to return nil
		}
	}

	return nil
}
