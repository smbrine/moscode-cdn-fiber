package iosync

import (
	"log"
	"os"
	"path/filepath"
)

func DeleteFile(fullPath string) error {
	if err := os.RemoveAll(fullPath); err != nil {
		return err
	}
	return nil
}

func EmptyDir(dirPath string) error {
	// Attempt to open the directory.
	d, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist, which is fine for our purpose.
			log.Printf("Directory %s does not exist, considered as empty.", dirPath)
			return nil
		}
		// Return errors other than "not exists".
		return err
	}
	defer func() {
		if cerr := d.Close(); cerr != nil {
			log.Printf("Warning: Failed to close directory handle for %s: %v", dirPath, cerr)
		}
	}()

	// Read directory contents.
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	// If directory is already empty, nothing to do.
	if len(names) == 0 {
		log.Printf("Directory %s is already empty.", dirPath)
		return nil
	}

	// Loop through and remove each item in the directory.
	for _, name := range names {
		fullPath := filepath.Join(dirPath, name)
		if err := DeleteFile(fullPath); err != nil {
			return err
		}
	}

	log.Printf("Successfully emptied directory: %s", dirPath)
	return nil
}
