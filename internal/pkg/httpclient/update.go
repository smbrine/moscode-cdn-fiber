package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"moscode-cdn-fiber/configs"
	"moscode-cdn-fiber/internal/pkg/iosync"
	"net/http"
)

type apiData struct {
	Files []string `json:"files"`
	Pages []string `json:"pages"`
}

func fetchURL(url string, res *apiData) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
		}
	}()

	// Ensure we got a successful response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal the JSON data into the update variable
	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("error unmarshalling response JSON: %w", err)
	}

	return nil
}

func UpdateStatic(url string) error {
	var TempDir = configs.GetTempDir()
	var StaticDir = configs.GetStaticDir()
	var update apiData

	if err := fetchURL(url, &update); err != nil {
		log.Printf("Error fetching url: %v", err)
		return nil
	}
	configs.SetUrlPages(update.Pages)

	if err := iosync.EmptyDir(TempDir); err != nil {
		log.Printf("Error emptying temp directory: %v", err)
		return nil
	}

	if err := iosync.FetchAndSaveFiles(configs.GetBaseURL(), update.Files, TempDir); err != nil {
		log.Printf("Error fetching and saving files: %v", err)
		return nil
	}

	if err := iosync.EmptyDir(StaticDir); err != nil {
		log.Printf("Error emptying static directory: %v", err)
		return nil
	}

	if err := iosync.CopyFilesInDir(TempDir, StaticDir); err != nil {
		log.Printf("Error copying files: %v", err)
		return nil
	}
	return nil
}
