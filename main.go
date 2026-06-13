package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"time"

	"github.com/andybalholm/brotli"
)

const (
	archiveFolder = "archive"
)

// Exchange represents a stock exchange
type Exchange string

const (
	NSE Exchange = "NSE"
	BSE Exchange = "BSE"
)

// DealType represents the type of deal
type DealType string

const (
	BulkDeal  DealType = "bulk"
	BlockDeal DealType = "block"
)

// Config holds the configuration for fetching deals
type Config struct {
	Exchange  Exchange
	DealType  DealType
	FromDate  string
	ToDate    string
}

// NSEDealData represents NSE API response structure
type NSEDealData struct {
	Data []map[string]interface{} `json:"data"`
}

// BSEDealData represents BSE API response structure
type BSEDealData struct {
	Table []map[string]interface{} `json:"Table"`
}

// buildURL constructs the API URL based on exchange and deal type
func buildURL(config Config) string {
	switch config.Exchange {
	case NSE:
		optionType := "bulk_deals"
		if config.DealType == BlockDeal {
			optionType = "block_deals"
		}
		return BuildNSEBulkBlockURL(optionType, config.FromDate, config.ToDate)
	case BSE:
		dealType := "1" // bulk
		if config.DealType == BlockDeal {
			dealType = "2"
		}
		return BuildBSEBulkBlockURL(dealType, config.FromDate, config.ToDate)
	default:
		return ""
	}
}

// fetchDealData fetches deal data from the API
func fetchDealData(config Config) ([]byte, error) {
	url := buildURL(config)
	
	// Create cookie jar for BSE
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	
	client := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	
	// Exchange-specific headers
	if config.Exchange == NSE {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	} else if config.Exchange == BSE {
		req.Header.Set("Referer", BSEReferer)
		req.Header.Set("Origin", BSEOrigin)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read raw data first
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check content encoding and decompress accordingly
	contentEncoding := resp.Header.Get("Content-Encoding")
	
	if contentEncoding == "br" || (len(data) > 0 && data[0] == 0xce) {
		// Brotli compression (NSE uses this)
		brReader := brotli.NewReader(bytes.NewReader(data))
		data, err = io.ReadAll(brReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress brotli data: %w", err)
		}
	} else if contentEncoding == "gzip" || (len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b) {
		// Gzip compression
		gzipReader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		
		data, err = io.ReadAll(gzipReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress gzip data: %w", err)
		}
	}

	return data, nil
}

// ensureArchiveFolder creates the archive folder structure
func ensureArchiveFolder(exchange Exchange, dealType DealType) (string, error) {
	folderPath := filepath.Join(archiveFolder, string(exchange), string(dealType))
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create archive folder: %w", err)
	}
	return folderPath, nil
}

// saveJSONToFile saves JSON data to a file
func saveJSONToFile(data []byte, folderPath, filename string) error {
	filepath := filepath.Join(folderPath, filename)
	err := os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

// parseAndDisplayData parses and displays the JSON data
func parseAndDisplayData(data []byte, config Config) error {
	fmt.Printf("\n=== %s %s Deal Data ===\n", config.Exchange, config.DealType)
	fmt.Printf("Date Range: %s to %s\n\n", config.FromDate, config.ToDate)

	var recordCount int
	
	switch config.Exchange {
	case NSE:
		var nseData NSEDealData
		if err := json.Unmarshal(data, &nseData); err != nil {
			return fmt.Errorf("failed to parse NSE data: %w", err)
		}
		recordCount = len(nseData.Data)
		
		if recordCount > 0 {
			fmt.Println("Sample Records (first 3):")
			fmt.Println("---")
			for i, record := range nseData.Data {
				if i >= 3 {
					break
				}
				fmt.Printf("\nRecord %d:\n", i+1)
				for key, value := range record {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
		}
		
	case BSE:
		var bseData BSEDealData
		if err := json.Unmarshal(data, &bseData); err != nil {
			return fmt.Errorf("failed to parse BSE data: %w", err)
		}
		recordCount = len(bseData.Table)
		
		if recordCount > 0 {
			fmt.Println("Sample Records (first 3):")
			fmt.Println("---")
			for i, record := range bseData.Table {
				if i >= 3 {
					break
				}
				fmt.Printf("\nRecord %d:\n", i+1)
				for key, value := range record {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
		}
	}

	fmt.Printf("\n---\nTotal records: %d\n", recordCount)
	return nil
}

// formatDateForExchange formats date based on exchange requirements
func formatDateForExchange(date time.Time, exchange Exchange) string {
	switch exchange {
	case NSE:
		return date.Format("02-01-2006") // DD-MM-YYYY
	case BSE:
		return date.Format("02/01/2006") // DD/MM/YYYY
	default:
		return date.Format("02-01-2006")
	}
}

// processDeal fetches and saves deal data for a specific configuration
func processDeal(config Config) error {
	fmt.Printf("\n📊 Processing %s %s deals...\n", config.Exchange, config.DealType)
	
	// Create archive folder
	folderPath, err := ensureArchiveFolder(config.Exchange, config.DealType)
	if err != nil {
		return err
	}

	// Generate filename
	filename := fmt.Sprintf("%s_%s_%s_to_%s.json",
		config.Exchange,
		config.DealType,
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02"))

	// Fetch data
	data, err := fetchDealData(config)
	if err != nil {
		return fmt.Errorf("failed to fetch %s %s data: %w", config.Exchange, config.DealType, err)
	}

	fmt.Printf("✓ Fetched %d bytes\n", len(data))

	// Save to file
	err = saveJSONToFile(data, folderPath, filename)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Saved to: %s\n", filepath.Join(folderPath, filename))

	// Parse and display
	err = parseAndDisplayData(data, config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Get current date
	now := time.Now()
	
	// Start date: January 1, 2026
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	
	// Define all configurations
	configs := []Config{
		{
			Exchange: NSE,
			DealType: BulkDeal,
			FromDate: formatDateForExchange(startDate, NSE),
			ToDate:   formatDateForExchange(now, NSE),
		},
		{
			Exchange: NSE,
			DealType: BlockDeal,
			FromDate: formatDateForExchange(startDate, NSE),
			ToDate:   formatDateForExchange(now, NSE),
		},
		{
			Exchange: BSE,
			DealType: BulkDeal,
			FromDate: formatDateForExchange(startDate, BSE),
			ToDate:   formatDateForExchange(now, BSE),
		},
		{
			Exchange: BSE,
			DealType: BlockDeal,
			FromDate: formatDateForExchange(startDate, BSE),
			ToDate:   formatDateForExchange(now, BSE),
		},
	}

	fmt.Println("🚀 FlowWatch - Bulk & Block Deal Tracker")
	fmt.Println("========================================")
	fmt.Printf("Fetching data from %s to %s\n", startDate.Format("2006-01-02"), now.Format("2006-01-02"))

	// Process each configuration
	successCount := 0
	for _, config := range configs {
		err := processDeal(config)
		if err != nil {
			fmt.Printf("❌ Error processing %s %s: %v\n", config.Exchange, config.DealType, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("\n✅ Successfully processed %d/%d data sources\n", successCount, len(configs))
	
	// Enrich all deals with shareholding classification
	if err := enrichAllDeals(); err != nil {
		fmt.Printf("⚠️  Error enriching deals: %v\n", err)
	}
}

// Made with Bob
