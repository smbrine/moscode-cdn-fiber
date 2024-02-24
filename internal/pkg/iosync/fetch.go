package iosync

import (
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func createBrotliVersion(filePath string) error {
	// Open the original file for reading
	originalFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening original file %s: %v", filePath, err)
	}
	defer func() {
		if err := originalFile.Close(); err != nil {
			log.Printf("error closing original file %s: %v", filePath, err)
		}
	}()

	// Create the .br compressed file for writing
	compressedFilePath := filePath + ".br"
	compressedFile, err := os.Create(compressedFilePath)
	if err != nil {
		return fmt.Errorf("error creating compressed file %s: %v", compressedFilePath, err)
	}
	defer func() {
		if err := compressedFile.Close(); err != nil {
			log.Printf("error closing compressed file %s: %v", filePath, err)
		}
	}()

	// Create a Brotli writer with max compression level
	brWriter := brotli.NewWriterLevel(compressedFile, brotli.BestCompression)
	defer func() {
		if err := brWriter.Close(); err != nil {
			log.Printf("error closing brotli writer %s: %v", filePath, err)
		}
	}()

	// Copy the original file content to the Brotli writer, compressing it in the process
	if _, err = io.Copy(brWriter, originalFile); err != nil {
		return fmt.Errorf("error compressing file %s: %v", filePath, err)
	}

	return nil
}

func createGzipVersion(filePath string) error {
	originalFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening original file %s: %v", filePath, err)
	}
	defer func() {
		if err := originalFile.Close(); err != nil {
			log.Printf("error closing original file %s: %v", filePath, err)
		}
	}()

	compressedFilePath := filePath + ".gz"
	compressedFile, err := os.Create(compressedFilePath)
	if err != nil {
		return fmt.Errorf("error creating compressed file %s: %v", compressedFilePath, err)
	}
	defer func() {
		if err := compressedFile.Close(); err != nil {
			log.Printf("error closing compressed file %s: %v", compressedFilePath, err)
		}
	}()

	gzWriter, _ := gzip.NewWriterLevel(compressedFile, gzip.BestCompression)
	defer func() {
		if err := gzWriter.Close(); err != nil {
			log.Printf("error closing gzip writer for file %s: %v", compressedFilePath, err)
		}
	}()

	if _, err = io.Copy(gzWriter, originalFile); err != nil {
		return fmt.Errorf("error compressing file %s: %v", filePath, err)
	}

	return nil
}

func FetchAndSaveFile(baseURL, filePath, targetDir string) error {
	fileURL := fmt.Sprintf("%s/%s", baseURL, filePath) // Construct the full URL to the file.
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("error fetching file %s: %v", filePath, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}

	// Calculate the full path to where the file should be saved, including any subdirectories.
	fullPath := filepath.Join(targetDir, filePath)

	// Ensure that the directory structure for fullPath exists.
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories for file %s: %v", fullPath, err)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", fullPath, err)
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	if _, err = io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("error writing file %s: %v", fullPath, err)
	}

	if err := createBrotliVersion(fullPath); err != nil {
		log.Printf("Error creating Brotli version of file %s: %v\n", fullPath, err)
	}

	if err := createGzipVersion(fullPath); err != nil {
		log.Printf("Error creating Gzip version of file %s: %v\n", fullPath, err)
	}

	return nil
}

func FetchAndSaveFiles(baseURL string, filePaths []string, targetDir string) error {
	for _, filePath := range filePaths {
		if err := FetchAndSaveFile(baseURL, filePath, targetDir); err != nil {
			fmt.Println("Error:", err)
			continue
		}
	}
	return nil
}
