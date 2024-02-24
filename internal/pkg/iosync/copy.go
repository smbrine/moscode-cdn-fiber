package iosync

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := sourceFile.Close(); cerr != nil {
			log.Printf("Failed to close source file %s: %v", src, cerr)
		}
	}()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := destinationFile.Close(); cerr != nil {
			log.Printf("Failed to close destination file %s: %v", dst, cerr)
		}
	}()

	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	return nil
}

func CopyFilesInDir(srcDir, dstDir string) error {
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return err
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories.
		}

		// Calculate relative path to preserve directory structure.
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %s: %v", path, err)
		}
		dstFilePath := filepath.Join(dstDir, relPath)

		// Ensure the destination directory exists.
		if err := os.MkdirAll(filepath.Dir(dstFilePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories for %s: %v", dstFilePath, err)
		}

		if err := CopyFile(path, dstFilePath); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %v", path, dstFilePath, err)
		}

		return nil
	})
}
