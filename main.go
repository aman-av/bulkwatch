package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// NSE Bulk Deal CSV URL
	nseBaseURL    = "https://nsearchives.nseindia.com/content/equities/bulk.csv"
	archiveFolder = "archive"
)

// fetchBulkDealCSV fetches the bulk deal CSV from NSE website
func fetchBulkDealCSV() ([]byte, error) {
	// Create a custom transport to force HTTP/1.1
	transport := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	// Create request
	req, err := http.NewRequest("GET", nseBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers for NSE
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CSV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Check if response is gzip compressed
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" || resp.Header.Get("Content-Type") == "application/x-gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	// Read response body (decompressed if needed)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}

// ensureArchiveFolder creates the archive folder if it doesn't exist
func ensureArchiveFolder() error {
	err := os.MkdirAll(archiveFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create archive folder: %w", err)
	}
	return nil
}

// fileExists checks if a file exists
func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

// saveCSVToFile saves the CSV data to a file in the archive folder
func saveCSVToFile(data []byte, filename string) error {
	filepath := fmt.Sprintf("%s/%s", archiveFolder, filename)
	err := os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

// readAndDisplayCSV reads and displays the CSV content from archive folder
func readAndDisplayCSV(filename string) error {
	filepath := fmt.Sprintf("%s/%s", archiveFolder, filename)
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Allow variable number of fields and handle quotes more flexibly
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	fmt.Println("\n=== NSE Bulk Deal Data ===")
	fmt.Println("\nColumns:", header)
	fmt.Println("\nData:")
	fmt.Println("---")

	// Read and display all records
	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		recordCount++
		fmt.Printf("\nRecord %d:\n", recordCount)
		for i, value := range record {
			if i < len(header) {
				fmt.Printf("  %s: %s\n", header[i], value)
			}
		}
	}

	fmt.Printf("\n---\nTotal records: %d\n", recordCount)
	return nil
}

func main() {
	// Ensure archive folder exists
	err := ensureArchiveFolder()
	if err != nil {
		fmt.Printf("Error creating archive folder: %v\n", err)
		os.Exit(1)
	}

	// Generate filename for today's data
	filename := fmt.Sprintf("bulk_deals_%s.csv", time.Now().Format("2006-01-02"))
	filepath := fmt.Sprintf("%s/%s", archiveFolder, filename)

	// Check if file already exists in archive
	if fileExists(filepath) {
		fmt.Printf("File already exists in archive: %s\n", filepath)
		fmt.Println("Reading from existing file...")
	} else {
		fmt.Println("Fetching NSE Bulk Deal CSV...")

		// Fetch CSV data
		data, err := fetchBulkDealCSV()
		if err != nil {
			fmt.Printf("Error fetching CSV: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully fetched %d bytes\n", len(data))

		// Save to archive folder
		err = saveCSVToFile(data, filename)
		if err != nil {
			fmt.Printf("Error saving CSV: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Saved to archive: %s\n", filepath)
	}

	// Read and display CSV content
	err = readAndDisplayCSV(filename)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		os.Exit(1)
	}
}

// Made with Bob
