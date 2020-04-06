package tagger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func downloadDb(apiKey, apiHost, datasetID string) string {
	dbPath := "in.sqlite"
	fmt.Printf("Downloading sqlite db to %s\n", dbPath)

	path := "admin/export/v1/dataset"
	values := url.Values{"datasetID": []string{datasetID}}
	err := download(apiKey, apiHost, path, values, dbPath)
	if err != nil {
		log.Fatalf("Failed to download DB for dataset: %v", err)
	}

	return dbPath
}

func download(apiKey, host, path string, q url.Values, output string) error {

	// Construct request.
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}
	u.RawQuery = q.Encode()

	log.Printf("Getting %s", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// Create the file
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", output, err)
	}
	defer out.Close()

	// Get the data
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download faled: [%d] %s", resp.StatusCode, resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}

	return nil
}
